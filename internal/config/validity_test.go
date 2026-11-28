package config

import "testing"

func TestValiditySecs(t *testing.T) {
	tests := []struct {
		in   string
		want int64
	}{
		{"", 0}, {"junk", 0}, {"30m", 1800}, {"1h", 3600}, {"1d", 86400},
		{"1w", 604800}, {"1w2d", 777600}, {"2h30m", 9000},
	}
	for _, tc := range tests {
		t.Run(tc.in, func(t *testing.T) {
			if got := ValiditySecs(tc.in); got != tc.want {
				t.Errorf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestValidityString(t *testing.T) {
	tests := []struct {
		secs int64
		want string
	}{
		{0, ""}, {1800, "30m"}, {3600, "1h"}, {18000, "5h"}, {86400, "1d"},
		{129600, "1d12h"}, {604800, "7d"}, {2592000, "30d"}, {9000, "2h30m"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			got := ValidityString(tc.secs)
			if got != tc.want {
				t.Errorf("ValidityString(%d) = %q, want %q", tc.secs, got, tc.want)
			}
			if ValiditySecs(got) != tc.secs {
				t.Errorf("%q reads back as %d seconds, want %d", got, ValiditySecs(got), tc.secs)
			}
		})
	}
}
