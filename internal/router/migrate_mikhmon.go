package router

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

// Mikhmon writes two comment formats onto hotspot users:
//   - "up-<validitySecs>-<MM.DD.YY>-"  — voucher generated, never logged in yet
//   - "YYYY-MM-DD HH:MM:SS"            — real expiry, stamped by its on-login script after first login
//
// Pikro only ever understands its own "exp:<epoch>" format. MigrateFromMikhmon
// rewrites both patterns into that format (or blanks the comment for unused
// vouchers, matching Pikro's own pre-login state) and replaces each hotspot
// profile's on-login script with Pikro's equivalent, so future logins stamp
// exp: instead of the Mikhmon formats.
var (
	mikhmonUnusedRe = regexp.MustCompile(`^up-\d+-`)
	mikhmonExpiryRe = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2}) (\d{2}):(\d{2}):(\d{2})$`)
	mikhmonOnLoginValidityRe = regexp.MustCompile(`:put \(",[^,]*,\d+,(\w+),`)
)

// MigrationResult reports what MigrateFromMikhmon changed.
type MigrationResult struct {
	UsersScanned      int `json:"usersScanned"`
	UsersUnused       int `json:"usersUnused"`       // comment cleared (never logged in)
	UsersConverted    int `json:"usersConverted"`    // comment converted to exp:<epoch>
	UsersSkipped      int `json:"usersSkipped"`      // comment format not recognized, left untouched
	ProfilesScanned   int `json:"profilesScanned"`
	ProfilesConverted int `json:"profilesConverted"` // on-login script replaced with Pikro's
	// SchedulersRemoved/ScriptsRemoved count Mikhmon's own leftover per-user
	// /system scheduler and /system script entries (see removeMikhmonArtifacts).
	SchedulersRemoved int  `json:"schedulersRemoved"`
	ScriptsRemoved    int  `json:"scriptsRemoved"`
	CleanupInstalled  bool `json:"cleanupInstalled"` // Pikro's own cleanup scheduler was (re-)installed
}

// MigrateFromMikhmon scans all hotspot users and profiles on the router,
// converting Mikhmon-style comments and on-login scripts to Pikro's own
// exp:<epoch> convention. Safe to run multiple times — already-migrated
// users/profiles are left untouched.
func (c *Client) MigrateFromMikhmon() (*MigrationResult, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	result := &MigrationResult{}
	rosVersion := c.rosVersion(conn)

	gmtOffset := time.Duration(0)
	if clockReply, err := conn.Run("/system/clock/print"); err == nil && len(clockReply.Re) > 0 {
		gmtOffset = parseGMTOffset(clockReply.Re[0].Map["gmt-offset"])
	}

	userReply, err := conn.Run("/ip/hotspot/user/print")
	if err != nil {
		return nil, fmt.Errorf("listing hotspot users: %w", err)
	}
	result.UsersScanned = len(userReply.Re)

	usernames := make(map[string]bool, len(userReply.Re))
	for _, s := range userReply.Re {
		if name := s.Map["name"]; name != "" {
			usernames[name] = true
		}
	}

	for _, s := range userReply.Re {
		id := s.Map[".id"]
		comment := s.Map["comment"]

		if strings.HasPrefix(comment, "exp:") {
			continue // already Pikro's format
		}

		if mikhmonUnusedRe.MatchString(comment) {
			if _, err := conn.Run("/ip/hotspot/user/set", "=.id="+id, "=comment="); err != nil {
				return result, fmt.Errorf("clearing comment for user %s: %w", id, err)
			}
			result.UsersUnused++
			continue
		}

		if epoch, ok := mikhmonExpiryToEpoch(comment, gmtOffset); ok {
			newComment := fmt.Sprintf("exp:%d", epoch)
			if _, err := conn.Run("/ip/hotspot/user/set", "=.id="+id, "=comment="+newComment); err != nil {
				return result, fmt.Errorf("converting comment for user %s: %w", id, err)
			}
			result.UsersConverted++
			continue
		}

		result.UsersSkipped++
	}

	profileReply, err := conn.Run("/ip/hotspot/user/profile/print")
	if err != nil {
		return result, fmt.Errorf("listing hotspot profiles: %w", err)
	}
	result.ProfilesScanned = len(profileReply.Re)

	for _, s := range profileReply.Re {
		id := s.Map[".id"]
		onLogin := s.Map["on-login"]

		validitySecs, ok := mikhmonProfileValiditySecs(onLogin)
		if !ok {
			continue // not a Mikhmon-style script, or already Pikro's — leave as-is
		}

		newScript := ""
		if validitySecs > 0 {
			newScript = buildOnLoginScript(validitySecs, rosVersion)
		}
		if _, err := conn.Run("/ip/hotspot/user/profile/set", "=.id="+id, "=on-login="+newScript); err != nil {
			return result, fmt.Errorf("converting profile %s: %w", id, err)
		}
		result.ProfilesConverted++
	}

	if err := removeMikhmonArtifacts(conn, usernames, result); err != nil {
		return result, fmt.Errorf("removing Mikhmon scheduler/script leftovers: %w", err)
	}

	if err := c.InstallCleanupScheduler("7d"); err == nil {
		result.CleanupInstalled = true
	}

	return result, nil
}

// removeMikhmonArtifacts removes Mikhmon's own leftover /system scheduler and
// /system script entries — created one-per-user by its on-login script at
// login time (name == the hotspot username, confirmed against Mikhmon's own
// source: laksa19/mikrotik-hotspot-monitor, mikhmon/app/uprofileadd.php).
// Mikhmon never removes these itself, so on a long-lived router hundreds can
// accumulate. Matches strictly by exact name equality against real hotspot
// usernames on this router, so an admin's own unrelated scheduler/script is
// never touched unless it happens to share a voucher's exact username.
//
// The remc/ntfc on-login variants additionally log each transaction as a
// /system script named "<date>-|-<time>-|-<user>-|-<price>" tagged
// comment=mikhmon (a literal, Mikhmon-specific marker) — those are matched
// on that comment instead, since their names aren't usernames.
func removeMikhmonArtifacts(conn *ros.Client, usernames map[string]bool, result *MigrationResult) error {
	schedReply, err := conn.Run("/system/scheduler/print")
	if err != nil {
		return fmt.Errorf("listing schedulers: %w", err)
	}
	for _, s := range schedReply.Re {
		name := s.Map["name"]
		if !usernames[name] {
			continue
		}
		if _, err := conn.Run("/system/scheduler/remove", "=.id="+s.Map[".id"]); err != nil {
			return fmt.Errorf("removing scheduler %s: %w", name, err)
		}
		result.SchedulersRemoved++
	}

	scriptReply, err := conn.Run("/system/script/print")
	if err != nil {
		return fmt.Errorf("listing scripts: %w", err)
	}
	for _, s := range scriptReply.Re {
		name := s.Map["name"]
		isPerUser := usernames[name]
		isTransactionLog := s.Map["comment"] == "mikhmon"
		if !isPerUser && !isTransactionLog {
			continue
		}
		if _, err := conn.Run("/system/script/remove", "=.id="+s.Map[".id"]); err != nil {
			return fmt.Errorf("removing script %s: %w", name, err)
		}
		result.ScriptsRemoved++
	}

	return nil
}

// mikhmonExpiryToEpoch parses Mikhmon's "YYYY-MM-DD HH:MM:SS" expiry comment
// (written in the router's local time) into a Unix epoch (UTC). gmtOffset is
// the router's configured offset from /system/clock's gmt-offset field, used
// to convert correctly regardless of the router's timezone.
func mikhmonExpiryToEpoch(comment string, gmtOffset time.Duration) (epoch int64, ok bool) {
	m := mikhmonExpiryRe.FindStringSubmatch(comment)
	if m == nil {
		return 0, false
	}
	year, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])
	hour, _ := strconv.Atoi(m[4])
	min, _ := strconv.Atoi(m[5])
	sec, _ := strconv.Atoi(m[6])
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.UTC).Add(-gmtOffset)
	return t.Unix(), true
}

// parseGMTOffset parses RouterOS's gmt-offset field, e.g. "+01:00" or "-05:30".
func parseGMTOffset(s string) time.Duration {
	if len(s) < 6 {
		return 0
	}
	sign := 1
	if s[0] == '-' {
		sign = -1
	}
	hours, err1 := strconv.Atoi(s[1:3])
	mins, err2 := strconv.Atoi(s[4:6])
	if err1 != nil || err2 != nil {
		return 0
	}
	return time.Duration(sign) * (time.Duration(hours)*time.Hour + time.Duration(mins)*time.Minute)
}

// mikhmonProfileValiditySecs extracts the validity duration from a Mikhmon
// on-login diagnostic string, e.g. `:put (",rem,50,2h,50,,Disable,")` -> 2h.
// The no-expiry marker `:put (",,0,,,noexp,Disable,")` returns (0, true) so
// callers can distinguish "convert to no-expiry" from "not Mikhmon's format".
func mikhmonProfileValiditySecs(onLogin string) (secs int64, ok bool) {
	if onLogin == "" {
		return 0, false
	}
	if strings.Contains(onLogin, ",noexp,") {
		return 0, true
	}
	m := mikhmonOnLoginValidityRe.FindStringSubmatch(onLogin)
	if m == nil {
		return 0, false
	}
	secs, err := parseRouterOSDuration(m[1])
	if err != nil {
		return 0, false
	}
	return secs, true
}

// parseRouterOSDuration parses durations like "2h", "24h", "7d", "14d", "30d".
func parseRouterOSDuration(s string) (int64, error) {
	if len(s) < 2 {
		return 0, fmt.Errorf("invalid duration %q", s)
	}
	unit := s[len(s)-1]
	n, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid duration %q: %w", s, err)
	}
	switch unit {
	case 's':
		return n, nil
	case 'm':
		return n * 60, nil
	case 'h':
		return n * 3600, nil
	case 'd':
		return n * 86400, nil
	case 'w':
		return n * 7 * 86400, nil
	default:
		return 0, fmt.Errorf("unknown duration unit in %q", s)
	}
}
