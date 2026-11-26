package router

import (
	"errors"
	"fmt"
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
	const buy = "https://cloud.test/r/abc-1?login=http://$(server-address)/login&mac=$(mac)"
	tests := []struct {
		name, html, cloudURL string
		wantButton           bool
		wantTail             string // what must follow the snippet, "" to skip
	}{
		{"inserted before the closing body tag", "<html><body><form name=\"login\"></form></body></html>", "https://cloud.test", true, "</body></html>"},
		{"closing tag in any case", "<BODY>x</BODY>", "https://cloud.test", true, "</BODY>"},
		{"page without a body tag gets it appended", "<form name=\"login\"></form>", "https://cloud.test", true, ""},
		{"empty url adds nothing", "<body>x</body>", "", false, "<body>x</body>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withCloudBuyButton(tt.html, tt.cloudURL, "abc-1")
			// The plain link that the buttons script, from Cloud, styles.
			for _, part := range []string{`<a id="pikro-buy" href="` + buy + `">Acheter ticket par Flooz / T-Money</a>`, `<script src="https://cloud.test/buttons.js" async></script>`} {
				if has := strings.Contains(got, part); has != tt.wantButton {
					t.Errorf("%s present = %v, want %v:\n%s", part, has, tt.wantButton, got)
				}
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
			return withCloudBuyButton(withCloudBuyButton(s, "https://c.test", "1"), "https://c.test", "1")
		}, withCloudBuyButton(page, "https://c.test", "1")},
		{"a new url replaces the old one", func(s string) string {
			return withCloudBuyButton(withCloudBuyButton(s, "https://old.test", "9"), "https://c.test", "1")
		}, withCloudBuyButton(page, "https://c.test", "1")},
		{"empty url restores the original page", func(s string) string {
			return withCloudBuyButton(withCloudBuyButton(s, "https://c.test", "1"), "", "1")
		}, page},
		{"removing from a clean page changes nothing", func(s string) string { return withCloudBuyButton(s, "", "1") }, page},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.do(page); got != tt.want {
				t.Errorf("got:\n%s\nwant:\n%s", got, tt.want)
			}
		})
	}
}

func TestLoginPageFile(t *testing.T) {
	tests := []struct{ name, dir, want string }{
		{"default folder", "hotspot", "hotspot/login.html"},
		{"owner's own template", "template-oklm", "template-oklm/login.html"},
		{"on flash", "flash/hotspot", "flash/hotspot/login.html"},
		{"leading slash", "/flash/hotspot", "flash/hotspot/login.html"},
		{"unset", "", "hotspot/login.html"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := loginPageFile(tt.dir); got != tt.want {
				t.Errorf("loginPageFile(%q) = %q, want %q", tt.dir, got, tt.want)
			}
		})
	}
}

func TestCloudBuyURL(t *testing.T) {
	tests := []struct{ name, base, id, want string }{
		{"plain", "https://cloud.test", "b9c9-44", "https://cloud.test/r/b9c9-44?login=http://$(server-address)/login&mac=$(mac)"},
		{"id is path-escaped", "https://cloud.test", "a/b", "https://cloud.test/r/a%2Fb?login=http://$(server-address)/login&mac=$(mac)"},
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

func TestCompletePage(t *testing.T) {
	page := strings.Repeat("0123456789", 900) // 9000 bytes
	reader := func(src string) func(int, int) (string, error) {
		return func(offset, n int) (string, error) { return src[offset : offset+n], nil }
	}
	tests := []struct {
		name    string
		size    int
		printed string
		read    func(int, int) (string, error)
		want    string
		wantOK  bool
	}{
		{"fits in print", 5, "hello", func(int, int) (string, error) { t.Fatal("read called"); return "", nil }, "hello", true},
		{"empty page", 0, "", nil, "", true},
		{"needs chunks (v7)", len(page), "", reader(page), page, true},
		{"chunk reader errors (v6)", len(page), "", func(int, int) (string, error) { return "", errors.New("no such command") }, "", false},
		{"empty chunk early", len(page), "", func(offset, n int) (string, error) {
			if offset > 0 {
				return "", nil
			}
			return page[:n], nil
		}, "", false},
		{"chunks longer than size", 10, "", func(int, int) (string, error) { return "0123456789abc", nil }, "", false},
		{"print is partial and read is wrong", 6, "abc", func(int, int) (string, error) { return "xyzw", nil }, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := completePage(tt.size, tt.printed, tt.read)
			if ok != tt.wantOK || got != tt.want {
				t.Errorf("completePage = (%d bytes, %v), want (%d bytes, %v)", len(got), ok, len(tt.want), tt.wantOK)
			}
		})
	}
}

func TestCompletePageReadsInChunksOf4000(t *testing.T) {
	var calls []string
	page := strings.Repeat("x", 9000)
	completePage(len(page), "", func(offset, n int) (string, error) {
		calls = append(calls, fmt.Sprintf("%d+%d", offset, n))
		return page[offset : offset+n], nil
	})
	if got, want := strings.Join(calls, " "), "0+4000 4000+4000 8000+1000"; got != want {
		t.Errorf("reads = %s, want %s", got, want)
	}
}

func TestCloudSnippetIsSmall(t *testing.T) {
	// RouterOS 6 refuses a login page over 4095 bytes. Pikro's largest page
	// (about 3.2 KB) plus this and its urls (about 130 bytes) must still fit.
	if n := len(fmt.Sprintf(cloudLoginHTML, "", "")); n > 450 {
		t.Errorf("snippet without its urls is %d bytes, want <= 450", n)
	}
}

func TestCloudButtonPlacementAndRemoval(t *testing.T) {
	tests := []struct {
		name, page string
		after      string // what must directly precede the snippet
	}{
		{"marker", "<html><body>a<!-- pikro-buy -->b</body></html>", "<!-- pikro-buy -->"},
		{"marker wins over the closing body tag", "</body><!-- pikro-buy --></body>", "<!-- pikro-buy -->"},
		{"marker wins over the form", "<form></form><!-- pikro-buy --></body>", "<!-- pikro-buy -->"},
		{"after the last form, inside its container", "<body><main><form name=sendin></form><form name=login></form><p>hint</p></main></body>", "<form name=login></form>"},
		{"form closing tag in any case", "<body><div><FORM></FORM></div></body>", "<FORM></FORM>"},
		{"no form: before the last body tag", "<html><body>a</body></html>", "a"},
		{"no form, no body tag", "<p>a</p>", "<p>a</p>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			added := withCloudBuyButton(tt.page, "https://c.test", "1")
			if !strings.Contains(added, tt.after+cloudLoginStart) {
				t.Errorf("snippet not right after %q:\n%s", tt.after, added)
			}
			if got := withCloudBuyButton(added, "", "1"); got != tt.page {
				t.Errorf("removal = %q, want the original %q", got, tt.page)
			}
			if again := withCloudBuyButton(added, "https://c.test", "1"); again != added {
				t.Errorf("re-adding changed the page:\n%s", again)
			}
		})
	}
}

func TestDeviceModeBlocksFetch(t *testing.T) {
	tests := []struct {
		name string
		mode map[string]string
		want bool
	}{
		{"7.17+ home mode, fetch off", map[string]string{"mode": "home", "fetch": "false", "scheduler": "true"}, true},
		{"fetch allowed", map[string]string{"mode": "home", "fetch": "true"}, false},
		{"device-mode without a fetch setting (before 7.17)", map[string]string{"mode": "enterprise"}, false},
		{"no device-mode at all (RouterOS 6)", nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deviceModeBlocksFetch(tt.mode); got != tt.want {
				t.Errorf("deviceModeBlocksFetch(%v) = %v, want %v", tt.mode, got, tt.want)
			}
		})
	}
}
