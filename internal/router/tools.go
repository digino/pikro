package router

// The Troubleshoot and Maintenance pages: checks and fixes run on the router.

import (
	"errors"
	"fmt"
	"math"
	"net/netip"
	"regexp"
	"strconv"
	"strings"
	"time"

	ros "github.com/go-routeros/routeros/v3"
)

const (
	internetProbe = "1.1.1.1"
	dnsProbe      = "google.com"
	ntpServer     = "pool.ntp.org"
)

// CheckStep is one link in the chain from the router to the internet.
type CheckStep struct {
	Key    string `json:"key"`              // cable | gateway | internet | dns
	State  string `json:"state"`            // ok | fail | skip
	Target string `json:"target"`           // the interface, address or name checked
	Reason string `json:"reason,omitempty"` // noRoute | noLink, for the cable step
	RTTMs  int    `json:"rttMs,omitempty"`
	Via    string `json:"via,omitempty"` // DNS servers asked
}

type InternetCheck struct {
	CheckedAt time.Time   `json:"checkedAt"`
	Steps     []CheckStep `json:"steps"`
}

// InternetCheck walks from the router's internet port outwards. The first
// step that fails is the problem; the steps after it are skipped, since
// they can only fail too.
func (c *Client) InternetCheck() (*InternetCheck, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	routes, err := conn.Run("/ip/route/print", "?dst-address=0.0.0.0/0")
	if err != nil {
		return nil, err
	}
	var route map[string]string
	for _, r := range reMaps(routes) {
		if r["active"] == "true" {
			route = r
			break
		}
	}
	wan := wanInterface(reMaps(routes))

	steps := []CheckStep{
		{Key: "cable", Target: wan},
		{Key: "gateway"},
		{Key: "internet", Target: internetProbe},
		{Key: "dns", Target: dnsProbe},
	}
	failed := false
	// run marks step i from check(), or skips it once an earlier step failed.
	run := func(i int, check func() bool) {
		if failed {
			steps[i].State = "skip"
			return
		}
		if check() {
			steps[i].State = "ok"
		} else {
			steps[i].State = "fail"
			failed = true
		}
	}

	run(0, func() bool {
		if wan == "" {
			steps[0].Reason = "noRoute"
			return false
		}
		r, err := conn.Run("/interface/print", "?name="+wan)
		if err != nil || len(r.Re) == 0 || r.Re[0].Map["running"] != "true" {
			steps[0].Reason = "noLink"
			return false
		}
		return true
	})
	gw := gatewayAddress(route)
	steps[1].Target = gw
	run(1, func() bool {
		if gw == "" {
			// PPPoE and similar: the link itself is the gateway, checked above.
			steps[1].Target = wan
			return true
		}
		received, rtt := c.ping(conn, gw)
		steps[1].RTTMs = rtt
		return received > 0
	})
	run(2, func() bool {
		received, rtt := c.ping(conn, internetProbe)
		steps[2].RTTMs = rtt
		return received > 0
	})
	if r, err := conn.Run("/ip/dns/print"); err == nil && len(r.Re) > 0 {
		dns := r.Re[0].Map
		steps[3].Via = dns["servers"]
		if steps[3].Via == "" {
			steps[3].Via = dns["dynamic-servers"]
		}
	}
	run(3, func() bool {
		_, err := conn.Run("/resolve", "=domain-name="+dnsProbe)
		return err == nil
	})

	return &InternetCheck{CheckedAt: time.Now().UTC(), Steps: steps}, nil
}

func (c *Client) ping(conn *ros.Client, address string) (received, avgMs int) {
	r, err := conn.Run("/ping", "=address="+address, "=count=3")
	if err != nil {
		return 0, 0
	}
	return pingSummary(reMaps(r))
}

// pingSummary reads the running totals RouterOS adds to every ping reply.
func pingSummary(replies []map[string]string) (received, avgMs int) {
	if len(replies) == 0 {
		return 0, 0
	}
	last := replies[len(replies)-1]
	received, _ = strconv.Atoi(last["received"])
	if received == 0 {
		return 0, 0
	}
	// Nearest millisecond, but an answer is never shown as 0 ms.
	return received, max(1, int(math.Round(float64(rttMicros(last["avg-rtt"]))/1000)))
}

var rttPart = regexp.MustCompile(`(\d+)(ms|us|s)`)

// rttMicros reads RouterOS round-trip times: "5ms" (v6), "12ms345us" (v7).
func rttMicros(s string) int {
	total := 0
	for _, m := range rttPart.FindAllStringSubmatch(s, -1) {
		n, _ := strconv.Atoi(m[1])
		switch m[2] {
		case "s":
			total += n * 1_000_000
		case "ms":
			total += n * 1000
		default:
			total += n
		}
	}
	return total
}

// gatewayAddress is the IP the default route goes through, or "" when the
// gateway is an interface (PPPoE).
func gatewayAddress(route map[string]string) string {
	gw, _, _ := strings.Cut(route["gateway"], ",")
	gw, _, _ = strings.Cut(gw, "%")
	if _, err := netip.ParseAddr(gw); err != nil {
		return ""
	}
	return gw
}

// UsePublicDNS points the router at public DNS servers, for when the
// provider's own servers stop answering.
func (c *Client) UsePublicDNS() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Run("/ip/dns/set", "=servers=1.1.1.1,8.8.8.8")
	return err
}

type RouterClock struct {
	UTC        time.Time `json:"utc"`
	Local      string    `json:"local"` // the router's own wall time, "2026-10-29T14:02:00"
	NTPEnabled bool      `json:"ntpEnabled"`
	NTPServers string    `json:"ntpServers"`
}

func (c *Client) Clock() (*RouterClock, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	r, err := conn.Run("/system/clock/print")
	if err != nil || len(r.Re) == 0 {
		return nil, fmt.Errorf("read clock: %w", err)
	}
	m := r.Re[0].Map
	utc, err := routerClockUTC(m["date"], m["time"], m["gmt-offset"])
	if err != nil {
		return nil, err
	}
	out := &RouterClock{
		UTC:   utc,
		Local: utc.Add(offsetDuration(m["gmt-offset"])).Format("2006-01-02T15:04:05"),
	}
	if r, err := conn.Run("/system/ntp/client/print"); err == nil && len(r.Re) > 0 {
		n := r.Re[0].Map
		out.NTPEnabled = n["enabled"] == "true"
		// RouterOS 7 lists "servers"; RouterOS 6 "server-dns-names" or a primary IP.
		out.NTPServers = n["servers"]
		if out.NTPServers == "" {
			out.NTPServers = n["server-dns-names"]
		}
		if out.NTPServers == "" && n["primary-ntp"] != "0.0.0.0" {
			out.NTPServers = n["primary-ntp"]
		}
	}
	return out, nil
}

// EnableAutoTime turns on the NTP client and sets the clock to now right
// away, since the first NTP sync can take a while.
func (c *Client) EnableAutoTime(now time.Time) error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	r, err := conn.Run("/system/ntp/client/print")
	if err != nil || len(r.Re) == 0 {
		return fmt.Errorf("read NTP client: %w", err)
	}
	if _, v7 := r.Re[0].Map["servers"]; v7 {
		_, err = conn.Run("/system/ntp/client/set", "=enabled=yes", "=servers="+ntpServer)
	} else {
		// RouterOS 6 refuses "mode" over the API but switches to unicast by
		// itself once it has a server name.
		_, err = conn.Run("/system/ntp/client/set", "=enabled=yes", "=server-dns-names="+ntpServer)
	}
	if err != nil {
		return err
	}

	clock, err := conn.Run("/system/clock/print")
	if err != nil || len(clock.Re) == 0 {
		return fmt.Errorf("read clock: %w", err)
	}
	m := clock.Re[0].Map
	date, tod := routerClockStrings(now, m["gmt-offset"], isoDate.MatchString(m["date"]))
	_, err = conn.Run("/system/clock/set", "=date="+date, "=time="+tod)
	return err
}

var (
	isoDate     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	clockOffset = regexp.MustCompile(`^([+-])(\d{2}):(\d{2})$`)
)

// routerClockUTC reads the router's wall clock: RouterOS 7.10+ writes the
// date as 2026-10-29, older versions as oct/29/2026.
func routerClockUTC(date, tod, offset string) (time.Time, error) {
	layout := "Jan/02/2006 15:04:05"
	if isoDate.MatchString(date) {
		layout = "2006-01-02 15:04:05"
	}
	local, err := time.Parse(layout, date+" "+tod)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown router date %q: %w", date, err)
	}
	return local.Add(-offsetDuration(offset)), nil
}

// routerClockStrings formats a moment as the date and time /system/clock/set
// takes, in the router's time zone and date style.
func routerClockStrings(utc time.Time, offset string, iso bool) (date, tod string) {
	local := utc.Add(offsetDuration(offset))
	if iso {
		date = local.Format("2006-01-02")
	} else {
		date = strings.ToLower(local.Format("Jan/02/2006"))
	}
	return date, local.Format("15:04:05")
}

func offsetDuration(offset string) time.Duration {
	m := clockOffset.FindStringSubmatch(offset)
	if m == nil {
		return 0
	}
	h, _ := strconv.Atoi(m[2])
	min, _ := strconv.Atoi(m[3])
	d := time.Duration(h)*time.Hour + time.Duration(min)*time.Minute
	if m[1] == "-" {
		return -d
	}
	return d
}

type UpdateInfo struct {
	Installed string `json:"installed"`
	Latest    string `json:"latest"`
	Channel   string `json:"channel"`
	Status    string `json:"status"`
	Available bool   `json:"available"`
}

// CheckUpdates asks MikroTik for the latest RouterOS on the router's channel.
// It only checks: installing restarts the router.
func (c *Client) CheckUpdates() (*UpdateInfo, error) {
	conn, err := c.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	r, err := conn.Run("/system/package/update/check-for-updates")
	if err != nil {
		return nil, err
	}
	replies := reMaps(r)
	if len(replies) == 0 {
		return nil, errors.New("no answer from the update check")
	}
	// The last reply holds the outcome; earlier ones are progress.
	last := replies[len(replies)-1]
	info := &UpdateInfo{
		Installed: last["installed-version"],
		Latest:    last["latest-version"],
		Channel:   last["channel"],
		Status:    last["status"],
	}
	if info.Latest == "" {
		return nil, fmt.Errorf("update check failed: %s", info.Status)
	}
	info.Available = info.Latest != info.Installed
	return info, nil
}

// Reboot restarts the router. The router drops the connection as it goes
// down, so only an error the router itself sends back counts as failure.
func (c *Client) Reboot() error {
	conn, err := c.connect()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Run("/system/reboot")
	var deviceErr *ros.DeviceError
	if errors.As(err, &deviceErr) {
		return err
	}
	return nil
}

func reMaps(r *ros.Reply) []map[string]string {
	out := make([]map[string]string, 0, len(r.Re))
	for _, s := range r.Re {
		out = append(out, s.Map)
	}
	return out
}
