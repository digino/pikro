package update

import "testing"

func TestIsNewer(t *testing.T) {
	tests := []struct {
		name            string
		latest, current string
		want            bool
	}{
		{"newer patch", "v0.4.1", "v0.4.0", true},
		{"newer minor beats a bigger patch", "v0.5.0", "v0.4.9", true},
		{"same version", "v0.4.0", "v0.4.0", false},
		{"older release", "v0.3.0", "v0.4.0", false},
		{"the vv0.4.0 tag reads as 0.4.0", "vv0.4.0", "v0.4.0", false},
		{"a vv0.4.0 build sees 0.4.1", "v0.4.1", "vv0.4.0", true},
		{"no v at all", "0.5.0", "0.4.0", true},
		{"a release beats its own release candidate", "v1.0.0", "v1.0.0-rc.2", true},
		{"a release candidate is older than the release", "v1.0.0-rc.2", "v1.0.0", false},
		{"development builds never ask to update", "v9.9.9", "dev", false},
		{"unreadable latest tag", "nightly", "v0.4.0", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isNewer(tt.latest, tt.current); got != tt.want {
				t.Errorf("isNewer(%q, %q) = %v, want %v", tt.latest, tt.current, got, tt.want)
			}
		})
	}
}
