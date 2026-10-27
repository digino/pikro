package sales

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	expDateRe  = regexp.MustCompile(`^exp:(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})`)
	expEpochRe = regexp.MustCompile(`^exp:(\d+)`)
)

// SyncActivations records an activation for every listed user that belongs to a
// batch, has been used, and has no activation yet. validityByProfile maps a
// profile name to its validity in seconds (0 or absent = unknown).
// Idempotent: activations are deduplicated on batch+name across the current
// and previous year.
func SyncActivations(routerID string, users []map[string]any, validityByProfile map[string]int64, now time.Time) error {
	now = now.UTC()
	mu.Lock()
	defer mu.Unlock()

	owner := map[string]Batch{} // voucher name -> batch
	done := map[[2]string]bool{}
	for _, year := range []int{now.Year() - 1, now.Year()} {
		batches, acts, err := readYear(routerID, year)
		if err != nil {
			return err
		}
		for _, b := range batches {
			for _, n := range b.Names {
				owner[n] = b
			}
		}
		for _, a := range acts {
			done[[2]string{a.Batch, a.Name}] = true
		}
	}

	for _, u := range users {
		name, _ := u["name"].(string)
		b, ok := owner[name]
		if !ok || done[[2]string{b.ID, name}] || !isActivated(u) {
			continue
		}
		profile, _ := u["profile"].(string)
		at := firstLogin(u, validityByProfile[profile], now)
		if err := appendLine(routerID, at.Year(), line{Type: typeActivation, Batch: b.ID, Name: name, At: at}); err != nil {
			return err
		}
		done[[2]string{b.ID, name}] = true
	}
	return nil
}

func isActivated(u map[string]any) bool {
	comment, _ := u["comment"].(string)
	if strings.HasPrefix(comment, "exp:") {
		return true
	}
	if in, _ := u["bytes-in"].(string); in != "" {
		if n, err := strconv.ParseInt(in, 10, 64); err == nil && n > 0 {
			return true
		}
	}
	// Any non-zero digit means time online, whichever form RouterOS reports:
	// "0s"/"00:00:00" never used, "5m"/"00:12:30" used.
	uptime, _ := u["uptime"].(string)
	return strings.ContainsAny(uptime, "123456789")
}

// firstLogin derives the first-login time from the "exp:" marker (expiry minus
// validity). Falls back to now when it cannot, and never returns a future time.
func firstLogin(u map[string]any, validitySecs int64, now time.Time) time.Time {
	comment, _ := u["comment"].(string)
	at := now
	if validitySecs > 0 {
		validity := time.Duration(validitySecs) * time.Second
		if m := expDateRe.FindStringSubmatch(comment); m != nil {
			if t, err := time.Parse("2006-01-02 15:04:05", m[1]); err == nil {
				at = t.UTC().Add(-validity)
			}
		} else if m := expEpochRe.FindStringSubmatch(comment); m != nil {
			if sec, err := strconv.ParseInt(m[1], 10, 64); err == nil {
				at = time.Unix(sec, 0).UTC().Add(-validity)
			}
		}
	}
	if at.After(now) {
		return now
	}
	return at
}
