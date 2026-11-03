package router

import "testing"

func TestWANInterface(t *testing.T) {
	tests := []struct {
		name   string
		routes []map[string]string
		want   string
	}{
		{"v7 gateway through an interface", []map[string]string{
			{"active": "true", "gateway": "10.0.2.2", "immediate-gw": "10.0.2.2%ether1"},
		}, "ether1"},
		{"v6 gateway through an interface", []map[string]string{
			{"active": "true", "gateway": "10.0.2.2", "gateway-status": "10.0.2.2 reachable via  ether1"},
		}, "ether1"},
		{"v7 PPPoE: the gateway is the interface", []map[string]string{
			{"active": "true", "gateway": "pppoe-out1", "immediate-gw": "pppoe-out1"},
		}, "pppoe-out1"},
		{"v6 PPPoE: the gateway is the interface", []map[string]string{
			{"active": "true", "gateway": "pppoe-out1", "gateway-status": "pppoe-out1 reachable"},
		}, "pppoe-out1"},
		{"v6 several gateways: the first one", []map[string]string{
			{"active": "true", "gateway-status": "10.0.0.1 reachable via  ether1,10.1.0.1 unreachable"},
		}, "ether1"},
		{"skips a backup route that isn't active", []map[string]string{
			{"active": "false", "immediate-gw": "10.9.0.1%lte1"},
			{"active": "true", "immediate-gw": "10.0.2.2%ether1"},
		}, "ether1"},
		{"no default route", nil, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := wanInterface(tt.routes); got != tt.want {
				t.Errorf("wanInterface() = %q, want %q", got, tt.want)
			}
		})
	}
}
