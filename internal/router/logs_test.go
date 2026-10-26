package router

import "testing"

func TestIsHotspotLog(t *testing.T) {
	tests := []struct {
		name    string
		topics  string
		message string
		want    bool
	}{
		{"hotspot login", "hotspot,account,info,debug", "gbiu (10.5.50.2): logged in", true},
		{"hotspot failed login", "hotspot,info,debug", "a3fm (10.5.50.2): login failed: invalid password", true},
		{"pikro cleanup script", "script,info", "pikro-cleanup: removing kx7m (exp: comment expired)", true},
		{"pikro on-login script", "script,info", "pikro-on-login: stamping kx7m (exp=1792900000)", true},
		{"someone else's script", "script,info", "backup: done", false},
		{"pikro text outside script topic", "system,info", "pikro-cleanup: removing kx7m", false},
		{"unrelated system line", "system,info,account", "user admin logged in from 10.0.0.2 via api", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := map[string]any{"topics": tt.topics, "message": tt.message}
			if got := IsHotspotLog(e); got != tt.want {
				t.Errorf("IsHotspotLog(%q, %q) = %v, want %v", tt.topics, tt.message, got, tt.want)
			}
		})
	}
}
