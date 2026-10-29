package router

import (
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
		{"custom alogin is uploaded", LoginPageContent{Alogin: "<p>custom</p>"}, "hotspot/alogin.html", "<p>custom</p>"},
		{"empty alogin falls back to the built-in page", LoginPageContent{}, "hotspot/alogin.html", string(assets.HotspotAloginHTML)},
		{"stylesheet is uploaded", LoginPageContent{Style: "body{}"}, "hotspot/css/style.css", "body{}"},
		{"no stylesheet leaves the router's one", LoginPageContent{}, "hotspot/css/style.css", ""},
		{"md5 script is uploaded", LoginPageContent{MD5: "function hexMD5(){}"}, "hotspot/md5.js", "function hexMD5(){}"},
		{"no md5 script leaves the router's one", LoginPageContent{}, "hotspot/md5.js", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := hotspotFiles(tt.content)
			if err != nil {
				t.Fatal(err)
			}
			got := ""
			for _, f := range files {
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
