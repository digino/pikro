package router

import (
	"fmt"
	"strings"

	"github.com/digino/pikro/internal/assets"
)

type hotspotFile struct {
	name    string
	content string
}

// hotspotFiles maps each page to its file. A page left empty is not uploaded:
// the router keeps its own. error.html is always the fixed built-in page.
func hotspotFiles(content LoginPageContent) []hotspotFile {
	var files []hotspotFile
	for _, f := range []hotspotFile{
		{"hotspot/login.html", content.Login},
		{"hotspot/logout.html", content.Logout},
		{"hotspot/status.html", content.Status},
		{"hotspot/alogin.html", content.Alogin},
		{"hotspot/error.html", string(assets.HotspotErrorHTML)},
		{"hotspot/redirect.html", content.Redirect},
		// Only pages that link them need these. Both already exist in a new
		// hotspot's folder, which matters on RouterOS 6: it can overwrite
		// files over the API but not create them.
		{"hotspot/css/style.css", content.Style},
		{"hotspot/md5.js", content.MD5},
	} {
		if f.content != "" {
			files = append(files, f)
		}
	}
	return files
}

// maxPageBytes is the largest file RouterOS 6 stores through the API; it
// answers "new contents too long" above it.
const maxPageBytes = 4095

// checkPageSizes refuses the upload, before anything is written, when a page
// is over maxPageBytes, so the hotspot is never left half updated.
func checkPageSizes(files []hotspotFile) error {
	for _, f := range files {
		if len(f.content) > maxPageBytes {
			return fmt.Errorf("%s is %d bytes, but the router accepts at most %d per page: shorten the title, subtitle or price list",
				strings.TrimPrefix(f.name, "hotspot/"), len(f.content), maxPageBytes)
		}
	}
	return nil
}

// LoginPageContent holds the fully-rendered HTML for each customizable
// hotspot page, built by the frontend (loginPageTemplates.ts), both at
// hotspot setup and from Settings > Login page.
type LoginPageContent struct {
	Login    string `json:"login"`
	Logout   string `json:"logout"`
	Status   string `json:"status"`
	Redirect string `json:"redirect"`
	Alogin   string `json:"alogin"` // shown right after a successful login
	Style    string `json:"css"`    // stylesheet the pages link, kept out of them for size
	MD5      string `json:"md5Js"`  // CHAP hash script the login page links

	// CloudURL is Pikro Cloud's base URL. When set and the router is connected
	// to Cloud, the login page gets the buy button; see withCloudBuyButton.
	CloudURL string `json:"-"`
}

// UploadLoginPage uploads the given rendered hotspot pages to the router's
// hotspot/ directory, plus the fixed error page.
func (c *Client) UploadLoginPage(profileName string, content LoginPageContent) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	// A page uploaded over a connected router keeps its buy button.
	if content.Login != "" {
		routerID := cloudRouterID(firstMap(conn, "/system/script/print", "?name="+cloudScriptName)["source"])
		content.Login = withCloudBuyButton(content.Login, cloudBuyURL(content.CloudURL, routerID))
	}
	files := hotspotFiles(content)
	if err := checkPageSizes(files); err != nil {
		return err
	}

	// List all files once to build a name → id/index map.
	idByName := map[string]string{}
	idxByName := map[string]int{}
	if listReply, listErr := conn.RunArgs([]string{"/file/print"}); listErr == nil {
		for i, re := range listReply.Re {
			name := re.Map["name"]
			idByName[name] = re.Map[".id"]
			idxByName[name] = i
		}
	}

	// Stop at the first file the router refuses: login.html goes first, so a
	// page RouterOS 6 finds too long (over 4 KB) leaves the hotspot as it was
	// instead of half-replaced.
	for _, f := range files {
		// Create file if it doesn't exist yet (ignore error — may already exist).
		conn.RunArgs([]string{"/file/add", "=name=" + f.name, "=type=.html"})

		fid := idByName[f.name]
		fidx, hasFidx := idxByName[f.name]

		// If not found in the initial listing, re-fetch after /file/add.
		if fid == "" && !hasFidx {
			if r2, e2 := conn.RunArgs([]string{"/file/print"}); e2 == nil {
				for i, re := range r2.Re {
					if re.Map["name"] == f.name {
						fid = re.Map[".id"]
						fidx = i
						hasFidx = true
						break
					}
				}
			}
		}

		var setErr error
		switch {
		case fid != "":
			_, setErr = conn.RunArgs([]string{"/file/set", "=.id=" + fid, "=contents=" + f.content})
		case hasFidx:
			_, setErr = conn.RunArgs([]string{"/file/set", fmt.Sprintf("=numbers=%d", fidx), "=contents=" + f.content})
		default:
			setErr = fmt.Errorf("%s not found after /file/add — ensure the hotspot is set up first", f.name)
		}
		if setErr != nil {
			return fmt.Errorf("file/set %s: %w", f.name, setErr)
		}
	}

	if profileName == "" {
		hsReply, hsErr := conn.RunArgs([]string{"/ip/hotspot/print"})
		if hsErr == nil && len(hsReply.Re) > 0 {
			profileName = hsReply.Re[0].Map["profile"]
		}
		if profileName == "" {
			profileName = "default" // RouterOS default profile name
		}
	}
	// numbers= expects a positional index into the last /print result, not a
	// name — resolve the profile's .id explicitly rather than passing the
	// name string as numbers= (which RouterOS rejects with "no such item").
	profReply, err := conn.RunArgs([]string{
		"/ip/hotspot/profile/print",
		"?name=" + profileName,
	})
	if err != nil {
		return err
	}
	if len(profReply.Re) == 0 {
		return fmt.Errorf("hotspot profile %q not found", profileName)
	}
	_, err = conn.RunArgs([]string{
		"/ip/hotspot/profile/set",
		"=.id=" + profReply.Re[0].Map[".id"],
		"=html-directory=hotspot",
	})
	return err
}

// GetLoginPageHTML reads back the contents of hotspot/login.html as currently
// stored on the router — i.e. the login page actually served to devices,
// independent of whatever template Pikro's local config thinks is selected.
func (c *Client) GetLoginPageHTML() (string, error) {
	conn, err := c.connect()
	if err != nil {
		return "", err
	}
	defer conn.Close()

	reply, err := conn.RunArgs([]string{
		"/file/print",
		"?name=hotspot/login.html",
		"=.proplist=contents",
	})
	if err != nil {
		return "", err
	}
	if len(reply.Re) == 0 {
		return "", fmt.Errorf("hotspot/login.html not found — run hotspot setup first")
	}
	return reply.Re[0].Map["contents"], nil
}
