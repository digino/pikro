package router

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

// Client wraps go-routeros to talk to the RouterOS native API (port 8728/8729).
// This is the same protocol WinBox uses — more reliable than the REST API and
// works on RouterOS v6 as well as v7.
type Client struct {
	host     string
	port     int
	username string
	password string
	useTLS   bool
}

func NewClient(host string, port int, username, password string, useTLS bool) *Client {
	return &Client{host: host, port: port, username: username, password: password, useTLS: useTLS}
}

func (c *Client) address() string {
	return fmt.Sprintf("%s:%d", c.host, c.port)
}

func (c *Client) connect() (*ros.Client, error) {
	addr := c.address()
	timeout := 10 * time.Second
	if c.useTLS {
		cfg := &tls.Config{InsecureSkipVerify: true}
		return ros.DialTLSTimeout(addr, c.username, c.password, cfg, timeout)
	}
	return ros.DialTimeout(addr, c.username, c.password, timeout)
}

// Ping tests reachability by fetching /system/identity.
func (c *Client) Ping() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Run("/system/identity/print")
	return err
}

// SystemResource returns CPU load, free memory, total memory, and uptime.
func (c *Client) SystemResource() (map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/system/resource/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return map[string]any{}, nil
	}
	return sentenceToMap(reply.Re[0].Map), nil
}

// HotspotUsers returns all hotspot users.
func (c *Client) HotspotUsers() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/hotspot/user/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// HotspotActive returns active hotspot sessions.
func (c *Client) HotspotActive() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/hotspot/active/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// HotspotHosts returns all devices the hotspot has ever seen on the network
// (/ip/hotspot/host) — MAC/IP/server/bridge-port for any connected client,
// not just currently-authenticated sessions (that's HotspotActive).
func (c *Client) HotspotHosts() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/hotspot/host/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// HotspotProfiles returns available hotspot user profiles.
func (c *Client) HotspotProfiles() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/hotspot/user/profile/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// AddressPools returns configured IP address pools (/ip/pool), so an admin
// can pick which pool new hotspot users are assigned to.
func (c *Client) AddressPools() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/pool/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// HotspotProfileParams holds the RouterOS fields for a hotspot user profile.
// ValiditySecs, if > 0, causes an on-login script to be set on the profile that
// stamps exp:<epoch> on the user's comment the first time they log in.
type HotspotProfileParams struct {
	Name         string
	AddressPool  string
	SharedUsers  string // e.g. "1", "2"
	RateLimit    string // e.g. "2M/4M"
	ValiditySecs int64  // 0 = no expiry, >0 = set on-login script
	ROSVersion   string // detected from router, used to pick on-login script variant
}

// buildOnLoginScript returns a RouterOS on-login script that stamps
// exp:<value> into the user comment the first time they log in, preserving
// any existing operator note by appending it after the marker (space-
// separated) rather than overwriting the whole field. extractExpEpoch and
// isUptimeExhausted (frontend) and the cleanup script (RouterOS side) only
// ever read the exp:<value> prefix, so anything after it is safe to carry
// along unexamined.
// On ROS 7.12+ it uses Unix epoch seconds; on v6/early v7 it uses
// "YYYY-MM-DD HH:MM:SS" which the v6 cleanup script parses as an integer.
func buildOnLoginScript(validitySecs int64, rosVersion string) string {
	major, minor := 7, 0
	fmt.Sscanf(rosVersion, "%d.%d", &major, &minor)
	if major > 7 || (major == 7 && minor >= 12) {
		return fmt.Sprintf(
			`:local nowEpoch ([:tonsec [:timestamp]] / 1000000000)
:local expEpoch ($nowEpoch + %d)
:local uid [/ip/hotspot/user/find where name=$user]
:if ([:len $uid] > 0) do={
  :local c [:tostr [/ip/hotspot/user/get $uid comment]]
  :if ([:len $c] < 5 || [:pick $c 0 4] != "exp:") do={
    :local newComment ("exp:" . $expEpoch)
    :if ([:len $c] > 0) do={ :set newComment ($newComment . " " . $c) }
    :log info ("pikro-on-login: stamping $user (exp=$expEpoch now=$nowEpoch)")
    /ip/hotspot/user/set $uid comment=$newComment
  }
}`, validitySecs)
	}
	// v6 fallback: use Mikhmon's proven temp-scheduler trick to get "now + validity"
	// as a RouterOS date string, then normalise it to YYYY-MM-DD HH:MM:SS.
	// A temporary scheduler named after the user is created with interval=validity,
	// its next-run field gives us the exact expiry datetime, then we remove it.
	return fmt.Sprintf(
		`:local uid [/ip/hotspot/user/find where name=$user]
:if ([:len $uid] > 0) do={
  :local c [:tostr [/ip/hotspot/user/get $uid comment]]
  :if ([:len $c] < 5 || [:pick $c 0 4] != "exp:") do={
    :local nd [/system clock get date]
    :local year [:pick $nd 7 11]
    /system scheduler add name=$user disabled=no start-date=$nd interval="%ds"
    :delay 1s
    :local exp [/system scheduler get [/system scheduler find where name=$user] next-run]
    :local xlen [:len $exp]
    :local expStr ""
    :if ($xlen = 15) do={
      :local mo [:tolower [:pick $exp 0 3]]
      :local months {jan="01";feb="02";mar="03";apr="04";may="05";jun="06";jul="07";aug="08";sep="09";oct="10";nov="11";dec="12"}
      :set expStr ($year."-".($months->$mo)."-".[:pick $exp 4 6]." ".[:pick $exp 7 15])
    }
    :if ($xlen > 15) do={ :set expStr $exp }
    :if ($expStr != "") do={
      :local newComment ("exp:" . $expStr)
      :if ([:len $c] > 0) do={ :set newComment ($newComment . " " . $c) }
      :log info ("pikro-on-login: stamping $user (exp=$expStr)")
      /ip/hotspot/user/set $uid comment=$newComment
    }
    /system scheduler remove [find where name=$user]
  }
}`, validitySecs)
}

func (c *Client) rosVersion(conn *ros.Client) string {
	if r, e := conn.Run("/system/resource/print"); e == nil && len(r.Re) > 0 {
		return r.Re[0].Map["version"]
	}
	return "7.0"
}

func (c *Client) CreateHotspotProfile(p HotspotProfileParams) (map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	p.ROSVersion = c.rosVersion(conn)
	args := append(buildProfileArgs("/ip/hotspot/user/profile/add", p), "=name="+p.Name)
	reply, err := conn.RunArgs(args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return map[string]any{}, nil
	}
	return sentenceToMap(reply.Re[0].Map), nil
}

func (c *Client) UpdateHotspotProfile(id string, p HotspotProfileParams) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	p.ROSVersion = c.rosVersion(conn)
	args := buildProfileArgs("/ip/hotspot/user/profile/set", p)
	args = append(args, "=.id="+id)
	_, err = conn.RunArgs(args)
	return err
}

func buildProfileArgs(path string, p HotspotProfileParams) []string {
	// Note: =name= is intentionally omitted here; callers add it only on create.
	args := []string{path}
	if p.AddressPool != "" {
		args = append(args, "=address-pool="+p.AddressPool)
	}
	if p.SharedUsers != "" {
		args = append(args, "=shared-users="+p.SharedUsers)
	}
	if p.RateLimit != "" {
		args = append(args, "=rate-limit="+p.RateLimit)
	}
	if p.ValiditySecs > 0 {
		args = append(args, "=on-login="+buildOnLoginScript(p.ValiditySecs, p.ROSVersion))
	} else {
		args = append(args, "=on-login=")
	}
	return args
}

func (c *Client) DeleteHotspotProfile(id string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Run("/ip/hotspot/user/profile/remove", "=.id="+id)
	return err
}

// HotspotUserParams holds all fields for creating a hotspot user.
type HotspotUserParams struct {
	Name            string
	Password        string
	Profile         string
	LimitUptime     string // e.g. "2h", "30m"
	LimitBytesTotal string // bytes as string, e.g. "1073741824"
	RateLimit       string // e.g. "2M/4M" (up/down)
	Comment         string
}

// CreateHotspotUser creates a new hotspot user with optional limits.
func (c *Client) CreateHotspotUser(p HotspotUserParams) (map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	args := []string{
		"/ip/hotspot/user/add",
		"=name=" + p.Name,
		"=password=" + p.Password,
	}
	if p.Profile != "" {
		args = append(args, "=profile="+p.Profile)
	}
	if p.LimitUptime != "" {
		args = append(args, "=limit-uptime="+p.LimitUptime)
	}
	if p.LimitBytesTotal != "" {
		args = append(args, "=limit-bytes-total="+p.LimitBytesTotal)
	}
	if p.RateLimit != "" {
		args = append(args, "=rate-limit="+p.RateLimit)
	}
	if p.Comment != "" {
		args = append(args, "=comment="+p.Comment)
	}
	reply, err := conn.RunArgs(args)
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return map[string]any{}, nil
	}
	return sentenceToMap(reply.Re[0].Map), nil
}

// ToggleHotspotUser sets disabled=yes or disabled=no on a hotspot user.
func (c *Client) ToggleHotspotUser(id string, disabled bool) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	val := "no"
	if disabled {
		val = "yes"
	}
	_, err = conn.Run("/ip/hotspot/user/set", "=.id="+id, "=disabled="+val)
	return err
}

// UpdateHotspotUserParams holds the editable fields for a hotspot user.
//
// Comment is intentionally not editable here: it doubles as Pikro's own
// exp:<value> expiry marker once a user has logged in (see
// buildOnLoginScript), and the Edit dialog only ever showed/edited the
// operator-note half of it — so a plain overwrite would silently discard
// the marker on every edit, permanently disabling time-based cleanup for
// that user. Until there's a design that safely preserves the marker while
// editing the note, comment editing is removed rather than left as a latent
// data-loss trap.
type UpdateHotspotUserParams struct {
	Password        string
	Profile         string
	LimitUptime     string
	LimitBytesTotal string
}

// UpdateHotspotUser updates an existing hotspot user by RouterOS ID.
func (c *Client) UpdateHotspotUser(id string, p UpdateHotspotUserParams) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	args := []string{"/ip/hotspot/user/set", "=.id=" + id}
	if p.Password != "" {
		args = append(args, "=password="+p.Password)
	}
	if p.Profile != "" {
		args = append(args, "=profile="+p.Profile)
	}
	// Always send limit-uptime so it can be cleared (empty string = no limit)
	args = append(args, "=limit-uptime="+p.LimitUptime)
	// Only send limit-bytes-total when non-empty — RouterOS requires an integer
	if p.LimitBytesTotal != "" {
		args = append(args, "=limit-bytes-total="+p.LimitBytesTotal)
	}
	_, err = conn.RunArgs(args)
	return err
}

// DeleteHotspotUser removes a hotspot user by its RouterOS ID (e.g. "*1").
//
// Deleting /ip/hotspot/user does NOT disconnect an already-authenticated
// session — RouterOS treats /ip/hotspot/active as a separate live-session
// table that only clears on timeout, quota exhaustion, or an explicit
// /ip/hotspot/active/remove. So a user removed here can still appear as
// "active" until their session naturally expires unless we also kick it.
func (c *Client) DeleteHotspotUser(id string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Look up the username before deleting, so we can find and kick any
	// matching active session by name (active sessions are keyed by their
	// own .id, unrelated to the user record's .id).
	var username string
	if reply, err := conn.Run("/ip/hotspot/user/print", "?.id="+id); err == nil && len(reply.Re) > 0 {
		username = reply.Re[0].Map["name"]
	}

	if _, err := conn.Run("/ip/hotspot/user/remove", "=.id="+id); err != nil {
		return err
	}

	if username != "" {
		if reply, err := conn.Run("/ip/hotspot/active/print", "?user="+username); err == nil {
			for _, re := range reply.Re {
				conn.Run("/ip/hotspot/active/remove", "=.id="+re.Map[".id"])
			}
		}
	}
	return nil
}

// DisconnectHotspotActive kicks a live session from /ip/hotspot/active by
// its own .id — this only ends the current connection, it does not touch
// the underlying /ip/hotspot/user record (the user can log back in).
func (c *Client) DisconnectHotspotActive(id string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Run("/ip/hotspot/active/remove", "=.id="+id)
	return err
}

// PollSnapshot holds all data collected in a single RouterOS connection for the dashboard poll.
type PollSnapshot struct {
	Resource   map[string]string   `json:"resource"`
	Traffic    []map[string]string `json:"traffic"`
	Addresses  []map[string]string `json:"addresses"`
	Clock      map[string]string   `json:"clock"`
	Interfaces []map[string]string `json:"interfaces"`
}

// Poll opens one connection and fetches resource, interface traffic, and clock together.
// monitor-traffic with =once= blocks ~1s on the router side before returning.
func (c *Client) Poll() (*PollSnapshot, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	snap := &PollSnapshot{
		Resource:   map[string]string{},
		Traffic:    []map[string]string{},
		Addresses:  []map[string]string{},
		Clock:      map[string]string{},
		Interfaces: []map[string]string{},
	}

	// Resource
	if r, err := conn.Run("/system/resource/print"); err == nil && len(r.Re) > 0 {
		snap.Resource = r.Re[0].Map
	}

	// Interfaces — link status (running/disabled/link-downs) plus names for traffic monitoring
	ifReply, err := conn.Run("/interface/print")
	if err == nil && len(ifReply.Re) > 0 {
		names := make([]string, 0, len(ifReply.Re))
		for _, s := range ifReply.Re {
			snap.Interfaces = append(snap.Interfaces, s.Map)
			if n := s.Map["name"]; n != "" {
				names = append(names, n)
			}
		}
		if len(names) > 0 {
			args := []string{
				"/interface/monitor-traffic",
				"=interface=" + strings.Join(names, ","),
				"=once=",
			}
			if tr, err := conn.Run(args...); err == nil {
				for _, s := range tr.Re {
					snap.Traffic = append(snap.Traffic, s.Map)
				}
			}
		}
	}

	// IP addresses per interface
	if r, err := conn.Run("/ip/address/print"); err == nil {
		for _, s := range r.Re {
			snap.Addresses = append(snap.Addresses, s.Map)
		}
	}

	// Clock
	if r, err := conn.Run("/system/clock/print"); err == nil && len(r.Re) > 0 {
		snap.Clock = r.Re[0].Map
	}

	return snap, nil
}

// WanIP returns the IP address assigned to the WAN interface (the default-route interface).
// Falls back to the first non-loopback, non-bridge address if no default route is found.
func (c *Client) WanIP() (string, error) {
	conn, err := c.connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	// Find the interface that carries the default route.
	routeReply, err := conn.Run("/ip/route/print", "?dst-address=0.0.0.0/0", "?active=true")
	if err != nil {
		return "", err
	}
	wanIface := ""
	if len(routeReply.Re) > 0 {
		wanIface = routeReply.Re[0].Map["interface"]
	}

	addrReply, err := conn.Run("/ip/address/print")
	if err != nil {
		return "", err
	}
	first := ""
	for _, s := range addrReply.Re {
		if s.Map["disabled"] == "true" || s.Map["dynamic"] == "false" {
			continue
		}
		iface := s.Map["interface"]
		addr := s.Map["address"]
		if addr == "" {
			continue
		}
		// Strip prefix length (e.g. "10.10.0.1/24" -> "10.10.0.1")
		if i := strings.Index(addr, "/"); i >= 0 {
			addr = addr[:i]
		}
		if wanIface != "" && iface == wanIface {
			return addr, nil
		}
		if first == "" {
			first = addr
		}
	}
	return first, nil
}

// InterfaceTraffic returns current TX/RX bits-per-second for each interface.
// It calls /interface/monitor-traffic with duration=1s which blocks for ~1s then
// returns one sample. The caller should run this in a goroutine.
func (c *Client) InterfaceTraffic() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// First get all interface names.
	ifReply, err := conn.Run("/interface/print")
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(ifReply.Re))
	for _, s := range ifReply.Re {
		if n := s.Map["name"]; n != "" {
			names = append(names, n)
		}
	}
	if len(names) == 0 {
		return nil, nil
	}

	args := []string{
		"/interface/monitor-traffic",
		"=interface=" + strings.Join(names, ","),
		"=once=",
	}
	reply, err := conn.Run(args...)
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// SystemClock returns the router's current date and time strings.
func (c *Client) SystemClock() (map[string]string, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/system/clock/print")
	if err != nil {
		return nil, err
	}
	if len(reply.Re) == 0 {
		return map[string]string{}, nil
	}
	return reply.Re[0].Map, nil
}

// SystemLogs returns log entries, newest first. If topic is non-empty, only
// entries whose topics field contains it are returned.
func (c *Client) SystemLogs(topic string) ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/log/print")
	if err != nil {
		return nil, err
	}
	entries := replyToList(reply)
	if topic != "" {
		filtered := entries[:0]
		for _, e := range entries {
			if topics, _ := e["topics"].(string); strings.Contains(topics, topic) {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}
	// RouterOS returns oldest first — reverse so newest is first.
	for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
		entries[i], entries[j] = entries[j], entries[i]
	}
	return entries, nil
}

// DHCPLeases returns all leases known to the router's DHCP server(s).
func (c *Client) DHCPLeases() ([]map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	reply, err := conn.Run("/ip/dhcp-server/lease/print")
	if err != nil {
		return nil, err
	}
	return replyToList(reply), nil
}

// BandwidthTest measures internet download speed from the router using /tool/fetch.
// The target parameter selects the file size: "10", "50", or "100" (MB).
// Downloads from Tele2's public speedtest CDN and computes throughput from the
// duration reported by RouterOS.
func (c *Client) BandwidthTest(sizeMB string) (map[string]any, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	type testFile struct {
		url  string
		size int64
	}
	files := map[string]testFile{
		"10":  {"http://speedtest.tele2.net/10MB.zip", 10 * 1024 * 1024},
		"50":  {"http://speedtest.tele2.net/50MB.zip", 50 * 1024 * 1024},
		"100": {"http://speedtest.tele2.net/100MB.zip", 100 * 1024 * 1024},
	}
	f, ok := files[sizeMB]
	if !ok {
		f = files["10"]
	}

	reply, err := conn.Run(
		"/tool/fetch",
		"=url="+f.url,
		"=mode=http",
		"=keep-result=no",
	)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	var durationStr string
	for _, s := range reply.Re {
		if d, ok := s.Map["duration"]; ok && d != "" {
			durationStr = d
			break
		}
	}

	downloadBps := 0
	if durationStr != "" {
		dur := parseROSDuration(durationStr)
		if dur > 0 {
			downloadBps = int(float64(f.size) / dur.Seconds())
		}
	}

	return map[string]any{
		"rx-speed":  fmt.Sprintf("%d", downloadBps),
		"tx-speed":  "0",
		"duration":  durationStr,
		"test-url":  f.url,
		"file-size": fmt.Sprintf("%d", f.size),
	}, nil
}

// parseROSDuration parses RouterOS duration strings like "2s430ms", "834ms", "1m2s".
func parseROSDuration(s string) time.Duration {
	var total time.Duration
	re := regexp.MustCompile(`(\d+)(h|m|s|ms)`)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		n, _ := strconv.Atoi(m[1])
		switch m[2] {
		case "h":
			total += time.Duration(n) * time.Hour
		case "m":
			total += time.Duration(n) * time.Minute
		case "s":
			total += time.Duration(n) * time.Second
		case "ms":
			total += time.Duration(n) * time.Millisecond
		}
	}
	return total
}

// sentenceToMap converts a proto.Sentence word map (all string values) to map[string]any.
func sentenceToMap(m map[string]string) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func replyToList(reply *ros.Reply) []map[string]any {
	out := make([]map[string]any, 0, len(reply.Re))
	for _, s := range reply.Re {
		out = append(out, sentenceToMap(s.Map))
	}
	return out
}

const cleanupScriptName = "pikro-cleanup"

// durationToSecsFn is a RouterOS script snippet defining a local function
// that parses RouterOS's native hotspot uptime/limit-uptime string into
// total seconds. Two forms are known to occur (confirmed via a live
// router's own :log output, which reflects the raw script-level value —
// Pikro's Go API client normalises this to the letter-suffixed form before
// it ever reaches the frontend, so the frontend/API never sees form 2):
//  1. Letter-suffixed, combinable: "1w3d", "1h30m", "11m44s", "1w", "0s".
//  2. Week/day letters followed by a colon-joined time-of-day: "00:00:00",
//     "1d00:00:00", "1w3d10:51:24".
//
// Bug history: an earlier version only handled form 1. A hotspot user's
// uptime/limit-uptime read directly via RouterOS script (form 2, e.g.
// "00:00:00") has no w/d/h/m/s letters at all, so the old parser silently
// returned 0 for every such value — comparing 0 >= 0 was always true,
// deleting every unused voucher regardless of its actual quota. This was
// invisible through Pikro's own API (which shows form 1) and only caught by
// reading RouterOS's own :log line directly on the router.
// Shared by both cleanup script variants. Assumes the caller has not already
// declared a $durationToSecs local.
const durationToSecsFn = `:local durationToSecs do={
  :local s $1
  :local slen [:len $s]
  :local total 0
  :local num ""
  :local i 0
  :while ($i < $slen) do={
    :local ch [:pick $s $i ($i + 1)]
    :if ($ch = "w") do={
      :set total ($total + [:tonum $num] * 604800)
      :set num ""
      :set i ($i + 1)
    } else={
      :if ($ch = "d") do={
        :set total ($total + [:tonum $num] * 86400)
        :set num ""
        :set i ($i + 1)
      } else={
        :if ($ch = "h") do={
          :set total ($total + [:tonum $num] * 3600)
          :set num ""
          :set i ($i + 1)
        } else={
          :if ($ch = "m") do={
            :set total ($total + [:tonum $num] * 60)
            :set num ""
            :set i ($i + 1)
          } else={
            :if ($ch = "s") do={
              :set total ($total + [:tonum $num])
              :set num ""
              :set i ($i + 1)
            } else={
              :if ($ch = ":") do={
                # Remainder from here on is HH:MM:SS — consume it whole.
                :local rest ($num . [:pick $s $i $slen])
                :local h [:tonum [:pick $rest 0 [:find $rest ":"]]]
                :local rest2 [:pick $rest ([:find $rest ":"] + 1) [:len $rest]]
                :local m [:tonum [:pick $rest2 0 [:find $rest2 ":"]]]
                :local sec [:tonum [:pick $rest2 ([:find $rest2 ":"] + 1) [:len $rest2]]]
                :set total ($total + $h * 3600 + $m * 60 + $sec)
                :set num ""
                :set i $slen
              } else={
                :set num ($num . $ch)
                :set i ($i + 1)
              }
            }
          }
        }
      }
    }
  }
  :return $total
}`

// cleanupScriptV7 uses [:tonsec [:timestamp]] available on RouterOS 7.12+.
// exp: comments store Unix epoch seconds written by the on-login profile script.
// Also removes any user whose uptime quota is exhausted (uptime >= limit-uptime),
// independent of the exp: comment — RouterOS already blocks their login at that
// point, so there's no reason to keep the account.
var cleanupScriptV7 = `:local nowEpoch ([:tonsec [:timestamp]] / 1000000000)
` + durationToSecsFn + `
:foreach u in=[/ip/hotspot/user/find] do={
  :local uname [:tostr [/ip/hotspot/user/get $u name]]
  :local comment [:tostr [/ip/hotspot/user/get $u comment]]
  :local removed false
  :if ([:len $comment] > 4 && [:pick $comment 0 4] = "exp:") do={
    :local expEpoch [:tonum [:pick $comment 4 [:len $comment]]]
    :if ($expEpoch > 0 && $expEpoch < $nowEpoch) do={
      :log info ("pikro-cleanup: removing $uname (exp: comment expired, exp=$expEpoch now=$nowEpoch)")
      /ip/hotspot/user/remove $u
      :set removed true
    }
  }
  :if (!$removed) do={
    :local limit [:tostr [/ip/hotspot/user/get $u limit-uptime]]
    :local used [:tostr [/ip/hotspot/user/get $u uptime]]
    :if ([:len $limit] > 0 && [:len $used] > 0) do={
      :if ([$durationToSecs $used] >= [$durationToSecs $limit]) do={
        :log info ("pikro-cleanup: removing $uname (uptime quota exhausted, uptime=$used limit=$limit)")
        /ip/hotspot/user/remove $u
      }
    }
  }
}`

// cleanupScriptV6 is the fallback for RouterOS 6 / early v7 without [:tonsec].
// exp: comments store "YYYY-MM-DD HH:MM:SS" (ISO, written by our on-login script).
// Uses Mikhmon-style dateint/timeint functions for proven v6 compatibility.
// Also removes any user whose uptime quota is exhausted (uptime >= limit-uptime),
// independent of the exp: comment — see cleanupScriptV7 for rationale.
var cleanupScriptV6 = `:local dateint do={
  :local days [:pick $d 8 10]
  :local month [:pick $d 5 7]
  :local year [:pick $d 0 4]
  :return [:tonum ("$year$month$days")]
}
:local timeint do={
  :local hours [:pick $t 0 2]
  :local minutes [:pick $t 3 5]
  :return ($hours * 60 + $minutes)
}
` + durationToSecsFn + `
:local date [/system clock get date]
:local time [/system clock get time]
:local today [$dateint d=$date]
:local curtime [$timeint t=$time]
:foreach u in=[/ip/hotspot/user/find] do={
  :local uname [:tostr [/ip/hotspot/user/get $u name]]
  :local comment [:tostr [/ip/hotspot/user/get $u comment]]
  :local removed false
  :if ([:len $comment] > 18 && [:pick $comment 0 4] = "exp:") do={
    :local s [:pick $comment 4 23]
    :if ([:pick $s 4 5] = "-" && [:pick $s 7 8] = "-" && [:pick $s 10 11] = " ") do={
      :local expd [$dateint d=$s]
      :local expt [$timeint t=[:pick $s 11 19]]
      :if (($expd < $today) || ($expd = $today && $expt < $curtime)) do={
        :log info ("pikro-cleanup: removing $uname (exp: comment expired, exp=$s now=$date $time)")
        /ip/hotspot/user/remove $u
        :set removed true
      }
    }
  }
  :if (!$removed) do={
    :local limit [:tostr [/ip/hotspot/user/get $u limit-uptime]]
    :local used [:tostr [/ip/hotspot/user/get $u uptime]]
    :if ([:len $limit] > 0 && [:len $used] > 0) do={
      :if ([$durationToSecs $used] >= [$durationToSecs $limit]) do={
        :log info ("pikro-cleanup: removing $uname (uptime quota exhausted, uptime=$used limit=$limit)")
        /ip/hotspot/user/remove $u
      }
    }
  }
}`

// cleanupScriptBody returns the appropriate script for the detected RouterOS version.
// rosVersion should be the string from /system/resource get version, e.g. "7.19.6 (stable)".
func cleanupScriptBody(rosVersion string) string {
	// Extract major.minor version number
	major, minor := 7, 0
	fmt.Sscanf(rosVersion, "%d.%d", &major, &minor)
	if major > 7 || (major == 7 && minor >= 12) {
		return cleanupScriptV7
	}
	return cleanupScriptV6
}

// CleanupSchedulerStatus returns whether the cleanup scheduler is installed and its interval.
func (c *Client) CleanupSchedulerStatus() (installed bool, interval string, err error) {
	conn, err := c.connect()
	if err != nil {
		return false, "", err
	}
	defer conn.Close()

	reply, err := conn.Run("/system/scheduler/print", "?name="+cleanupScriptName)
	if err != nil {
		return false, "", err
	}
	if len(reply.Re) == 0 {
		return false, "", nil
	}
	return true, reply.Re[0].Map["interval"], nil
}

// InstallCleanupScheduler creates or replaces the cleanup scheduler entry.
// It detects the RouterOS version and picks the appropriate script variant.
func (c *Client) InstallCleanupScheduler(interval string) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// Detect RouterOS version to pick the right script variant.
	rosVersion := "7.0"
	if r, e := conn.Run("/system/resource/print"); e == nil && len(r.Re) > 0 {
		rosVersion = r.Re[0].Map["version"]
	}

	// Remove existing entry if present.
	existing, err := conn.Run("/system/scheduler/print", "?name="+cleanupScriptName)
	if err != nil {
		return err
	}
	if len(existing.Re) > 0 {
		id := existing.Re[0].Map[".id"]
		if _, err := conn.Run("/system/scheduler/remove", "=.id="+id); err != nil {
			return err
		}
	}

	_, err = conn.Run("/system/scheduler/add",
		"=name="+cleanupScriptName,
		"=interval="+interval,
		"=on-event="+cleanupScriptBody(rosVersion),
		"=comment=Auto-cleanup expired hotspot users (by Pikro)",
	)
	return err
}

// RemoveCleanupScheduler removes the cleanup scheduler entry.
func (c *Client) RemoveCleanupScheduler() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	existing, err := conn.Run("/system/scheduler/print", "?name="+cleanupScriptName)
	if err != nil {
		return err
	}
	if len(existing.Re) == 0 {
		return nil // already gone
	}
	id := existing.Re[0].Map[".id"]
	_, err = conn.Run("/system/scheduler/remove", "=.id="+id)
	return err
}
