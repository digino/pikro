package sales

import (
	"strconv"
	"testing"
	"time"
)

func TestIsActivated(t *testing.T) {
	tests := []struct {
		name string
		user map[string]any
		want bool
	}{
		{"fresh", map[string]any{"uptime": "0s", "bytes-in": "0", "comment": ""}, false},
		{"fresh, clock-style uptime", map[string]any{"uptime": "00:00:00", "bytes-in": "0"}, false},
		{"clock-style uptime", map[string]any{"uptime": "00:12:30"}, true},
		{"no fields", map[string]any{}, false},
		{"exp comment", map[string]any{"comment": "exp:1700000000"}, true},
		{"bytes in", map[string]any{"bytes-in": "10"}, true},
		{"uptime", map[string]any{"uptime": "5m"}, true},
		{"other comment", map[string]any{"comment": "batch x"}, false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := isActivated(tc.user); got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFirstLogin(t *testing.T) {
	now := at("2026-06-10T12:00:00Z")
	hour := int64(3600)
	exp := func(tm string) string { return strconv.FormatInt(at(tm).Unix(), 10) }
	tests := []struct {
		name     string
		comment  string
		validity int64
		want     time.Time
	}{
		{"epoch", "exp:" + exp("2026-06-10T13:00:00Z") + " note", hour, at("2026-06-10T12:00:00Z").Add(0)},
		{"epoch earlier", "exp:" + exp("2026-06-10T10:00:00Z"), hour, at("2026-06-10T09:00:00Z")},
		{"v6 date", "exp:2026-06-10 10:00:00 note", hour, at("2026-06-10T09:00:00Z")},
		{"unknown validity", "exp:" + exp("2026-06-10T10:00:00Z"), 0, now},
		{"no exp", "", hour, now},
		{"future clamp", "exp:" + exp("2026-06-11T10:00:00Z"), hour, now},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := firstLogin(map[string]any{"comment": tc.comment}, tc.validity, now)
			if !got.Equal(tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestSyncActivationsIdempotentAcrossYears(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	now := at("2026-01-02T12:00:00Z")
	if err := AppendBatch("r", "b1", at("2025-12-31T10:00:00Z"), "1d", "500", "XOF", []string{"a", "b", "c"}); err != nil {
		t.Fatal(err)
	}
	users := []map[string]any{
		{"name": "a", "profile": "1d", "comment": "exp:2026-01-02 00:00:00"}, // first login 2026-01-01 00:00
		{"name": "b", "profile": "1d", "uptime": "0s", "bytes-in": "0"},
		{"name": "manual", "profile": "1d", "uptime": "1h"}, // not in any batch
		{"name": "c", "profile": "1d", "uptime": "3m"},
	}
	for i := 0; i < 2; i++ {
		if err := SyncActivations("r", users, map[string]int64{"1d": 86400}, now); err != nil {
			t.Fatal(err)
		}
	}
	got, err := Load("r", 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.Activations) != 2 {
		t.Fatalf("activations = %+v, want 2", got.Activations)
	}
	a := got.Activations[0]
	if a.Name != "a" || a.Batch != "b1" || !a.At.Equal(at("2026-01-01T00:00:00Z")) ||
		a.Profile != "1d" || a.Price != "500" || a.Currency != "XOF" {
		t.Errorf("activation a = %+v", a)
	}
	if c := got.Activations[1]; c.Name != "c" || !c.At.Equal(now) {
		t.Errorf("activation c = %+v", c)
	}
	// The 2025 report has the batch but no activations.
	prev, _ := Load("r", 2025)
	if len(prev.Batches) != 1 || prev.Batches[0].Count != 3 || len(prev.Activations) != 0 {
		t.Errorf("2025 ledger = %+v", prev)
	}
}
