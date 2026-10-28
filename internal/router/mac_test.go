package router

import "testing"

func TestPickMAC(t *testing.T) {
	addrs := []map[string]string{
		{"address": "10.5.50.1/24", "interface": "bridge-hotspot"},
		{"address": "192.168.1.5/24", "interface": "ether2"},
	}
	ifaces := []map[string]string{
		{"name": "ether1", "type": "ether", "mac-address": "AA:AA:AA:AA:AA:01"},
		{"name": "ether2", "type": "ether", "mac-address": "AA:AA:AA:AA:AA:02"},
		{"name": "bridge-hotspot", "type": "bridge", "mac-address": "AA:AA:AA:AA:AA:03"},
	}
	tests := []struct {
		name   string
		host   string
		ifaces []map[string]string
		want   string
	}{
		{"interface carrying the host IP", "192.168.1.5", ifaces, "AA:AA:AA:AA:AA:02"},
		{"host IP on a bridge", "10.5.50.1", ifaces, "AA:AA:AA:AA:AA:03"},
		{"host not on the router (NAT, DNS name) falls back to first ethernet", "127.0.0.1", ifaces, "AA:AA:AA:AA:AA:01"},
		{"no ethernet at all", "127.0.0.1", ifaces[2:], ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pickMAC(tt.host, addrs, tt.ifaces); got != tt.want {
				t.Errorf("pickMAC(%q) = %q, want %q", tt.host, got, tt.want)
			}
		})
	}
}
