package router

import "testing"

func TestAddressIP(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"with prefix length", "10.5.50.1/24", "10.5.50.1"},
		{"bare address", "192.168.88.1", "192.168.88.1"},
		{"empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addressIP(tt.in); got != tt.want {
				t.Errorf("addressIP(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
