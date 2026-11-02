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
