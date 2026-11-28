package router

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// durationToSecsGo is a pure-Go mirror of the RouterOS `durationToSecs` local
// function embedded in durationToSecsFn (see client.go). It exists purely to make
// the parsing ALGORITHM testable in Go — the actual RouterOS script text in
// client.go is never modified or executed here. Keep this in lockstep with the
// RSC source by hand, control-flow for control-flow; if you change one, change
// the other.
//
// RSC source being mirrored (see durationToSecsFn in client.go for the real,
// authoritative copy):
//
//	:local durationToSecs do={
//	  :local s $1
//	  :local slen [:len $s]
//	  :local total 0
//	  :local num ""
//	  :local i 0
//	  :while ($i < $slen) do={
//	    :local ch [:pick $s $i ($i + 1)]
//	    :if ($ch = "w") do={
//	      :set total ($total + [:tonum $num] * 604800); :set num ""; :set i ($i + 1)
//	    } else={ :if ($ch = "d") do={
//	      :set total ($total + [:tonum $num] * 86400); :set num ""; :set i ($i + 1)
//	    } else={ :if ($ch = "h") do={
//	      :set total ($total + [:tonum $num] * 3600); :set num ""; :set i ($i + 1)
//	    } else={ :if ($ch = "m") do={
//	      :set total ($total + [:tonum $num] * 60); :set num ""; :set i ($i + 1)
//	    } else={ :if ($ch = "s") do={
//	      :set total ($total + [:tonum $num]); :set num ""; :set i ($i + 1)
//	    } else={ :if ($ch = ":") do={
//	      # Remainder from here on is HH:MM:SS — consume it whole.
//	      :local rest ($num . [:pick $s $i $slen])
//	      :local h [:tonum [:pick $rest 0 [:find $rest ":"]]]
//	      :local rest2 [:pick $rest ([:find $rest ":"] + 1) [:len $rest]]
//	      :local m [:tonum [:pick $rest2 0 [:find $rest2 ":"]]]
//	      :local sec [:tonum [:pick $rest2 ([:find $rest2 ":"] + 1) [:len $rest2]]]
//	      :set total ($total + $h * 3600 + $m * 60 + $sec)
//	      :set num ""; :set i $slen
//	    } else={
//	      :set num ($num . $ch); :set i ($i + 1)
//	    } } } } } }
//	  }
//	  :return $total
//	}
//
// [:pick $s $start $end] is 0-based with an EXCLUSIVE end index (confirmed
// against MikroTik's own docs and forum examples: [:pick "abcde" 1 3] = "bc").
// [:find $s substr] returns the 0-based index of the first match. Both are
// mirrored below using Go string slicing, which has identical semantics for
// the byte-oriented, ASCII-only inputs RouterOS hands us here.
func durationToSecsGo(s string) int64 {
	slen := len(s)
	var total int64
	var num strings.Builder
	i := 0
	toNum := func(str string) int64 {
		// RouterOS [:tonum ""] returns 0, matching strconv on an empty accumulator.
		var n int64
		for _, d := range str {
			if d < '0' || d > '9' {
				return 0
			}
			n = n*10 + int64(d-'0')
		}
		return n
	}
	for i < slen {
		ch := s[i]
		switch ch {
		case 'w':
			total += toNum(num.String()) * 604800
			num.Reset()
			i++
		case 'd':
			total += toNum(num.String()) * 86400
			num.Reset()
			i++
		case 'h':
			total += toNum(num.String()) * 3600
			num.Reset()
			i++
		case 'm':
			total += toNum(num.String()) * 60
			num.Reset()
			i++
		case 's':
			total += toNum(num.String())
			num.Reset()
			i++
		case ':':
			// Remainder from here on is HH:MM:SS -- consume it whole, mirroring
			// the RSC colon branch's :find/:pick chain exactly.
			rest := num.String() + s[i:slen]
			restColon := strings.Index(rest, ":")
			h := toNum(rest[0:restColon])
			rest2 := rest[restColon+1:]
			rest2Colon := strings.Index(rest2, ":")
			m := toNum(rest2[0:rest2Colon])
			sec := toNum(rest2[rest2Colon+1:])
			total += h*3600 + m*60 + sec
			num.Reset()
			i = slen
		default:
			num.WriteByte(ch)
			i++
		}
	}
	return total
}

func TestDurationToSecsGo(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		// --- Form 1: letter-suffixed, combinable (Pikro's Go API client's
		// normalized form; the only form ever visible via the frontend/API). ---
		{name: "minutes", input: "5m", want: 300},
		{name: "hours", input: "1h", want: 3600},
		{name: "days", input: "1d", want: 86400},
		{name: "weeks", input: "1w", want: 604800},
		{name: "seconds", input: "30s", want: 30},
		{name: "hours and minutes", input: "1h30m", want: 5400},
		{name: "minutes and seconds", input: "11m44s", want: 704},
		{name: "days and hours", input: "1d12h", want: 129600},
		{name: "zero seconds", input: "0s", want: 0},
		{name: "empty string means no limit set", input: "", want: 0},

		// Real limit-uptime / uptime values observed from a live router (form 1).
		{name: "live limit-uptime: 1d", input: "1d", want: 86400},
		{name: "live limit-uptime: 1h30m", input: "1h30m", want: 5400},
		{name: "live limit-uptime: 1w", input: "1w", want: 604800},
		{name: "live limit-uptime: 2w", input: "2w", want: 1209600},
		{name: "live limit-uptime: 4h", input: "4h", want: 14400},
		{name: "live limit-uptime: 5m", input: "5m", want: 300},
		{name: "live uptime: 11m44s", input: "11m44s", want: 704},
		{name: "live uptime: 48m16s", input: "48m16s", want: 2896},
		{name: "live uptime: 53m20s", input: "53m20s", want: 3200},

		// --- Form 2: raw RouterOS script value, week/day letters followed by a
		// colon-joined HH:MM:SS time-of-day. This is what a .rsc script actually
		// reads via [/ip/hotspot/user/get $u uptime] -- confirmed via a live
		// router's own :log output, NOT via Pikro's API (which never sees this
		// form). The earlier parser (form-1-only) silently returned 0 for all of
		// these, which caused the production incident this test guards against. ---
		{name: "under a day, zero uptime", input: "00:00:00", want: 0},
		{name: "under a day, ten minutes", input: "00:10:00", want: 600},
		{name: "under a day, with seconds", input: "01:02:03", want: 3723},
		{name: "multi-day prefix, zero time-of-day", input: "1d00:00:00", want: 86400},
		{name: "multi-day prefix, with time-of-day", input: "2d05:15:30", want: 191730},
		// 1 week + 3 days + (10h51m24s) = 604800 + 259200 + 39084 = 903084.
		{name: "multi-week prefix, weeks+days+time-of-day", input: "1w3d10:51:24", want: 903084},
		// 1 week only, no day letter, straight into time-of-day.
		{name: "week prefix without day, with time-of-day", input: "1w00:05:00", want: 604800 + 300},

		// --- The real production incident: uptime="00:00:00" (never-used
		// voucher), limit="00:10:00" (10-minute quota). The old parser returned 0
		// for BOTH (no letters to match), so 0 >= 0 was true and every fresh
		// voucher got deleted on the next scheduler tick. This is the single most
		// important assertion in this file. ---
		{name: "incident: fresh voucher uptime", input: "00:00:00", want: 0},
		{name: "incident: ten-minute limit", input: "00:10:00", want: 600},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := durationToSecsGo(tc.input)
			if got != tc.want {
				t.Errorf("durationToSecsGo(%q) = %d, want %d", tc.input, got, tc.want)
			}
		})
	}
}

// TestDurationToSecsGo_incidentRegression is the direct regression test for the
// production incident: a freshly-generated, never-used voucher has
// uptime="00:00:00" and some non-zero limit-uptime like "00:10:00". The cleanup
// scheduler deletes a user when durationToSecs(uptime) >= durationToSecs(limit).
// The bug was that both values parsed to 0 (no letter suffixes matched), so
// 0 >= 0 was true and the voucher was deleted despite never being used. This test
// asserts the fix holds: an unused voucher's uptime must compare as LESS THAN its
// limit, not equal to or greater than it.
func TestDurationToSecsGo_incidentRegression(t *testing.T) {
	uptime := durationToSecsGo("00:00:00")
	limit := durationToSecsGo("00:10:00")

	if !(uptime < limit) {
		t.Fatalf("incident regression: durationToSecsGo(uptime=00:00:00)=%d, durationToSecsGo(limit=00:10:00)=%d; "+
			"want uptime < limit (a fresh, unused voucher must NOT look quota-exhausted)", uptime, limit)
	}
}

// TestCleanupScripts_containUptimeQuotaLogic is a coarse smoke test verifying the
// embedded RouterOS script constants literally contain the new uptime-quota-exhausted
// removal logic (checking `uptime >= limit-uptime`, independent of the `exp:` comment).
// This does NOT execute or interpret the RSC — Go cannot run RouterOS scripting
// language. It only guards against a future edit silently regressing the feature back
// out (e.g. someone "cleaning up" the script and dropping the uptime branch). Real
// correctness must be validated via the project's RouterOS linter
// (.agents/skills/mikrotik-routeros-rsc/scripts/lint_rsc.py) and, ultimately, on real
// hardware.
func TestCleanupScripts_containUptimeQuotaLogic(t *testing.T) {
	scripts := []struct {
		name   string
		script string
	}{
		{name: "cleanupScriptV7", script: cleanupScriptV7},
		{name: "cleanupScriptV6", script: cleanupScriptV6},
	}

	wantSubstrings := []string{
		"limit-uptime",
		"uptime",
		"$durationToSecs $used",
		"$durationToSecs $limit",
		"!$removed",
	}

	for _, s := range scripts {
		t.Run(s.name, func(t *testing.T) {
			for _, want := range wantSubstrings {
				if !strings.Contains(s.script, want) {
					t.Errorf("script does not contain expected marker %q", want)
				}
			}
		})
	}
}

// Removing a hotspot user leaves its open session running, and a device
// logged in by MAC cookie stays online past its validity. Found on a 6.49.13
// router: an expired voucher was removed and still connected an hour later.
func TestCleanupScripts_endSessionsOfRemovedUsers(t *testing.T) {
	tests := []struct {
		name, script, active, cookie, user string
	}{
		{"cleanupScriptV7", cleanupScriptV7, "/ip/hotspot/active/remove [find where user=$who]",
			"/ip/hotspot/cookie/remove [find where user=$who]", "/ip/hotspot/user/remove $u"},
		{"cleanupScriptV6", cleanupScriptV6, "/ip hotspot active remove [find where user=$who]",
			"/ip hotspot cookie remove [find where user=$who]", "/ip hotspot user remove $u"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, want := range []string{tc.active, tc.cookie} {
				if !strings.Contains(tc.script, want) {
					t.Errorf("script lacks %q", want)
				}
			}
			if strings.Contains(tc.script, tc.user) {
				t.Errorf("script still removes a user alone (%q): its session would stay online", tc.user)
			}
		})
	}
}

// TestCleanupScripts_defineDurationToSecsOnce verifies both script variants embed
// exactly one copy of the shared durationToSecsFn snippet. A future refactor could
// accidentally duplicate or drop this shared local-function definition (e.g. if
// someone inlines it per-variant); RouterOS would error on a duplicate `:local`
// function of the same name in the same scope, or fail outright if it's missing when
// the uptime-quota check tries to call it.
func TestCleanupScripts_defineDurationToSecsOnce(t *testing.T) {
	tests := []struct {
		name   string
		script string
	}{
		{name: "cleanupScriptV7", script: cleanupScriptV7},
		{name: "cleanupScriptV6", script: cleanupScriptV6},
	}

	const def = ":local durationToSecs do="
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := strings.Count(tc.script, def); got != 1 {
				t.Errorf("script contains %d copies of durationToSecs definition, want exactly 1", got)
			}
		})
	}
}

// TestCleanupScripts_lintClean runs the project's RouterOS linter against both
// embedded script variants. RSC bodies cannot be unit-tested (no RouterOS
// interpreter in Go) so this is the closest thing to static verification: it dumps
// each script constant to a temp .rsc file, lints it with --strict, and cleans up
// afterward. Skips gracefully if python3 or the linter script is unavailable so this
// suite doesn't become flaky in environments without them.
func TestCleanupScripts_lintClean(t *testing.T) {
	linter := filepath.Join("..", "..", ".agents", "skills", "mikrotik-routeros-rsc", "scripts", "lint_rsc.py")
	if _, err := os.Stat(linter); err != nil {
		t.Skipf("lint_rsc.py not found at %s: %v", linter, err)
	}
	if _, err := exec.LookPath("python3"); err != nil {
		t.Skip("python3 not available")
	}

	scripts := []struct {
		name   string
		script string
	}{
		{name: "cleanupScriptV7", script: cleanupScriptV7},
		{name: "cleanupScriptV6", script: cleanupScriptV6},
	}

	for _, s := range scripts {
		t.Run(s.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, s.name+".rsc")
			if err := os.WriteFile(path, []byte(s.script), 0o644); err != nil {
				t.Fatalf("failed to write temp rsc file: %v", err)
			}

			cmd := exec.Command("python3", linter, "--strict", path)
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Errorf("lint_rsc.py --strict flagged %s:\n%s", s.name, out)
			}
		})
	}
}

// RouterOS 6 only understands space-separated command paths inside scripts
// ("/ip hotspot user find"); the slash form ("/ip/hotspot/user/find") is
// RouterOS 7 syntax, and 6.49 stops the whole script with "expected command
// name". Found on a CHR 6.49.23: no voucher was stamped or cleaned up.
func TestPreV712ScriptsUseSpacedCommandPaths(t *testing.T) {
	// A command starts a line or follows "[" (comments may mention "mmm/dd/yyyy").
	slashPath := regexp.MustCompile(`(?m)(^\s*|\[)/[a-z-]+/[a-z-]+`)
	scripts := []struct {
		name, body string
	}{
		{"on-login, RouterOS 6", buildOnLoginScript(3600, "6.49.23 (stable)")},
		{"on-login, RouterOS 7.11", buildOnLoginScript(3600, "7.11.2 (stable)")},
		{"cleanup, RouterOS 6", cleanupScriptBody("6.49.23 (stable)")},
		{"cleanup, RouterOS 7.11", cleanupScriptBody("7.11.2 (stable)")},
	}
	for _, s := range scripts {
		t.Run(s.name, func(t *testing.T) {
			if m := slashPath.FindAllString(s.body, -1); len(m) > 0 {
				t.Errorf("slash command paths RouterOS 6 rejects: %v", m)
			}
		})
	}
}

// Routers set up by an older Pikro keep its old scripts until they are
// rewritten, which needs each plan's validity read back from its script.
func TestPikroOnLoginValiditySecs(t *testing.T) {
	for _, version := range []string{"7.25 (stable)", "7.11.2 (stable)", "6.49.23 (stable)"} {
		for _, secs := range []int64{120, 3600, 2592000} {
			t.Run(version, func(t *testing.T) {
				got, ok := PikroOnLoginValiditySecs(buildOnLoginScript(secs, version))
				if !ok || got != secs {
					t.Errorf("validity = %d, %v, want %d", got, ok, secs)
				}
			})
		}
	}
	for _, script := range []string{"", `:put (",rem,500,1d,500,,Disable,")`} {
		if _, ok := PikroOnLoginValiditySecs(script); ok {
			t.Errorf("%q read as a Pikro script", script)
		}
	}
}

// TestScripts_noPropertyNamedVariableInFind guards against a RouterOS trap: in
// `find where name=$name` RouterOS reads $name as the item's own name
// property, so the filter matches every item. On RouterOS 6 that made the
// cleanup try to remove all hotspot users at once (refused because of
// default-trial), so it removed none and stopped at the first expired voucher.
func TestScripts_noPropertyNamedVariableInFind(t *testing.T) {
	filter := regexp.MustCompile(`find where ([a-z-]+)=\$([a-zA-Z-]+)`)
	tests := []struct{ name, script string }{
		{"cleanupScriptV7", cleanupScriptV7},
		{"cleanupScriptV6", cleanupScriptV6},
		{"on-login 7.25", buildOnLoginScript(3600, "7.25")},
		{"on-login 6.49", buildOnLoginScript(3600, "6.49.13")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			for _, m := range filter.FindAllStringSubmatch(tc.script, -1) {
				if m[1] == m[2] {
					t.Errorf("%q compares the %s property with itself: it matches every item", m[0], m[1])
				}
			}
		})
	}
}

// TestCleanupV7_readsStampBeforeNote: the on-login script keeps an existing
// comment after the stamp ("exp:1793794523 pikro-cloud"), so the cleanup must
// read the 10 digits only; the rest of the comment isn't a number.
func TestCleanupV7_readsStampBeforeNote(t *testing.T) {
	if !strings.Contains(cleanupScriptV7, "[:tonum [:pick $comment 4 14]]") {
		t.Error("cleanupScriptV7 doesn't read exactly the 10-digit epoch after exp:")
	}
}

// TestOnLoginBefore712_readsEveryNextRunForm: the pre-7.12 on-login script
// turns a scheduler's next-run into a date; RouterOS writes it in a different
// form per version and distance (seen on CHR 6.48–7.11). 7.10–7.11 write
// "11-26 17:02:15" later in the same year: without that form, a voucher
// lasting past midnight was never stamped, so it never expired.
func TestOnLoginBefore712_readsEveryNextRunForm(t *testing.T) {
	script := buildOnLoginScript(86400, "7.11.2")
	for _, tc := range []struct{ form, check string }{
		{"17:04:13 (today)", ":if ($xlen = 8)"},
		{"11-26 17:02:15 (7.10-7.11, this year)", ":if ($xlen = 14)"},
		{"nov/26 17:02:01 (this year)", ":if ($xlen = 15)"},
		{"2027-01-04 17:02:16 (7.10+, another year)", ":if ($xlen = 19)"},
		{"jan/04/2027 17:02:03 (another year)", ":if ($xlen = 20)"},
	} {
		if !strings.Contains(script, tc.check) {
			t.Errorf("no branch for next-run %s", tc.form)
		}
	}
}
