package router

import "testing"

func TestPoolSize(t *testing.T) {
	tests := []struct {
		name   string
		ranges string
		want   int
	}{
		{"one range", "10.5.50.2-10.5.50.254", 253},
		{"several ranges", "192.168.88.10-192.168.88.100,192.168.88.200-192.168.88.209", 101},
		{"single address", "10.0.0.5", 1},
		{"CIDR", "10.0.0.0/24", 256},
		{"range across octets", "10.0.0.250-10.0.1.5", 12},
		{"spaces after commas", "10.0.0.1-10.0.0.10, 10.0.0.20", 11},
		{"empty", "", 0},
		{"IPv6 is not counted", "2001:db8::1-2001:db8::ff", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := poolSize(tt.ranges); got != tt.want {
				t.Errorf("poolSize(%q) = %d, want %d", tt.ranges, got, tt.want)
			}
		})
	}
}
