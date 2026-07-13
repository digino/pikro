package router

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/digino/pikro/internal/assets"
)

// LoginPageParams holds the visual customisation options for the hotspot login page.
type LoginPageParams struct {
	Title       string // page heading / logo text, default "Sign in to continue"
	Subtitle    string // subtitle below the heading, default "$(hostname)"
	AccentColor string // CSS hex color for button + focus ring, default "#111827"
}

// renderLoginPage applies params to the embedded hotspot_login.html template.
// Uses text/template so user-supplied values are not HTML-escaped.
func renderLoginPage(p LoginPageParams) (string, error) {
	if p.Title == "" {
		p.Title = "Sign in to continue"
	}
	if p.Subtitle == "" {
		p.Subtitle = "$(hostname)"
	}
	if p.AccentColor == "" {
		p.AccentColor = "#111827"
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

// UploadLoginPage renders the login page template with the given params and
// uploads all hotspot pages (login, logout, status, alogin, error, redirect)
// to the router's hotspot/ directory.
func (c *Client) UploadLoginPage(profileName string, p LoginPageParams) error {
	html, err := renderLoginPage(p)
	if err != nil {
		return err
	}

	files := []hotspotFile{
		{"hotspot/login.html", html},
		{"hotspot/logout.html", string(assets.HotspotLogoutHTML)},
		{"hotspot/status.html", string(assets.HotspotStatusHTML)},
		{"hotspot/alogin.html", string(assets.HotspotAloginHTML)},
		{"hotspot/error.html", string(assets.HotspotErrorHTML)},
		{"hotspot/redirect.html", string(assets.HotspotRedirectHTML)},
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
		profileName = "pikro-profile"
	}
	_, err = conn.RunArgs([]string{
		"/ip/hotspot/profile/set",
		"=numbers=" + profileName,
		"=html-directory=hotspot",
	})
	return err
}
