package router

import (
	"testing"
	"time"
)

func TestMikhmonExpiryToEpoch(t *testing.T) {
	tests := []struct {
		name      string
		comment   string
		gmtOffset time.Duration
		wantEpoch int64
		wantOK    bool
	}{
		{
			name:      "zero offset (Lome, UTC+0)",
			comment:   "2026-03-15 12:00:00",
			gmtOffset: 0,
			wantEpoch: 1773576000,
			wantOK:    true,
		},
		{
			name:      "positive offset (Lagos, +01:00) shifts epoch earlier",
			comment:   "2026-03-15 12:00:00",
			gmtOffset: 1 * time.Hour,
			wantEpoch: 1773572400,
			wantOK:    true,
		},
		{
			name:      "negative offset (-05:00) shifts epoch later",
			comment:   "2026-03-15 12:00:00",
			gmtOffset: -5 * time.Hour,
			wantEpoch: 1773594000,
			wantOK:    true,
		},
		{
			name:      "leap day Feb 29 at midnight boundary minus one second",
			comment:   "2028-02-29 23:59:59",
			gmtOffset: 0,
			wantEpoch: 1835481599,
			wantOK:    true,
		},
		{
			name:      "midnight boundary",
			comment:   "2026-06-01 00:00:00",
			gmtOffset: 0,
			wantEpoch: 1780272000,
			wantOK:    true,
		},
		{
			name:      "year boundary, zero offset",
			comment:   "2026-12-31 23:59:59",
			gmtOffset: 0,
			wantEpoch: 1798761599,
			wantOK:    true,
		},
		{
			name:      "year boundary crosses into next UTC year with positive offset",
			comment:   "2026-12-31 23:59:59",
			gmtOffset: 1 * time.Hour,
			wantEpoch: 1798757999,
			wantOK:    true,
		},
		{
			name:      "Mikhmon unused-voucher format is not an expiry timestamp",
			comment:   "up-541-06.23.26-",
			gmtOffset: 0,
			wantOK:    false,
		},
		{
			name:      "bare vc- prefix does not match",
			comment:   "vc-",
			gmtOffset: 0,
			wantOK:    false,
		},
		{
			name:      "empty comment",
			comment:   "",
			gmtOffset: 0,
			wantOK:    false,
		},
		{
			name:      "free text does not match",
			comment:   "not a date",
			gmtOffset: 0,
			wantOK:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotEpoch, gotOK := mikhmonExpiryToEpoch(tc.comment, tc.gmtOffset)
			if gotOK != tc.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tc.wantOK)
			}
			if gotOK && gotEpoch != tc.wantEpoch {
				t.Errorf("epoch = %d, want %d", gotEpoch, tc.wantEpoch)
			}
		})
	}
}

func TestMikhmonProfileValiditySecs(t *testing.T) {
	tests := []struct {
		name     string
		onLogin  string
		wantSecs int64
		wantOK   bool
	}{
		{
			name:     "no-expiry marker",
			onLogin:  `:put (",,0,,,noexp,Disable,")`,
			wantSecs: 0,
			wantOK:   true,
		},
		{
			name:     "2 hours",
			onLogin:  `:put (",rem,50,2h,50,,Disable,")`,
			wantSecs: 7200,
			wantOK:   true,
		},
		{
			name:     "24 hours",
			onLogin:  `:put (",rem,200,24h,200,,Disable,")`,
			wantSecs: 86400,
			wantOK:   true,
		},
		{
			name:     "7 days",
			onLogin:  `:put (",rem,500,7d,500,,Disable,")`,
			wantSecs: 604800,
			wantOK:   true,
		},
		{
			name:     "14 days",
			onLogin:  `:put (",rem,1000,14d,1000,,Disable,")`,
			wantSecs: 1209600,
			wantOK:   true,
		},
		{
			name:     "30 days",
			onLogin:  `:put (",rem,1500,30d,1500,,Disable,")`,
			wantSecs: 2592000,
			wantOK:   true,
		},
		{
			name:    "empty on-login script",
			onLogin: "",
			wantOK:  false,
		},
		{
			name:    "already Pikro's own on-login script is not re-migrated",
			onLogin: ":local nowEpoch ([:tonsec [:timestamp]] / 1000000000)\n:local expEpoch ($nowEpoch + 7200)",
			wantOK:  false,
		},
		{
			name:    "unrelated garbage on-login script",
			onLogin: `:local x 1; :log info "hello"`,
			wantOK:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotSecs, gotOK := mikhmonProfileValiditySecs(tc.onLogin)
			if gotOK != tc.wantOK {
				t.Fatalf("ok = %v, want %v", gotOK, tc.wantOK)
			}
			if gotOK && gotSecs != tc.wantSecs {
				t.Errorf("secs = %d, want %d", gotSecs, tc.wantSecs)
			}
		})
	}
}

func TestParseRouterOSDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "seconds", input: "45s", want: 45},
		{name: "minutes", input: "10m", want: 600},
		{name: "hours", input: "2h", want: 7200},
		{name: "days", input: "7d", want: 604800},
		{name: "weeks", input: "1w", want: 604800},
		{name: "empty string is an error", input: "", wantErr: true},
		{name: "non-numeric garbage is an error", input: "abc", wantErr: true},
		{name: "unknown unit is an error", input: "5x", wantErr: true},
		{name: "unit with no number is an error", input: "h", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parseRouterOSDuration(tc.input)
			if (err != nil) != tc.wantErr {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantErr)
			}
			if err == nil && got != tc.want {
				t.Errorf("secs = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestParseGMTOffset(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  time.Duration
	}{
		{name: "zero offset", input: "+00:00", want: 0},
		{name: "positive hour offset", input: "+01:00", want: 1 * time.Hour},
		{name: "negative offset with half-hour minutes", input: "-05:30", want: -5*time.Hour - 30*time.Minute},
		{name: "too short falls back to zero", input: "+1:0", want: 0},
		{name: "garbage falls back to zero", input: "garbage", want: 0},
		{name: "empty string falls back to zero", input: "", want: 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := parseGMTOffset(tc.input)
			if got != tc.want {
				t.Errorf("parseGMTOffset(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}
