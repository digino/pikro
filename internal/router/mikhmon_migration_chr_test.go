//go:build chr

// End-to-end test of the Mikhmon v3 migration on a RouterOS test router
// (CHR). It writes plans, vouchers and scripts, so it only runs against a
// router on this machine, never one on the network:
//
//	make chr6   # or chr7, in another terminal
//	PIKRO_CHR_ADDR=127.0.0.1:28728 PIKRO_CHR_PASS=<admin password> \
//	  go test -tags chr -run TestMikhmonMigrationCHR -v ./internal/router
//
// Worth running on each RouterOS family before a release: 6.x, 7.0–7.9,
// 7.10–7.11 and 7.12+ write dates differently.

package router

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
	"time"
)

// mikhmonOnLogin and mikhmonMonitor are Mikhmon v3's own scripts, from
// laksa19/mikhmonv3 hotspot/adduserprofile.php (0743da9).
func mikhmonOnLogin(expmode, price, validity, lock, plan string) string {
	lockScript := ""
	if lock == "Enable" {
		lockScript = `; [:local mac $"mac-address"; /ip hotspot user set mac-address=$mac [find where name=$user]]`
	}
	record := `; :local mac $"mac-address"; :local time [/system clock get time ]; /system script add name="$date-|-$time-|-$user-|-` + price + `-|-$address-|-$mac-|-` + validity + `-|-` + plan + `-|-$comment" owner="$month$year" source="$date" comment="mikhmon"`
	onlogin := `:put (",` + expmode + `,` + price + `,` + validity + `,,,` + lock + `,"); {:local comment [ /ip hotspot user get [/ip hotspot user find where name="$user"] comment]; :local ucode [:pic $comment 0 2]; :if ($ucode = "vc" or $ucode = "up" or $comment = "") do={ :local date [ /system clock get date ];:local year [ :pick $date 7 11 ];:local month [ :pick $date 0 3 ]; /sys sch add name="$user" disable=no start-date=$date interval="` + validity + `"; :delay 5s; :local exp [ /sys sch get [ /sys sch find where name="$user" ] next-run]; :local getxp [len $exp]; :if ($getxp = 15) do={ :local d [:pic $exp 0 6]; :local t [:pic $exp 7 16]; :local s ("/"); :local exp ("$d$s$year $t"); /ip hotspot user set comment="$exp" [find where name="$user"];}; :if ($getxp = 8) do={ /ip hotspot user set comment="$date $exp" [find where name="$user"];}; :if ($getxp > 15) do={ /ip hotspot user set comment="$exp" [find where name="$user"];};:delay 5s; /sys sch remove [find where name="$user"]`
	switch expmode {
	case "rem", "ntf":
		return onlogin + lockScript + "}}"
	case "remc", "ntfc":
		return onlogin + record + lockScript + "}}"
	}
	return `:put (",,` + price + `,,,noexp,` + lock + `,")` + lockScript
}

func mikhmonMonitor(plan, expmode string) string {
	mode := "remove"
	if expmode == "ntf" || expmode == "ntfc" {
		mode = "set limit-uptime=1s"
	}
	return `:local dateint do={:local montharray ( "jan","feb","mar","apr","may","jun","jul","aug","sep","oct","nov","dec" );:local days [ :pick $d 4 6 ];:local month [ :pick $d 0 3 ];:local year [ :pick $d 7 11 ];:local monthint ([ :find $montharray $month]);:local month ($monthint + 1);:if ( [len $month] = 1) do={:local zero ("0");:return [:tonum ("$year$zero$month$days")];} else={:return [:tonum ("$year$month$days")];}}; :local timeint do={ :local hours [ :pick $t 0 2 ]; :local minutes [ :pick $t 3 5 ]; :return ($hours * 60 + $minutes) ; }; :local date [ /system clock get date ]; :local time [ /system clock get time ]; :local today [$dateint d=$date] ; :local curtime [$timeint t=$time] ; :foreach i in [ /ip hotspot user find where profile="` + plan + `" ] do={ :local comment [ /ip hotspot user get $i comment]; :local name [ /ip hotspot user get $i name]; :local gettime [:pic $comment 12 20]; :if ([:pic $comment 3] = "/" and [:pic $comment 6] = "/") do={:local expd [$dateint d=$comment] ; :local expt [$timeint t=$gettime] ; :if (($expd < $today and $expt < $curtime) or ($expd < $today and $expt > $curtime) or ($expd = $today and $expt < $curtime)) do={ [ /ip hotspot user ` + mode + ` $i ]; [ /ip hotspot active remove [find where user=$name] ];}}}`
}

func TestMikhmonMigrationCHR(t *testing.T) {
	addr := os.Getenv("PIKRO_CHR_ADDR")
	if addr == "" {
		t.Skip("PIKRO_CHR_ADDR not set")
	}
	host, portStr, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatal(err)
	}
	if host != "127.0.0.1" && host != "localhost" {
		t.Fatalf("%s is not a router on this machine: this test rewrites the router it runs on", addr)
	}
	var port int
	fmt.Sscan(portStr, &port)
	c := NewClient(host, port, "admin", os.Getenv("PIKRO_CHR_PASS"), false)
	for i := 0; c.Ping() != nil; i++ {
		if i > 100 {
			t.Fatal("router not up")
		}
		time.Sleep(3 * time.Second)
	}
	conn, err := c.connect()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	must := func(args ...string) {
		if _, err := conn.RunArgs(args); err != nil {
			t.Fatalf("%v: %v", args[:2], err)
		}
	}
	get := func(path, name string) map[string]string {
		r, _ := conn.Run(path+"/print", "?name="+name)
		if len(r.Re) == 0 {
			return nil
		}
		return r.Re[0].Map
	}
	// A previous run's plans, vouchers and scripts all start with "mk".
	for _, path := range []string{"/ip/hotspot/user", "/ip/hotspot/user/profile", "/system/scheduler", "/system/script"} {
		r, _ := conn.Run(path + "/print")
		for _, s := range r.Re {
			if strings.HasPrefix(s.Map["name"], "mk") || s.Map["comment"] == "mikhmon" {
				conn.RunArgs([]string{path + "/remove", "=.id=" + s.Map[".id"]})
			}
		}
	}

	res, _ := conn.Run("/system/resource/print")
	version := res.Re[0].Map["version"]
	clk, _ := conn.Run("/system/clock/print")
	date := clk.Re[0].Map["date"]
	offset := parseGMTOffset(clk.Re[0].Map["gmt-offset"])
	t.Logf("RouterOS %s, clock %s %s", version, date, clk.Re[0].Map["time"])

	// ── Mikhmon's plans and their monitors ──────────────────────────────────
	plans := []struct {
		name, expmode, price, validity, lock string
		secs                                 int64
	}{
		{"mk2m", "rem", "100", "2m", "Disable", 120},
		{"mk1d", "ntfc", "500", "1d", "Enable", 86400},
		{"mk1d12h", "remc", "700", "1d12h", "Disable", 129600},
		{"mk1w", "ntf", "2000", "1w", "Disable", 604800},
		{"mk40d", "rem", "9000", "40d", "Disable", 3456000},
		{"mknoexp", "0", "200", "", "Disable", 0},
	}
	for _, p := range plans {
		must("/ip/hotspot/user/profile/add", "=name="+p.name, "=shared-users=1", "=on-login="+mikhmonOnLogin(p.expmode, p.price, p.validity, p.lock, p.name))
		if p.expmode != "0" {
			// Mikhmon runs these every ~2 min; a day keeps them from acting mid-test.
			must("/system/scheduler/add", "=name="+p.name, "=interval=1d", "=start-time=03:00:00", "=on-event="+mikhmonMonitor(p.name, p.expmode), "=comment=Monitor Profile "+p.name)
		}
	}
	must("/system/scheduler/add", "=name=mk-admin-backup", "=interval=1d", "=on-event=:log info backup")

	// ── Vouchers, as Mikhmon's generator writes them ────────────────────────
	gen := time.Now().Format("01.02.06")
	for _, p := range plans {
		must("/ip/hotspot/user/add", "=name="+p.name+"-used", "=password=x", "=profile="+p.name, "=comment=vc-123-"+gen+"-")
		must("/ip/hotspot/user/add", "=name="+p.name+"-new", "=password=x", "=profile="+p.name, "=comment=up-456-"+gen+"-batch A")
	}
	must("/ip/hotspot/user/add", "=name=mk-staff", "=password=x", "=profile=default", "=comment=reception PC")
	// An expired voucher Mikhmon's notice mode already cut off, in the date
	// form Mikhmon wrote on this version.
	week := time.Now().UTC().Add(offset).Add(-7 * 24 * time.Hour)
	past := strings.ToLower(week.Format("Jan/02/2006 15:04:05"))
	if isoDate.MatchString(date) {
		past = week.Format("2006-01-02 15:04:05")
	}
	must("/ip/hotspot/user/add", "=name=mk1w-cutoff", "=password=x", "=profile=mk1w", "=comment="+past, "=limit-uptime=1s", "=disabled=yes")

	// ── Logins: the plan's on-login runs with $user, $address and $mac ──────
	login := func(user, plan string) time.Time {
		p := get("/ip/hotspot/user/profile", plan)
		name := fmt.Sprintf("mklogin-%d", time.Now().UnixNano()%1000000)
		src := `:local user "` + user + `"; :local address "192.168.88.50"; :local "mac-address" "AA:BB:CC:00:00:01";` + "\n" + p["on-login"]
		must("/system/script/add", "=name="+name, "=source="+src)
		at := time.Now()
		if _, err := conn.RunArgs([]string{"/system/script/run", "=number=" + name}); err != nil {
			t.Errorf("login %s: %v", user, err)
		}
		must("/system/script/remove", "=numbers="+name)
		return at
	}
	loggedIn := map[string]time.Time{}
	for _, p := range plans {
		loggedIn[p.name+"-used"] = login(p.name+"-used", p.name)
	}
	// A per-user scheduler Mikhmon left behind (its login was interrupted).
	must("/system/scheduler/add", "=name=mk1w-used", "=interval=1w", "=on-event=:log info x")
	mikhmonComment := map[string]string{}
	for _, p := range plans {
		u := get("/ip/hotspot/user", p.name+"-used")
		mikhmonComment[p.name+"-used"] = u["comment"]
		t.Logf("Mikhmon stamped %-13s %q", p.name+"-used", u["comment"])
	}

	// ── Migrate ─────────────────────────────────────────────────────────────
	result, err := c.MigrateFromMikhmon()
	if err != nil {
		t.Fatalf("migration: %v", err)
	}
	t.Logf("migration: %+v", *result)

	expiry := func(comment string) (time.Time, bool) {
		if !strings.HasPrefix(comment, "exp:") {
			return time.Time{}, false
		}
		s := strings.TrimPrefix(comment, "exp:")
		if hasTimestamp(version) {
			var sec int64
			if _, err := fmt.Sscan(s, &sec); err != nil {
				return time.Time{}, false
			}
			return time.Unix(sec, 0), true
		}
		lt, err := time.Parse("2006-01-02 15:04:05", s)
		if err != nil {
			return time.Time{}, false
		}
		return lt.Add(-offset), true
	}
	for _, p := range plans {
		used := get("/ip/hotspot/user", p.name+"-used")
		switch {
		case used == nil:
			t.Errorf("%s-used disappeared", p.name)
		case p.secs == 0:
			// Mikhmon never stamps a no-expiry plan: the voucher is back to unused.
			if used["comment"] != "" {
				t.Errorf("%s-used (no expiry) comment %q, want it cleared", p.name, used["comment"])
			}
		case strings.HasPrefix(mikhmonComment[p.name+"-used"], "vc-"):
			// Mikhmon couldn't stamp it (7.10–7.11 past midnight): it counts as
			// unused, and its validity starts at its next login.
			t.Logf("%s-used never stamped by Mikhmon on this version", p.name)
		default:
			if exp, ok := expiry(used["comment"]); !ok {
				t.Errorf("%s-used not converted: %q (Mikhmon wrote %q)", p.name, used["comment"], mikhmonComment[p.name+"-used"])
			} else if want := loggedIn[p.name+"-used"].Add(time.Duration(p.secs) * time.Second); exp.Sub(want).Abs() > 20*time.Second {
				t.Errorf("%s-used expiry %v, want about %v (login + %s)", p.name, exp.UTC(), want.UTC(), p.validity)
			}
		}
		if p.lock == "Enable" && used != nil && used["mac-address"] != "AA:BB:CC:00:00:01" {
			t.Errorf("%s-used lost its device lock: mac=%q", p.name, used["mac-address"])
		}
		if n := get("/ip/hotspot/user", p.name+"-new"); n == nil || n["comment"] != "" {
			t.Errorf("%s-new unused voucher: %v", p.name, n)
		}
		prof := get("/ip/hotspot/user/profile", p.name)
		if secs, ok := pikroOnLoginValiditySecs(prof["on-login"]); p.secs > 0 && (!ok || secs != p.secs) {
			t.Errorf("plan %s: Pikro on-login validity %d (ok=%v), want %d", p.name, secs, ok, p.secs)
		} else if p.secs == 0 && prof["on-login"] != "" {
			t.Errorf("plan %s (no expiry) still has an on-login", p.name)
		}
		if p.expmode != "0" && get("/system/scheduler", p.name) != nil {
			t.Errorf("monitor of %s still there", p.name)
		}
	}
	if get("/system/scheduler", "mk1w-used") != nil {
		t.Error("Mikhmon's leftover per-user scheduler still there")
	}
	if get("/system/scheduler", "mk-admin-backup") == nil {
		t.Error("the admin's own scheduler was removed")
	}
	if left, _ := conn.Run("/system/script/print", "?comment=mikhmon"); len(left.Re) > 0 {
		t.Errorf("%d Mikhmon sales records still there", len(left.Re))
	}
	if s := get("/ip/hotspot/user", "mk-staff"); s == nil || s["comment"] != "reception PC" {
		t.Errorf("staff user changed: %v", s)
	}
	if get("/system/scheduler", cleanupScriptName) == nil {
		t.Error("cleanup not installed")
	}

	// ── Again: nothing left to change ──────────────────────────────────────
	again, err := c.MigrateFromMikhmon()
	if err != nil || again.UsersConverted+again.UsersUnused+again.ProfilesConverted+again.SchedulersRemoved+again.ScriptsRemoved > 0 {
		t.Errorf("second migration changed things: %+v err=%v", again, err)
	}

	// ── A first login after migration stamps Pikro's format ────────────────
	login("mk1d-new", "mk1d")
	if exp, ok := expiry(get("/ip/hotspot/user", "mk1d-new")["comment"]); !ok || exp.Sub(time.Now().Add(24*time.Hour)).Abs() > time.Minute {
		t.Errorf("first Pikro login on a migrated voucher: %q", get("/ip/hotspot/user", "mk1d-new")["comment"])
	}

	// ── Pikro's cleanup once the 2-minute voucher has expired ───────────────
	wait := time.Until(loggedIn["mk2m-used"].Add(3*time.Minute + 5*time.Second))
	t.Logf("waiting %v for mk2m-used to expire", wait.Round(time.Second))
	time.Sleep(wait)
	must("/system/script/add", "=name=mk-cleanup-now", "=source="+cleanupScriptBody(version))
	if _, err := conn.RunArgs([]string{"/system/script/run", "=number=mk-cleanup-now"}); err != nil {
		t.Errorf("cleanup run: %v", err)
	}
	for name, wantGone := range map[string]bool{"mk2m-used": true, "mk1w-cutoff": true, "mk1d-used": false, "mk40d-used": false, "mknoexp-used": false, "mk1d-new": false, "mk2m-new": false, "mk-staff": false, "default-trial": false} {
		if gone := get("/ip/hotspot/user", name) == nil; gone != wantGone {
			t.Errorf("after cleanup %s gone=%v, want %v", name, gone, wantGone)
		}
	}
}
