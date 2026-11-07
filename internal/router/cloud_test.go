package router

import (
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
