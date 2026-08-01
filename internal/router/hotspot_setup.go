package router

import (
	"fmt"
	"net"
	"regexp"
)

// ─── Preflight types ──────────────────────────────────────────────────────────

type InterfaceInfo struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Running bool   `json:"running"`
	Comment string `json:"comment"`
}

type PreflightResult struct {
	Interfaces      []InterfaceInfo `json:"interfaces"`
	HotspotExists   bool            `json:"hotspotExists"`
	HotspotOnIface  string          `json:"hotspotOnIface"`
	HotspotName     string          `json:"hotspotName"`
	HotspotProfile  string          `json:"hotspotProfile"`
	HotspotDNSName  string          `json:"hotspotDnsName"`
	HotspotAddPool  string          `json:"hotspotAddressPool"`
}

// ─── Setup types ──────────────────────────────────────────────────────────────

type SetupRequest struct {
	LANIface    string `json:"lanIface"`
	WANIface    string `json:"wanIface"`
	Subnet      string `json:"subnet"`
	HotspotName string `json:"hotspotName"` // user-chosen label, e.g. "myspot" → DNS: myspot.spot
}

type SetupStepResult struct {
	Name    string `json:"name"`
	OK      bool   `json:"ok"`
	Error   string `json:"error,omitempty"`
	Skipped bool   `json:"skipped,omitempty"`
}

type SetupResult struct {
	Steps   []SetupStepResult `json:"steps"`
	Success bool              `json:"success"`
}

// ─── Preflight ────────────────────────────────────────────────────────────────

// HotspotPreflight checks for existing hotspot config and lists interfaces.
func (c *Client) HotspotPreflight() (*PreflightResult, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Fetch all interfaces
	ifaceReply, err := conn.Run("/interface/print")
	if err != nil {
		return nil, fmt.Errorf("interface/print: %w", err)
	}

	var ifaces []InterfaceInfo
	for _, s := range ifaceReply.Re {
		ifaces = append(ifaces, InterfaceInfo{
			Name:    s.Map["name"],
			Type:    s.Map["type"],
			Running: s.Map["running"] == "true",
			Comment: s.Map["comment"],
		})
	}

	// Check for existing hotspot
	hsReply, err := conn.Run("/ip/hotspot/print")
	if err != nil {
		return nil, fmt.Errorf("hotspot/print: %w", err)
	}

	result := &PreflightResult{Interfaces: ifaces}
	if len(hsReply.Re) == 0 {
		return result, nil
	}

	hs := hsReply.Re[0].Map
	result.HotspotExists = true
	result.HotspotOnIface = hs["interface"]
	result.HotspotName = hs["name"]
	result.HotspotProfile = hs["profile"]
	result.HotspotAddPool = hs["address-pool"]

	// Fetch the hotspot server profile to get dns-name
	profileName := result.HotspotProfile
	if profileName == "" {
		profileName = "default" // RouterOS default profile name
	}
	profReply, err := conn.RunArgs([]string{
		"/ip/hotspot/profile/print",
		"?name=" + profileName,
	})
	if err == nil && len(profReply.Re) > 0 {
		result.HotspotDNSName = profReply.Re[0].Map["dns-name"]
	}

	return result, nil
}

// ─── Setup orchestrator ───────────────────────────────────────────────────────

var validIfaceName = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]+$`)
var validHotspotName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9\-]{0,61}$`)

// SetupHotspot runs all setup steps and always returns a populated SetupResult.
// It returns (nil, err) only for a connection failure or invalid subnet.
func (c *Client) SetupHotspot(req SetupRequest) (*SetupResult, error) {
	if !validIfaceName.MatchString(req.LANIface) || !validIfaceName.MatchString(req.WANIface) {
		return nil, fmt.Errorf("invalid interface name")
	}

	gw, poolStart, poolEnd, ipWithPrefix, networkAddr, err := subnetDetails(req.Subnet)
	if err != nil {
		return &SetupResult{
			Steps: []SetupStepResult{{Name: "Parse subnet", OK: false, Error: err.Error()}},
		}, nil
	}

	// Build DNS name: user label + forced .spot TLD
	dnsName := ""
	if req.HotspotName != "" {
		if !validHotspotName.MatchString(req.HotspotName) {
			return nil, fmt.Errorf("invalid hotspot name: use letters, numbers and hyphens only")
		}
		dnsName = req.HotspotName + ".spot"
	}

	poolName := "hotspot-pool"
	profileName := "pikro-profile"
	dhcpName := "hotspot-dhcp"
	hotspotName := "hotspot1"

	var steps []SetupStepResult
	run := func(s SetupStepResult) {
		steps = append(steps, s)
	}

	run(c.stepAssignIP(req.LANIface, ipWithPrefix))
	run(c.stepCreateIPPool(poolName, poolStart, poolEnd))
	run(c.stepCreateDHCPNetwork(networkAddr, gw))
	run(c.stepCreateDHCPServer(dhcpName, req.LANIface, poolName))
	run(c.stepCreateHotspotProfile(profileName, gw, dnsName))
	run(c.stepEnableHotspot(hotspotName, req.LANIface, poolName, profileName))
	// Walled garden immediately after hotspot activation so port 8728 is
	// whitelisted before any further API calls — otherwise the captive portal
	// blocks the connection for all remaining steps.
	run(c.stepAddWalledGarden())
	// Upload the custom login page after the hotspot server is running —
	// the hotspot/ directory only exists once the hotspot server is active.
	run(c.stepUploadLoginPage(profileName))
	run(c.stepAddNATMasquerade(req.WANIface))
	run(c.stepEnableDNS())
	run(c.stepInstallCleanupScheduler())

	success := true
	for _, s := range steps {
		if !s.OK && !s.Skipped {
			success = false
			break
		}
	}
	return &SetupResult{Steps: steps, Success: success}, nil
}

// ─── Setup steps ─────────────────────────────────────────────────────────────

// isAlreadyExists returns true for RouterOS "already have such" / "already exists" errors.
// These mean the resource is already present — we treat them as Skipped (idempotent).
func isAlreadyExists(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	return containsStr(msg, "already have such") ||
		containsStr(msg, "already exists") ||
		containsStr(msg, "such network already exists") ||
		containsStr(msg, "server or relay with such interface already exists") ||
		containsStr(msg, "pool with such name exists")
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func (c *Client) stepAssignIP(iface, ipWithPrefix string) SetupStepResult {
	name := "Assign IP to LAN interface"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{"/ip/address/add", "=interface=" + iface, "=address=" + ipWithPrefix})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepCreateIPPool(poolName, start, end string) SetupStepResult {
	name := "Create IP pool"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{"/ip/pool/add", "=name=" + poolName, "=ranges=" + start + "-" + end})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepCreateDHCPNetwork(networkCIDR, gateway string) SetupStepResult {
	name := "Create DHCP network"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{
		"/ip/dhcp-server/network/add",
		"=address=" + networkCIDR,
		"=gateway=" + gateway,
		"=dns-server=" + gateway,
	})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepCreateDHCPServer(dhcpName, iface, poolName string) SetupStepResult {
	name := "Create DHCP server"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{
		"/ip/dhcp-server/add",
		"=name=" + dhcpName,
		"=interface=" + iface,
		"=address-pool=" + poolName,
		"=disabled=no",
	})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepCreateHotspotProfile(profileName, gateway, dnsName string) SetupStepResult {
	name := "Create hotspot profile"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()

	// http-pap is included alongside http-chap so the "Business" voucher
	// template's QR code can encode username/password as plain URL query
	// params — chap requires a per-session challenge hash a QR scan can't
	// provide, so without pap a QR-based auto-login silently never works.
	args := []string{
		"/ip/hotspot/profile/add",
		"=name=" + profileName,
		"=hotspot-address=" + gateway,
		"=login-by=cookie,http-chap,http-pap,mac-cookie",
	}
	if dnsName != "" {
		args = append(args, "=dns-name="+dnsName)
	}
	if _, err = conn.RunArgs(args); err != nil {
		if !isAlreadyExists(err) {
			return SetupStepResult{Name: name, Error: err.Error()}
		}
		// Profile exists — update dns-name and login-by so re-runs apply changes.
		setArgs := []string{
			"/ip/hotspot/profile/set",
			"=numbers=" + profileName,
			"=login-by=cookie,http-chap,http-pap,mac-cookie",
		}
		if dnsName != "" {
			setArgs = append(setArgs, "=dns-name="+dnsName)
		}
		if _, serr := conn.RunArgs(setArgs); serr != nil {
			return SetupStepResult{Name: name, Error: serr.Error()}
		}
		return SetupStepResult{Name: name, OK: true, Skipped: true}
	}
	return SetupStepResult{Name: name, OK: true}
}

// stepInstallCleanupScheduler installs the cleanup scheduler with a weekly interval.
// Non-fatal if it fails (e.g. device-mode restriction).
func (c *Client) stepInstallCleanupScheduler() SetupStepResult {
	name := "Install cleanup scheduler"
	if err := c.InstallCleanupScheduler("7d"); err != nil {
		return SetupStepResult{Name: name, OK: true, Skipped: true}
	}
	return SetupStepResult{Name: name, OK: true}
}

// stepUploadLoginPage uploads the Pikro login page after the hotspot server
// is active (the hotspot/ directory only exists once /ip/hotspot/add has run).
func (c *Client) stepUploadLoginPage(profileName string) SetupStepResult {
	name := "Upload login page"
	if err := c.UploadLoginPage(profileName, LoginPageParams{}); err != nil {
		// Non-fatal: the hotspot works with the default page if this fails.
		return SetupStepResult{Name: name, OK: true, Skipped: true}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepAddWalledGarden() SetupStepResult {
	name := "Allow router management (walled garden)"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()

	// Allow unauthenticated clients to reach the router's management ports
	// so admins are never locked out after hotspot is enabled.
	rules := [][]string{
		{"/ip/hotspot/walled-garden/ip/add", "=action=accept", "=dst-port=8728", "=protocol=tcp", "=comment=pikro-api"},
		{"/ip/hotspot/walled-garden/ip/add", "=action=accept", "=dst-port=8291", "=protocol=tcp", "=comment=pikro-winbox"},
	}
	for _, args := range rules {
		if _, err = conn.RunArgs(args); err != nil && !isAlreadyExists(err) {
			return SetupStepResult{Name: name, Error: err.Error()}
		}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepEnableHotspot(hotspotName, iface, poolName, profileName string) SetupStepResult {
	name := "Enable hotspot on interface"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{
		"/ip/hotspot/add",
		"=name=" + hotspotName,
		"=interface=" + iface,
		"=address-pool=" + poolName,
		"=profile=" + profileName,
		"=disabled=no",
	})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepAddNATMasquerade(wanIface string) SetupStepResult {
	name := "Add NAT masquerade"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{
		"/ip/firewall/nat/add",
		"=chain=srcnat",
		"=action=masquerade",
		"=out-interface=" + wanIface,
		"=comment=pikro-hotspot-masq",
	})
	if err != nil {
		if isAlreadyExists(err) {
			return SetupStepResult{Name: name, OK: true, Skipped: true}
		}
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

func (c *Client) stepEnableDNS() SetupStepResult {
	name := "Enable DNS remote requests"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	defer conn.Close()
	_, err = conn.RunArgs([]string{"/ip/dns/set", "=allow-remote-requests=yes"})
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}
	}
	return SetupStepResult{Name: name, OK: true}
}

// ─── Teardown (dev/reset) ─────────────────────────────────────────────────────

type TeardownResult struct {
	Steps []SetupStepResult `json:"steps"`
}

// TeardownHotspot removes only the captive-portal layer created by Pikro:
// the hotspot server, its profile, and the NAT masquerade rule.
// It intentionally leaves IP address, DHCP server/network, IP pool, and DNS
// untouched — those are infrastructure the router needs to stay reachable.
func (c *Client) TeardownHotspot() (*TeardownResult, error) {
	var steps []SetupStepResult

	steps = append(steps, c.teardownRemove("/ip/hotspot", "name", "hotspot1", "Remove hotspot server"))
	steps = append(steps, c.teardownRemove("/ip/hotspot/profile", "name", "pikro-profile", "Remove hotspot profile"))
	steps = append(steps, c.teardownRemoveNAT("Remove NAT masquerade rule"))
	steps = append(steps, c.teardownRemoveWalledGarden("Remove walled garden rules"))

	return &TeardownResult{Steps: steps}, nil
}

func (c *Client) teardownRemove(path, key, val, label string) SetupStepResult {
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: label, Error: err.Error()}
	}
	defer conn.Close()

	reply, err := conn.RunArgs([]string{path + "/print", "?.id", "?" + key + "=" + val})
	if err != nil || len(reply.Re) == 0 {
		return SetupStepResult{Name: label, OK: true, Skipped: true}
	}
	id := reply.Re[0].Map[".id"]
	if _, err = conn.RunArgs([]string{path + "/remove", "=.id=" + id}); err != nil {
		return SetupStepResult{Name: label, Error: err.Error()}
	}
	return SetupStepResult{Name: label, OK: true}
}

func (c *Client) teardownRemoveNAT(label string) SetupStepResult {
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: label, Error: err.Error()}
	}
	defer conn.Close()

	reply, err := conn.RunArgs([]string{"/ip/firewall/nat/print", "?comment=pikro-hotspot-masq"})
	if err != nil || len(reply.Re) == 0 {
		return SetupStepResult{Name: label, OK: true, Skipped: true}
	}
	id := reply.Re[0].Map[".id"]
	if _, err = conn.RunArgs([]string{"/ip/firewall/nat/remove", "=.id=" + id}); err != nil {
		return SetupStepResult{Name: label, Error: err.Error()}
	}
	return SetupStepResult{Name: label, OK: true}
}

func (c *Client) teardownRemoveWalledGarden(label string) SetupStepResult {
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: label, Error: err.Error()}
	}
	defer conn.Close()

	for _, comment := range []string{"pikro-api", "pikro-winbox"} {
		reply, err := conn.RunArgs([]string{"/ip/hotspot/walled-garden/ip/print", "?comment=" + comment})
		if err != nil || len(reply.Re) == 0 {
			continue
		}
		id := reply.Re[0].Map[".id"]
		conn.RunArgs([]string{"/ip/hotspot/walled-garden/ip/remove", "=.id=" + id})
	}
	return SetupStepResult{Name: label, OK: true}
}

// ─── Subnet helpers ───────────────────────────────────────────────────────────

// subnetDetails parses a CIDR like "192.168.88.0/24" and derives all values
// needed to configure the hotspot network. Assumes /24 or larger prefix.
func subnetDetails(cidr string) (gateway, poolStart, poolEnd, ipWithPrefix, networkAddr string, err error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("invalid CIDR %q: %w", cidr, err)
	}
	_ = ip

	// Network base IP (e.g. 192.168.88.0)
	base := ipNet.IP.To4()
	if base == nil {
		return "", "", "", "", "", fmt.Errorf("only IPv4 supported")
	}

	networkAddr = ipNet.String() // "192.168.88.0/24"

	// Gateway = base + 1 (e.g. 192.168.88.1)
	gwIP := make(net.IP, 4)
	copy(gwIP, base)
	gwIP[3] = base[3] + 1
	gateway = gwIP.String()

	ones, bits := ipNet.Mask.Size()
	totalHosts := (1 << uint(bits-ones)) - 2 // excludes network + broadcast

	// Pool start = base + 2
	startIP := make(net.IP, 4)
	copy(startIP, base)
	startIP[3] = base[3] + 2
	poolStart = startIP.String()

	// Pool end = broadcast - 1
	endIP := make(net.IP, 4)
	copy(endIP, base)
	// broadcast is base | ~mask
	for i := range endIP {
		endIP[i] = base[i] | ^ipNet.Mask[i]
	}
	endIP[3]-- // one before broadcast
	poolEnd = endIP.String()

	// Sanity: need at least one host in pool
	if totalHosts < 2 {
		return "", "", "", "", "", fmt.Errorf("subnet too small")
	}

	prefix, _ := ipNet.Mask.Size()
	ipWithPrefix = fmt.Sprintf("%s/%d", gateway, prefix)

	return gateway, poolStart, poolEnd, ipWithPrefix, networkAddr, nil
}
