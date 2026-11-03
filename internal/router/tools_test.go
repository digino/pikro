package router

import (
	"testing"
	"time"
)

func TestPingSummary(t *testing.T) {
	tests := []struct {
		name         string
		replies      []map[string]string
		wantReceived int
		wantAvgMs    int
	}{
		{"v6 answers", []map[string]string{
			{"sent": "1", "received": "1", "avg-rtt": "5ms"},
			{"sent": "2", "received": "2", "avg-rtt": "4ms"},
		}, 2, 4},
		{"v7 sub-millisecond parts", []map[string]string{
			{"sent": "3", "received": "3", "avg-rtt": "12ms345us"},
		}, 3, 12},
		{"under a millisecond rounds up to 1", []map[string]string{
			{"sent": "1", "received": "1", "avg-rtt": "800us"},
		}, 1, 1},
		{"no answer", []map[string]string{
			{"status": "timeout", "sent": "1", "received": "0"},
			{"status": "timeout", "sent": "2", "received": "0"},
		}, 0, 0},
		{"no replies at all", nil, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			received, avg := pingSummary(tt.replies)
			if received != tt.wantReceived || avg != tt.wantAvgMs {
				t.Errorf("pingSummary() = (%d, %d), want (%d, %d)", received, avg, tt.wantReceived, tt.wantAvgMs)
			}
		})
	}
}

func TestRouterClockUTC(t *testing.T) {
	want := time.Date(2026, 10, 29, 22, 36, 12, 0, time.UTC)
	tests := []struct {
		name              string
		date, tod, offset string
		want              time.Time
		wantErr           bool
	}{
		{"v6 date", "oct/29/2026", "22:36:12", "+00:00", want, false},
		{"v7 date", "2026-10-29", "22:36:12", "+00:00", want, false},
		{"ahead of UTC", "2026-10-29", "23:36:12", "+01:00", want, false},
		{"behind UTC", "2026-10-29", "17:36:12", "-05:00", want, false},
		{"after a power cut", "jan/02/1970", "00:41:00", "+00:00", time.Date(1970, 1, 2, 0, 41, 0, 0, time.UTC), false},
		{"unknown date format", "29.10.2026", "22:36:12", "+00:00", time.Time{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := routerClockUTC(tt.date, tt.tod, tt.offset)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if !got.Equal(tt.want) {
				t.Errorf("routerClockUTC() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRouterClockStrings(t *testing.T) {
	utc := time.Date(2026, 10, 29, 22, 36, 12, 0, time.UTC)
	tests := []struct {
		name               string
		offset             string
		iso                bool
		wantDate, wantTime string
	}{
		{"v6 style", "+00:00", false, "oct/29/2026", "22:36:12"},
		{"v7 style", "+00:00", true, "2026-10-29", "22:36:12"},
		{"local time crosses midnight", "+02:00", true, "2026-10-30", "00:36:12"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, tod := routerClockStrings(utc, tt.offset, tt.iso)
			if d != tt.wantDate || tod != tt.wantTime {
				t.Errorf("routerClockStrings() = (%q, %q), want (%q, %q)", d, tod, tt.wantDate, tt.wantTime)
			}
		})
	}
}

func TestGatewayAddress(t *testing.T) {
	tests := []struct {
		name  string
		route map[string]string
		want  string
	}{
		{"v7", map[string]string{"gateway": "10.0.2.2", "immediate-gw": "10.0.2.2%ether1"}, "10.0.2.2"},
		{"v6", map[string]string{"gateway": "10.0.2.2", "gateway-status": "10.0.2.2 reachable via  ether1"}, "10.0.2.2"},
		{"PPPoE has no address to ping", map[string]string{"gateway": "pppoe-out1"}, ""},
		{"several gateways: the first", map[string]string{"gateway": "10.0.0.1,10.1.0.1"}, "10.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := gatewayAddress(tt.route); got != tt.want {
				t.Errorf("gatewayAddress() = %q, want %q", got, tt.want)
			}
		})
	}
}
