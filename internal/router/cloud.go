package router

import (
	_ "embed"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
	"strings"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

// Pikro Cloud: the router keeps a pool of unsold vouchers that Cloud issued
// and syncs them every 10 minutes with the pikro-cloud-sync script. The
// script and its token come from Cloud; Pikro installs them, then the router
// works on its own. See pikro-cloud's docs/cloud-setup.md.
const (
	cloudScriptName = "pikro-cloud-sync"
	cloudComment    = "pikro-cloud" // on pool vouchers and walled-garden entries
	cloudInterval   = 10 * time.Minute
	cloudCertName   = "ISRG Root X1"
)

// isrgRootX1 is Let's Encrypt's root, which Cloud's certificate chains to.
// RouterOS 6 has no trust store, so fetch with check-certificate=yes needs it
// imported. SHA-256 96:BC:EC:06:26:49:76:F3:74:60:77:9A:CF:28:C5:A7:CF:E8:A3:C0:AA:E1:1A:8F:FC:EE:05:C0:BD:DF:08:C6.
//
//go:embed isrgrootx1.pem
var isrgRootX1 string

var cloudTokenRe = regexp.MustCompile(`:local token "([^".]+)\.[^"]+"`)

// cloudRouterID reads the Cloud router id from an installed sync script.
func cloudRouterID(source string) string {
	if m := cloudTokenRe.FindStringSubmatch(source); m != nil {
		return m[1]
	}
	return ""
}

// pickHardwareID is what identifies a router to Cloud across re-setups: the
// RouterBOARD serial, or on a CHR (no RouterBOARD) its licence id.
func pickHardwareID(routerboard, license map[string]string) string {
	for _, id := range []string{routerboard["serial-number"], license["software-id"], license["system-id"]} {
		if id != "" {
			return id
		}
	}
	return ""
}

// cloudStartTime spreads routers' syncs over the interval, so a hundred
// routers don't all reach Cloud on the same minute.
func cloudStartTime(offset time.Duration) string {
	s := int(offset.Seconds())
	return fmt.Sprintf("%02d:%02d:%02d", s/3600, s/60%60, s%60)
}

type CloudStatus struct {
	Connected       bool   `json:"connected"`     // sync script installed
	CloudRouterID   string `json:"cloudRouterId"` // from the script's token
	PoolVouchers    int    `json:"poolVouchers"`  // unsold Cloud vouchers on the router
	HardwareID      string `json:"hardwareId"`
	RouterOSVersion string `json:"routerOsVersion"`
	HasHotspot      bool   `json:"hasHotspot"`
}

func firstMap(conn *ros.Client, path string, args ...string) map[string]string {
	r, err := conn.RunArgs(append([]string{path}, args...))
	if err != nil || len(r.Re) == 0 {
		return map[string]string{}
	}
	return r.Re[0].Map
}

// CloudStatus reports whether the router is connected to Cloud, and what
// setup needs to send for it.
func (c *Client) CloudStatus() (*CloudStatus, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	s := &CloudStatus{
		HardwareID:      pickHardwareID(firstMap(conn, "/system/routerboard/print"), firstMap(conn, "/system/license/print")),
		RouterOSVersion: strings.Fields(c.rosVersion(conn) + " ")[0], // "6.49.23 (long-term)" -> "6.49.23"
		HasHotspot:      len(firstMap(conn, "/ip/hotspot/print")) > 0,
	}
	if script := firstMap(conn, "/system/script/print", "?name="+cloudScriptName); len(script) > 0 {
		s.Connected = true
		s.CloudRouterID = cloudRouterID(script["source"])
	}
	s.PoolVouchers, err = poolVouchers(conn)
	return s, err
}

func poolVouchers(conn *ros.Client) (int, error) {
	r, err := conn.RunArgs([]string{"/ip/hotspot/user/print", "?comment=" + cloudComment, "=count-only="})
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(r.Done.Map["ret"])
}

// CloudImportCert imports Let's Encrypt's root if the router lacks it.
func (c *Client) CloudImportCert() SetupStepResult {
	name := "Install the security certificate"
	fail := func(err error) SetupStepResult { return SetupStepResult{Name: name, Error: err.Error()} }
	conn, err := c.connect()
	if err != nil {
		return fail(err)
	}
	defer conn.Close()

	if len(firstMap(conn, "/certificate/print", "?common-name="+cloudCertName)) > 0 {
		return SetupStepResult{Name: name, OK: true, Skipped: true}
	}
	f, err := createFile(conn, "pikro-isrgrootx1")
	if err != nil {
		return fail(err)
	}
	if _, err := conn.RunArgs([]string{"/file/set", "=.id=" + f[".id"], "=contents=" + isrgRootX1}); err != nil {
		return fail(err)
	}
	_, err = conn.RunArgs([]string{"/certificate/import", "=file-name=" + f["name"], "=passphrase="})
	conn.RunArgs([]string{"/file/remove", "=.id=" + f[".id"]})
	if err != nil {
		return fail(err)
	}
	return SetupStepResult{Name: name, OK: true}
}

// createFile makes an empty file and returns its /file entry. RouterOS 7 has
// /file/add (base.pem); RouterOS 6 doesn't, but "/file print file=base"
// creates base.txt, a moment later.
func createFile(conn *ros.Client, base string) (map[string]string, error) {
	if _, err := conn.RunArgs([]string{"/file/add", "=name=" + base + ".pem"}); err == nil {
		if f := firstMap(conn, "/file/print", "?name="+base+".pem"); f[".id"] != "" {
			return f, nil
		}
	}
	conn.RunArgs([]string{"/file/print", "=file=" + base})
	for range 10 {
		if f := firstMap(conn, "/file/print", "?name="+base+".txt"); f[".id"] != "" {
			return f, nil
		}
		time.Sleep(300 * time.Millisecond)
	}
	return nil, fmt.Errorf("could not create %s on the router", base)
}

// removeAll removes every item at path matching query (e.g. "?comment=x").
func removeAll(conn *ros.Client, path, query string) error {
	r, err := conn.RunArgs([]string{path + "/print", query, "=.proplist=.id"})
	if err != nil {
		return err
	}
	for _, re := range r.Re {
		if _, err := conn.RunArgs([]string{path + "/remove", "=.id=" + re.Map[".id"]}); err != nil {
			return err
		}
	}
	return nil
}

// CloudInstall installs (or replaces) the sync script from Cloud with its
// 10-minute scheduler, and lets hotspot guests reach cloudHost before login.
func (c *Client) CloudInstall(script, cloudHost string) SetupStepResult {
	name := "Install the sync script"
	fail := func(err error) SetupStepResult { return SetupStepResult{Name: name, Error: err.Error()} }
	conn, err := c.connect()
	if err != nil {
		return fail(err)
	}
	defer conn.Close()

	for _, path := range []string{"/system/scheduler", "/system/script"} {
		if err := removeAll(conn, path, "?name="+cloudScriptName); err != nil {
			return fail(err)
		}
	}
	if err := removeAll(conn, "/ip/hotspot/walled-garden", "?comment="+cloudComment); err != nil {
		return fail(err)
	}
	steps := [][]string{
		{"/system/script/add", "=name=" + cloudScriptName, "=source=" + script,
			"=comment=Pikro Cloud voucher sync (by Pikro)"},
		{"/system/scheduler/add", "=name=" + cloudScriptName, "=interval=10m",
			"=start-time=" + cloudStartTime(rand.N(cloudInterval)),
			"=on-event=/system script run " + cloudScriptName,
			"=comment=Pikro Cloud voucher sync every 10 minutes (by Pikro)"},
		// The buy page lives on Cloud: guests open it before they have a voucher.
		{"/ip/hotspot/walled-garden/add", "=dst-host=" + cloudHost, "=comment=" + cloudComment},
	}
	for _, args := range steps {
		if _, err := conn.RunArgs(args); err != nil {
			return fail(fmt.Errorf("%s: %w", args[0], err))
		}
	}
	return SetupStepResult{Name: name, OK: true}
}

// CloudFirstSync runs the sync once, so the pool is filled before the
// technician leaves, and reports how many vouchers arrived.
func (c *Client) CloudFirstSync() (SetupStepResult, int) {
	name := "First sync with Pikro Cloud"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}, 0
	}
	defer conn.Close()

	if _, err := conn.RunArgs([]string{"/system/script/run", "=number=" + cloudScriptName}); err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}, 0
	}
	n, err := poolVouchers(conn)
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}, 0
	}
	if n == 0 {
		return SetupStepResult{Name: name, Error: "no vouchers arrived: the router could not reach Pikro Cloud"}, 0
	}
	return SetupStepResult{Name: name, OK: true}, n
}

// CloudRemove takes the router off Cloud: scheduler, script, walled-garden
// entry, the sync's temp file and the unsold pool vouchers. The certificate
// stays; it is harmless and other HTTPS fetches use it.
func (c *Client) CloudRemove() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, rm := range [][2]string{
		{"/system/scheduler", "?name=" + cloudScriptName},
		{"/system/script", "?name=" + cloudScriptName},
		{"/ip/hotspot/walled-garden", "?comment=" + cloudComment},
		{"/ip/hotspot/user", "?comment=" + cloudComment},
		{"/file", "?name=" + cloudScriptName + ".txt"},
	} {
		if err := removeAll(conn, rm[0], rm[1]); err != nil {
			return fmt.Errorf("%s: %w", rm[0], err)
		}
	}
	return nil
}
