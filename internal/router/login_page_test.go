package router

import (
	"strings"
	"testing"

	"github.com/digino/pikro/internal/assets"
)

func TestHotspotFiles(t *testing.T) {
	tests := []struct {
		name    string
		content LoginPageContent
		file    string
		want    string // "" means the file must not be uploaded
	}{
		{"alogin is uploaded", LoginPageContent{Alogin: "<p>custom</p>"}, "hotspot/alogin.html", "<p>custom</p>"},
		{"no alogin leaves the router's one", LoginPageContent{}, "hotspot/alogin.html", ""},
		{"no login page leaves the router's one", LoginPageContent{}, "hotspot/login.html", ""},
		{"error page is always the built-in one", LoginPageContent{}, "hotspot/error.html", string(assets.HotspotErrorHTML)},
		{"stylesheet is uploaded", LoginPageContent{Style: "body{}"}, "hotspot/css/style.css", "body{}"},
		{"no stylesheet leaves the router's one", LoginPageContent{}, "hotspot/css/style.css", ""},
		{"md5 script is uploaded", LoginPageContent{MD5: "function hexMD5(){}"}, "hotspot/md5.js", "function hexMD5(){}"},
		{"no md5 script leaves the router's one", LoginPageContent{}, "hotspot/md5.js", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ""
			for _, f := range hotspotFiles(tt.content) {
				if f.name == tt.file {
					got = f.content
				}
			}
			if got != tt.want {
				t.Errorf("%s = %q, want %q", tt.file, got, tt.want)
			}
		})
	}
}

func TestCheckPageSizes(t *testing.T) {
	page := func(n int) string { return strings.Repeat("a", n) }
	tests := []struct {
		name  string
		files []hotspotFile
		want  string // "" means accepted; otherwise part of the error
	}{
		{"all at the limit", []hotspotFile{{"hotspot/login.html", page(4095)}, {"hotspot/status.html", page(10)}}, ""},
		{"login page over", []hotspotFile{{"hotspot/login.html", page(4096)}}, "login.html is 4096 bytes"},
		{"a later page over", []hotspotFile{{"hotspot/login.html", page(10)}, {"hotspot/status.html", page(5000)}}, "status.html is 5000 bytes"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkPageSizes(tt.files)
			if tt.want == "" && err != nil || tt.want != "" && (err == nil || !strings.Contains(err.Error(), tt.want)) {
				t.Errorf("checkPageSizes: %v, want %q", err, tt.want)
			}
		})
	}
}
