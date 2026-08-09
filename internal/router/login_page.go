package router

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/digino/pikro/internal/assets"
)

// defaultLoginPageParams holds the visual customisation options for the
// built-in fallback login page used during initial hotspot setup, before
// the admin has picked a template in Settings.
type defaultLoginPageParams struct {
	Title       string // page heading / logo text, default "Sign in to continue"
	Subtitle    string // subtitle below the heading, default "$(hostname)"
	AccentColor string // CSS hex color for button + focus ring, default "#111827"
}

// renderDefaultLoginPage applies params to the embedded hotspot_login.html
// template. Uses text/template so user-supplied values are not HTML-escaped.
func renderDefaultLoginPage(p defaultLoginPageParams) (string, error) {
	if p.Title == "" {
		p.Title = "Sign in to continue"
	}
	if p.Subtitle == "" {
		p.Subtitle = "$(hostname)"
	}
	if p.AccentColor == "" {
		p.AccentColor = "#ffbf00"
	}

	tmpl, err := template.New("login").Parse(string(assets.HotspotLoginHTML))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, p); err != nil {
		return "", err
	}
	return buf.String(), nil
}

type hotspotFile struct {
	name    string
	content string
}

// LoginPageContent holds the fully-rendered HTML for each customizable
// hotspot page (built by the frontend from one of its template presets).
// Any field left empty falls back to the built-in fixed default for that
// page — this is the path taken during initial hotspot setup, before an
// admin has picked a template in Settings, and also covers alogin/error,
// which are never customizable and always use their fixed default.
type LoginPageContent struct {
	Login    string
	Logout   string
	Status   string
	Redirect string
}

// UploadLoginPage uploads the given (or default, if empty) rendered hotspot
// pages — login, logout, status, alogin, error, redirect — to the router's
// hotspot/ directory. alogin and error are never customizable and always use
// their fixed built-in design.
func (c *Client) UploadLoginPage(profileName string, content LoginPageContent) error {
	if content.Login == "" {
		rendered, err := renderDefaultLoginPage(defaultLoginPageParams{})
		if err != nil {
			return err
		}
		content.Login = rendered
	}
	if content.Logout == "" {
		content.Logout = string(assets.HotspotLogoutHTML)
	}
	if content.Status == "" {
		content.Status = string(assets.HotspotStatusHTML)
	}
	if content.Redirect == "" {
		content.Redirect = string(assets.HotspotRedirectHTML)
	}

	files := []hotspotFile{
		{"hotspot/login.html", content.Login},
		{"hotspot/logout.html", content.Logout},
		{"hotspot/status.html", content.Status},
		{"hotspot/alogin.html", string(assets.HotspotAloginHTML)},
		{"hotspot/error.html", string(assets.HotspotErrorHTML)},
		{"hotspot/redirect.html", content.Redirect},
	}

	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

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

	var firstErr error
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
		if setErr != nil && firstErr == nil {
			firstErr = fmt.Errorf("file/set %s: %w", f.name, setErr)
		}
	}
	if firstErr != nil {
		return firstErr
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
