package router

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// durationToSecsGo is a pure-Go mirror of the RouterOS `durationToSecs` local
// function embedded in durationToSecsFn (see client.go). It exists purely to make
// the parsing ALGORITHM testable in Go — the actual RouterOS script text in
// client.go is never modified or executed here. Keep this in lockstep with the
// RSC source by hand; if you change one, change the other.
//
// RSC source being mirrored:
//
//	:local durationToSecs do={
//	  :local s $1
//	  :local total 0
//	  :local num ""
//	  for i from=0 to=([:len $s] - 1) do={
//	    :local ch [:pick $s $i ($i + 1)]
//	    :if ($ch = "w" || $ch = "d" || $ch = "h" || $ch = "m" || $ch = "s") do={
//	      :local n [:tonum $num]
//	      :if ($ch = "w") do={ :set total ($total + $n * 604800) }
//	      :if ($ch = "d") do={ :set total ($total + $n * 86400) }
//	      :if ($ch = "h") do={ :set total ($total + $n * 3600) }
//	      :if ($ch = "m") do={ :set total ($total + $n * 60) }
//	      :if ($ch = "s") do={ :set total ($total + $n) }
//	      :set num ""
//	    } else={
//	      :set num ($num . $ch)
//	    }
//	  }
//	  :return $total
//	}
func durationToSecsGo(s string) int64 {
	var total int64
	var num strings.Builder
	multiplier := map[byte]int64{
		'w': 604800,
		'd': 86400,
		'h': 3600,
		'm': 60,
		's': 1,
	}
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if mult, isUnit := multiplier[ch]; isUnit {
			// RouterOS [:tonum ""] returns 0, matching strconv on an empty accumulator.
			var n int64
			if num.Len() > 0 {
				for _, d := range num.String() {
					if d < '0' || d > '9' {
						n = 0
						break
					}
					n = n*10 + int64(d-'0')
				}
			}
			total += n * mult
			num.Reset()
		} else {
			num.WriteByte(ch)
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
		// Single units.
		{name: "minutes", input: "5m", want: 300},
		{name: "hours", input: "1h", want: 3600},
		{name: "days", input: "1d", want: 86400},
		{name: "weeks", input: "1w", want: 604800},
		{name: "seconds", input: "30s", want: 30},

		// Combined units.
		{name: "hours and minutes", input: "1h30m", want: 5400},
		{name: "minutes and seconds", input: "11m44s", want: 704},
		{name: "days and hours", input: "1d12h", want: 129600},

		// Edge cases.
		{name: "zero seconds", input: "0s", want: 0},
		{name: "empty string means no limit set", input: "", want: 0},

		// Real limit-uptime values observed from a live router.
		{name: "live limit-uptime: 1d", input: "1d", want: 86400},
		{name: "live limit-uptime: 1h30m", input: "1h30m", want: 5400},
		{name: "live limit-uptime: 1w", input: "1w", want: 604800},
		{name: "live limit-uptime: 2w", input: "2w", want: 1209600},
		{name: "live limit-uptime: 4h", input: "4h", want: 14400},
		{name: "live limit-uptime: 5m", input: "5m", want: 300},

		// Real uptime values observed from a live router.
		{name: "live uptime: 0s", input: "0s", want: 0},
		{name: "live uptime: 11m44s", input: "11m44s", want: 704},
		{name: "live uptime: 48m16s", input: "48m16s", want: 2896},
		{name: "live uptime: 53m20s", input: "53m20s", want: 3200},
		{name: "live uptime: 5m", input: "5m", want: 300},
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
	linter := "/Users/gino/Projects/web/pikro/.agents/skills/mikrotik-routeros-rsc/scripts/lint_rsc.py"
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
