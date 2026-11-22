package router

import (
	_ "embed"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

// Pikro Cloud
const (
	cloudScriptName = "pikro-cloud-sync"
	CloudComment    = "pikro-cloud" // on unsold pool vouchers and walled-garden entries
	cloudInterval   = 10 * time.Minute
	cloudCertName   = "ISRG Root X1"
)

// isrgRootX1 is Let's Encrypt's root, which Cloud's certificate chains to.
// RouterOS 6 has no trust store, so fetch with check-certificate=yes needs it
// imported. SHA-256 96:BC:EC:06:26:49:76:F3:74:60:77:9A:CF:28:C5:A7:CF:E8:A3:C0:AA:E1:1A:8F:FC:EE:05:C0:BD:DF:08:C6.
//
//go:embed isrgrootx1.pem
var isrgRootX1 string

var cloudTokenRe = regexp.MustCompile(`:local token "([^".]+\.[^"]+)"`)

// cloudToken reads "<router-id>.<secret>" from an installed sync script.
func cloudToken(source string) string {
	if m := cloudTokenRe.FindStringSubmatch(source); m != nil {
		return m[1]
	}
	return ""
}

// cloudRouterID reads the Cloud router id from an installed sync script.
func cloudRouterID(source string) string {
	id, _, _ := strings.Cut(cloudToken(source), ".")
	return id
}

// CloudToken is the router's Cloud token, from its sync script; "" when the
// script isn't installed. It authenticates the router to Cloud: keep it out
// of the browser.
func (c *Client) CloudToken() (string, error) {
	conn, err := c.connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()
	return cloudToken(firstMap(conn, "/system/script/print", "?name="+cloudScriptName)["source"]), nil
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

// cloudStartTime spreads routers' syncs over the interval.
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
	BuyURL          string `json:"buyUrl"`    // guests open it before logging in; "" when not connected
	BuyButton       string `json:"buyButton"` // login page: "present", "missing" or "unknown" (not readable completely)
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
func (c *Client) CloudStatus(cloudURL string) (*CloudStatus, error) {
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
		s.BuyURL = strings.TrimSuffix(cloudBuyURL(cloudURL, s.CloudRouterID), cloudLoginParam)
		s.BuyButton = buyButtonState(conn)
	}
	s.PoolVouchers, err = poolVouchers(conn)
	return s, err
}

func buyButtonState(conn *ros.Client) string {
	_, html, ok := readLoginPage(conn)
	switch {
	case !ok:
		return "unknown"
	case strings.Contains(html, cloudLoginStart):
		return "present"
	}
	return "missing"
}

func poolVouchers(conn *ros.Client) (int, error) {
	r, err := conn.RunArgs([]string{"/ip/hotspot/user/print", "?comment=" + CloudComment, "=count-only="})
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
	if err := removeAll(conn, "/ip/hotspot/walled-garden", "?comment="+CloudComment); err != nil {
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
		{"/ip/hotspot/walled-garden/add", "=dst-host=" + cloudHost, "=comment=" + CloudComment},
	}
	for _, args := range steps {
		if _, err := conn.RunArgs(args); err != nil {
			return fail(fmt.Errorf("%s: %w", args[0], err))
		}
	}
	return SetupStepResult{Name: name, OK: true}
}

// The buy button and voucher auto-login that a connected router's login page
// carries, between these markers so they can be replaced or removed. Works
// with every template: a "code" page has login.code, the others
// login.username and login.password. It submits through the form's own
// onsubmit, so CHAP hashing still applies. Skipped when the page is showing
// a login error, so a wrong voucher can't loop.
const (
	cloudLoginStart = "<!-- pikro-cloud -->"
	cloudLoginEnd   = "<!-- /pikro-cloud -->"
	cloudLoginHTML  = cloudLoginStart + `<a href="%s" style="display:table;margin:16px auto;padding:12px 18px;background:#ffbf00;color:#111;border:2px solid #111;box-shadow:3px 3px 0 #111;font:700 15px sans-serif;text-decoration:none;white-space:nowrap">Acheter ticket par Flooz / T-Money</a>` +
		`<script>!function(){var m=/[?&]voucher=([^&#]+)/.exec(location),f=document.login;` +
		`if($(if error)1$(else)0$(endif)||!m||!f)return;m=decodeURIComponent(m[1].replace(/\+/g," "));` +
		`f.code?f.code.value=m:f.username.value=f.password.value=m;` +
		`f.onsubmit&&f.onsubmit()==0||f.submit()}()</script>` + cloudLoginEnd

	// A page that holds this comment gets the snippet right after it, wherever
	// the owner wants the button. Removal leaves the marker, so the page comes
	// back byte for byte.
	cloudBuyMarker = "<!-- pikro-buy -->"
)

var (
	cloudLoginRe = regexp.MustCompile(`(?s)` + regexp.QuoteMeta(cloudLoginStart) + `.*?` + regexp.QuoteMeta(cloudLoginEnd) + `\n?`)
	bodyEndRe    = regexp.MustCompile(`(?i)</body>`)
	// The login form is the page's last form, both in MikroTik's default page
	// and in Pikro's. Right after it the button sits in the form's own column,
	// where before </body> a flex body would put it beside the page.
	formEndRe = regexp.MustCompile(`(?i)</form>`)
)

// cloudLoginParam is how the router tells the buy page where to log in: the
// hotspot's own address, which works with or without a DNS name. $(mac) is
// the guest's device, so the buy page can find its tickets again.
const cloudLoginParam = "?login=http://$(server-address)/login&mac=$(mac)"

// cloudBuyURL is the Cloud buy page for a router; the router fills in its own
// login URL. "" when Cloud's URL or the router id is unknown.
func cloudBuyURL(cloudURL, routerID string) string {
	if cloudURL == "" || routerID == "" {
		return ""
	}
	return cloudURL + "/r/" + url.PathEscape(routerID) + cloudLoginParam
}

// withCloudBuyButton returns the login page with the Cloud buy button and
// auto-login for buyURL, replacing any earlier one; with "" it returns the
// page without them.
func withCloudBuyButton(html, buyURL string) string {
	html = cloudLoginRe.ReplaceAllString(html, "")
	if buyURL == "" {
		return html
	}
	snippet := fmt.Sprintf(cloudLoginHTML, buyURL) + "\n"
	if at := strings.Index(html, cloudBuyMarker); at >= 0 {
		at += len(cloudBuyMarker)
		return html[:at] + snippet + html[at:]
	}
	if loc := formEndRe.FindAllStringIndex(html, -1); len(loc) > 0 {
		at := loc[len(loc)-1][1]
		return html[:at] + snippet + html[at:]
	}
	if loc := bodyEndRe.FindAllStringIndex(html, -1); len(loc) > 0 {
		at := loc[len(loc)-1][0]
		return html[:at] + snippet + html[at:]
	}
	return html + snippet
}

// loginPageChunk is how much /file/read returns per call.
const loginPageChunk = 4000

// completePage returns a file's whole content, or false when it can't be had.
// RouterOS only prints a file's contents up to 4095 bytes; for bigger files
// (RouterOS 7) read gets the rest in chunks, and RouterOS 6 has no read at all.
func completePage(size int, printed string, read func(offset, n int) (string, error)) (string, bool) {
	if len(printed) == size {
		return printed, true
	}
	var page strings.Builder
	for page.Len() < size {
		chunk, err := read(page.Len(), min(loginPageChunk, size-page.Len()))
		if err != nil || chunk == "" {
			return "", false
		}
		page.WriteString(chunk)
	}
	if page.Len() != size {
		return "", false
	}
	return page.String(), true
}

// loginPagePath is the login page the hotspot actually serves: login.html in
// the HTML directory of the first hotspot server's profile, the profile
// UploadLoginPage writes to. Owners often keep their own template folder.
func loginPagePath(conn *ros.Client) string {
	profile := firstMap(conn, "/ip/hotspot/print")["profile"]
	if profile == "" {
		profile = "default"
	}
	return loginPageFile(firstMap(conn, "/ip/hotspot/profile/print", "?name="+profile)["html-directory"])
}

// loginPageFile is login.html in an HTML directory as RouterOS reports it.
func loginPageFile(dir string) string {
	dir = strings.Trim(dir, "/")
	if dir == "" {
		dir = "hotspot"
	}
	return dir + "/login.html"
}

// readLoginPage returns the login page's /file id ("" when there is none) and
// its content; ok is false when the content can't be read completely.
func readLoginPage(conn *ros.Client) (id, html string, ok bool) {
	file := loginPagePath(conn)
	page := firstMap(conn, "/file/print", "?name="+file, "=.proplist=.id,size,contents")
	size, err := strconv.Atoi(page["size"])
	if page[".id"] == "" || err != nil {
		return page[".id"], "", false
	}
	html, ok = completePage(size, page["contents"], func(offset, n int) (string, error) {
		r, err := conn.RunArgs([]string{"/file/read", "=file=" + file, "=offset=" + strconv.Itoa(offset), "=chunk-size=" + strconv.Itoa(n)})
		if err != nil || len(r.Re) == 0 {
			return "", errors.New("cannot read")
		}
		return r.Re[0].Map["data"], nil
	})
	return page[".id"], html, ok
}

// CloudLoginPage puts the Cloud buy button on the router's login page when
// cloudURL is set, and takes it off when it is "". The owner's page is never
// written unless it was read completely; when the button can't be added the
// step is OK with a Warning and the page stays as it is.
func (c *Client) CloudLoginPage(cloudURL string) SetupStepResult {
	name := "Add the buy button to the login page"
	fail := func(err error) SetupStepResult { return SetupStepResult{Name: name, Error: err.Error()} }
	conn, err := c.connect()
	if err != nil {
		return fail(err)
	}
	defer conn.Close()

	routerID := cloudRouterID(firstMap(conn, "/system/script/print", "?name="+cloudScriptName)["source"])
	id, current, ok := readLoginPage(conn)
	if id == "" {
		return fail(fmt.Errorf("%s not found: run hotspot setup first", loginPagePath(conn)))
	}
	if !ok {
		return SetupStepResult{Name: name, OK: true, Warning: "The login page is a custom page that Pikro can't read completely on this router, so it was left as it is, without the buy button."}
	}
	html := withCloudBuyButton(current, cloudBuyURL(cloudURL, routerID))
	if html == current {
		return SetupStepResult{Name: name, OK: true, Skipped: true}
	}
	if _, err := conn.RunArgs([]string{"/file/set", "=.id=" + id, "=contents=" + html}); err != nil {
		if strings.Contains(err.Error(), "too long") {
			return SetupStepResult{Name: name, OK: true, Warning: "The login page is too large for this router to hold the buy button, so it was left as it is."}
		}
		return fail(err)
	}
	return SetupStepResult{Name: name, OK: true}
}

// firstSyncTries and firstSyncWait bound the first sync: right after setup
// it can fail for a moment (the clock just set by NTP, the certificate just
// imported) and work a few seconds later.
const (
	firstSyncTries = 3
	firstSyncWait  = 5 * time.Second
)

// CloudFirstSync runs the sync, retrying a little, so the pool is filled
// before the technician leaves, and reports how many vouchers arrived.
func (c *Client) CloudFirstSync() (SetupStepResult, int) {
	name := "First sync with Pikro Cloud"
	conn, err := c.connect()
	if err != nil {
		return SetupStepResult{Name: name, Error: err.Error()}, 0
	}
	defer conn.Close()

	for try := 1; ; try++ {
		if _, err := conn.RunArgs([]string{"/system/script/run", "=number=" + cloudScriptName}); err != nil {
			return SetupStepResult{Name: name, Error: err.Error()}, 0
		}
		n, err := poolVouchers(conn)
		if err != nil {
			return SetupStepResult{Name: name, Error: err.Error()}, 0
		}
		if n > 0 {
			return SetupStepResult{Name: name, OK: true}, n
		}
		if try == firstSyncTries {
			return SetupStepResult{Name: name, Error: "no vouchers arrived yet: the router retries by itself every 10 minutes. If they still don't come, check its log for \"pikro-cloud: sync failed\""}, 0
		}
		time.Sleep(firstSyncWait)
	}
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
		{"/ip/hotspot/walled-garden", "?comment=" + CloudComment},
		{"/ip/hotspot/user", "?comment=" + CloudComment},
		{"/file", "?name=" + cloudScriptName + ".txt"},
	} {
		if err := removeAll(conn, rm[0], rm[1]); err != nil {
			return fmt.Errorf("%s: %w", rm[0], err)
		}
	}
	return nil
}

// IsCloudVoucher reports whether a hotspot user is an unsold Cloud pool voucher.
func (c *Client) IsCloudVoucher(id string) (bool, error) {
	conn, err := c.connect()
	if err != nil {
		return false, err
	}
	defer conn.Close()
	r, err := conn.RunArgs([]string{"/ip/hotspot/user/print", "?.id=" + id, "=.proplist=comment"})
	if err != nil {
		return false, err
	}
	return len(r.Re) > 0 && r.Re[0].Map["comment"] == CloudComment, nil
}
