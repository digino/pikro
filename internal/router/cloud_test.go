package router

import (
	"strings"
	"testing"
	"time"
)

func TestCloudRouterID(t *testing.T) {
	tests := []struct {
		name, source, want string
	}{
		{"installed script", ":local url \"https://cloud.test/sync\"\n:local token \"b9c9-44.SECRET\"\n", "b9c9-44"},
		{"no token line", ":local url \"https://cloud.test/sync\"\n", ""},
		{"token without secret", ":local token \"b9c9-44\"\n", ""},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cloudRouterID(tt.source); got != tt.want {
				t.Errorf("cloudRouterID = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCloudToken(t *testing.T) {
	tests := []struct{ name, source, want string }{
		{"installed script", ":local token \"b9c9-44.SECRET\"\n", "b9c9-44.SECRET"},
		{"no token line", ":local url \"https://cloud.test/sync\"\n", ""},
		{"token without secret", ":local token \"b9c9-44\"\n", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cloudToken(tt.source); got != tt.want {
				t.Errorf("cloudToken = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestPickHardwareID(t *testing.T) {
	tests := []struct {
		name             string
		routerboard, lic map[string]string
		want             string
	}{
		{"RouterBOARD serial", map[string]string{"routerboard": "true", "serial-number": "HXB0A1"}, map[string]string{"software-id": "AB12-CD34"}, "HXB0A1"},
		{"CHR v6: software-id", map[string]string{"routerboard": "false"}, map[string]string{"software-id": "AB12-CD34"}, "AB12-CD34"},
		{"CHR v7: system-id", map[string]string{"routerboard": "false"}, map[string]string{"system-id": "Zx9y"}, "Zx9y"},
		{"nothing", map[string]string{}, map[string]string{}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickHardwareID(tt.routerboard, tt.lic); got != tt.want {
				t.Errorf("pickHardwareID = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestCloudStartTime(t *testing.T) {
	tests := []struct {
		offset time.Duration
		want   string
	}{
		{0, "00:00:00"},
		{3*time.Minute + 27*time.Second, "00:03:27"},
		{9*time.Minute + 59*time.Second, "00:09:59"},
	}
	for _, tt := range tests {
		if got := cloudStartTime(tt.offset); got != tt.want {
			t.Errorf("cloudStartTime(%s) = %q, want %q", tt.offset, got, tt.want)
		}
	}
}

func TestWithCloudBuyButton(t *testing.T) {
	const buy = "https://cloud.test/r/abc-1?login=$(link-login-only-esc)"
	tests := []struct {
		name, html, buyURL string
		wantButton         bool
		wantTail           string // what must follow the snippet, "" to skip
	}{
		{"inserted before the closing body tag", "<html><body><form name=\"login\"></form></body></html>", buy, true, "</body></html>"},
		{"closing tag in any case", "<BODY>x</BODY>", buy, true, "</BODY>"},
		{"page without a body tag gets it appended", "<form name=\"login\"></form>", buy, true, ""},
		{"empty url adds nothing", "<body>x</body>", "", false, "<body>x</body>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withCloudBuyButton(tt.html, tt.buyURL)
			if has := strings.Contains(got, `href="`+buy+`"`); has != tt.wantButton {
				t.Errorf("button present = %v, want %v:\n%s", has, tt.wantButton, got)
			}
			if tt.wantButton && !strings.Contains(got, "Acheter avec mobile money") {
				t.Errorf("button label missing:\n%s", got)
			}
			if tt.wantTail != "" && !strings.HasSuffix(got, tt.wantTail) {
				t.Errorf("page = %q, want suffix %q", got, tt.wantTail)
			}
		})
	}
}

func TestWithCloudBuyButtonIsIdempotentAndRemovable(t *testing.T) {
	page := "<html><body>\n<form name=\"login\"></form>\n</body></html>"
	tests := []struct {
		name string
		do   func(string) string
		want string
	}{
		{"applying twice equals once", func(s string) string {
			return withCloudBuyButton(withCloudBuyButton(s, "https://c.test/r/1?login=$(x)"), "https://c.test/r/1?login=$(x)")
		}, withCloudBuyButton(page, "https://c.test/r/1?login=$(x)")},
		{"a new url replaces the old one", func(s string) string {
			return withCloudBuyButton(withCloudBuyButton(s, "https://old.test/r/9"), "https://c.test/r/1?login=$(x)")
		}, withCloudBuyButton(page, "https://c.test/r/1?login=$(x)")},
		{"empty url restores the original page", func(s string) string {
			return withCloudBuyButton(withCloudBuyButton(s, "https://c.test/r/1"), "")
		}, page},
		{"removing from a clean page changes nothing", func(s string) string { return withCloudBuyButton(s, "") }, page},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.do(page); got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestCloudBuyURL(t *testing.T) {
	tests := []struct{ name, base, id, want string }{
		{"plain", "https://cloud.test", "b9c9-44", "https://cloud.test/r/b9c9-44?login=$(link-login-only-esc)"},
		{"id is path-escaped", "https://cloud.test", "a/b", "https://cloud.test/r/a%2Fb?login=$(link-login-only-esc)"},
		{"no cloud url", "", "b9c9-44", ""},
		{"not connected", "https://cloud.test", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cloudBuyURL(tt.base, tt.id); got != tt.want {
				t.Errorf("cloudBuyURL = %q, want %q", got, tt.want)
			}
		})
	}
}
