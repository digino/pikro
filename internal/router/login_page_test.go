package router

import (
	"testing"

	"github.com/digino/pikro/internal/assets"
)

func TestHotspotFilesAlogin(t *testing.T) {
	tests := []struct {
		name    string
		content LoginPageContent
		want    string
	}{
		{"custom page is uploaded", LoginPageContent{Alogin: "<p>custom</p>"}, "<p>custom</p>"},
		{"empty falls back to the built-in page", LoginPageContent{}, string(assets.HotspotAloginHTML)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files, err := hotspotFiles(tt.content)
			if err != nil {
				t.Fatal(err)
			}
			for _, f := range files {
				if f.name == "hotspot/alogin.html" {
					if f.content != tt.want {
						t.Errorf("alogin.html = %q, want %q", f.content, tt.want)
					}
					return
				}
			}
			t.Error("alogin.html not in the uploaded files")
		})
	}
}
