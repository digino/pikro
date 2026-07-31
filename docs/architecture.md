# Architecture

## Project structure

```
pikro/
├── main.go                          # Entry point
├── internal/
│   ├── assets/                      # Embedded hotspot HTML pages
│   │   ├── assets.go                # go:embed declarations
│   │   ├── hotspot_login.html       # Go template — customizable via Settings
│   │   ├── hotspot_logout.html
│   │   ├── hotspot_status.html
│   │   ├── hotspot_alogin.html
│   │   ├── hotspot_error.html
│   │   └── hotspot_redirect.html
│   ├── config/config.go             # Router profiles → ~/.pikro.json
│   ├── discovery/                   # MNDP broadcast scanner (UDP :5678)
│   │   ├── mndp.go                  # Packet encoding/decoding
│   │   ├── scan_unix.go             # Linux/macOS raw socket implementation
│   │   └── scan_windows.go          # Windows net.ListenPacket implementation
│   ├── handlers/                    # HTTP handler functions
│   │   ├── routers.go               # CRUD for router profiles
│   │   ├── hotspot.go               # Hotspot user management
│   │   ├── hotspot_profiles.go      # Hotspot user profiles
│   │   ├── hotspot_settings.go      # Hotspot settings + login page upload
│   │   ├── hotspot_setup.go         # Guided hotspot setup wizard
│   │   ├── cleanup_scheduler.go     # Auto-cleanup scheduler management
│   │   ├── system.go                # System resource metrics + app version
│   │   ├── speedtest.go             # Bandwidth test trigger
│   │   ├── discovery.go             # Network scan endpoint
│   │   └── helpers.go               # JSON response helpers
│   ├── router/                      # RouterOS native API client
│   │   ├── client.go                # go-routeros wrapper, cleanup scripts
│   │   ├── hotspot_setup.go         # Hotspot setup orchestration steps
│   │   └── login_page.go            # Login page rendering + file upload
│   └── server/server.go             # HTTP routes + embedded SPA handler
└── web/                             # Vue 3 frontend
    ├── src/
    │   ├── api/index.ts             # Typed API client (axios)
    │   ├── stores/routers.ts        # Pinia: router profiles + UI state
    │   ├── router/index.ts          # Vue Router routes
    │   ├── components/
    │   │   ├── PageLayout.vue       # Shared page wrapper
    │   │   ├── AppDialog.vue        # Reka UI modal wrapper
    │   │   ├── AddRouterDialog.vue  # Add router form + MNDP discovery UI
    │   │   └── StatCard.vue         # Metric card for dashboard
    │   └── views/
    │       ├── Dashboard.vue        # CPU, memory, uptime, hotspot summary
    │       ├── Routers.vue          # Router table with test/delete actions
    │       ├── SpeedTest.vue        # Router-side bandwidth test
    │       └── hotspot/
    │           ├── Setup.vue        # Guided hotspot setup wizard
    │           ├── Users.vue        # Hotspot users + active sessions
    │           ├── Profiles.vue     # Hotspot user profiles
    │           └── Settings.vue     # Hotspot settings, login page, vouchers
    └── dist/                        # Built output — embedded into Go binary
```

## How it works

1. **Go binary starts** → finds a free port (default 8080) → opens the browser
2. **Browser loads** the Vue SPA embedded in the binary via `//go:embed`
3. **Vue makes API calls** to `localhost:PORT/api/*`
4. **Go calls** the MikroTik router's native API (port 8728/8729)
5. **Router credentials** are stored locally in `~/.pikro.json`

```
Browser  ←→  Go (localhost:8080)  ←→  RouterOS Native API (port 8728)
              └─ serves SPA             (same protocol as WinBox)
              └─ stores config (~/.pikro.json)
              └─ MNDP scanner (UDP broadcast :5678)
```

## RouterOS Native API

Uses the **RouterOS native API protocol** — the same protocol WinBox uses, implemented via [go-routeros v3](https://github.com/go-routeros/routeros).
Works on RouterOS **v6 and v7**. No extra service needs to be enabled on the router.

| Mode  | Port     | When to use |
|-------|----------|-------------|
| Plain | **8728** | Default — same LAN, no TLS needed |
| TLS   | **8729** | Enable "Use TLS" when adding a router |

Authentication uses the same username/password as WinBox/console.
The API service is enabled by default on RouterOS; ensure port 8728 is not blocked by the firewall.

## Router discovery (MNDP)

Clicking **Scan network** sends a 4-byte UDP broadcast to `255.255.255.255:5678`.
MikroTik routers reply with a TLV-encoded packet containing identity, IP, MAC, board model and RouterOS version.
Discovery is Layer 2 — only routers on the **same subnet** are found, identical to WinBox neighbour discovery.

## Hotspot voucher expiry system

A voucher has **two independent limits**, either of which ends it — confirmed against
real West-Africa operator usage: an admin can set validity=24h + uptime=1h30m (customer
gets 24h of *calendar time* to use up to 1h30m of *connected time*), or validity=24h +
uptime=24h (voucher simply can't outlast 24h regardless of usage pattern). Neither
implies the other.

### 1. Time-based validity — `exp:<epoch>` comment

- **`on-login` script on the hotspot user profile** — fires on first login, writes
  `exp:<unix_epoch>` (or `exp:YYYY-MM-DD HH:MM:SS` on the v6 script variant) into the
  user's comment. The validity window starts from **first use, not from voucher
  creation** — an unused voucher has no comment at all and shows as "Waiting" in the UI,
  by design (mirrors Mikhmon's own `up-...` unused-voucher convention; a voucher sitting
  unsold shouldn't have its clock already running).
- The validity duration is **fixed per profile**, baked into that profile's on-login
  script text at profile-creation time — it does not read the individual voucher's own
  `limit-uptime` override. This is intentional: validity and uptime are orthogonal
  limits (see above), so the on-login script only ever needs the profile's configured
  validity, never the per-voucher uptime override.
- `pikro-cleanup` scheduler removes any user whose `exp:` epoch is in the past.

### 2. Usage-based quota — `uptime` vs `limit-uptime`

- RouterOS enforces `limit-uptime` at the network level natively — once a user's
  accumulated `uptime` reaches it, RouterOS itself refuses further logins ("your uptime
  limit is reached" in `/log`). Pikro's cleanup script independently checks
  `uptime >= limit-uptime` and removes the account, regardless of the `exp:` comment —
  there's no reason to keep an account RouterOS has already cut off.
- RouterOS reports `uptime`/`limit-uptime` in **two different string shapes** depending
  on context, both of which the cleanup script's `durationToSecs` function must parse:
  letter-suffixed (`"5m"`, `"1h30m"`, `"1w3d"` — what Pikro's Go API client normalises
  values to before the frontend ever sees them) and colon-joined time-of-day, optionally
  week/day-prefixed (`"00:00:00"`, `"1d00:00:00"`, `"1w3d10:51:24"` — the **raw**
  script-level value, never visible through Pikro's own API/UI). An incident on
  2026-07-30 was caused by a duration parser that only handled the first form: a fresh,
  never-logged-in voucher's `uptime` reads as `"00:00:00"` at the script level, which the
  broken parser silently read as `0`, making `0 >= limit` always true and deleting every
  unused voucher in a batch on the very first scheduler tick.

### Both mechanisms share one RouterOS version split

- **v7.12+**: uses `[:tonsec [:timestamp]]` for Unix epoch comparison.
- **v6 / early v7**: uses Mikhmon-style `dateint`/`timeint` inline functions with date
  string comparison (proven v6-compatible pattern, reused from Mikhmon's own approach).

### Known failure mode: router clock correctness

Time-based expiry (`exp:`) only works if the router's clock is consistent between the
moment a voucher is stamped and the moment cleanup checks it — the *absolute* value
doesn't matter, only that it doesn't change in between. RouterBoard hardware has no
battery-backed clock, so every unclean reboot (frequent in outage-prone deployments —
Pikro's target market) resets the clock to RouterOS's factory default and only corrects
it via NTP or MikroTik's cloud-based time-zone detection, both of which require internet
access the router may not have. A reboot between stamping and checking an `exp:` value
breaks the comparison in either direction (wrongly-early or wrongly-late deletion). No
in-script signal reliably distinguishes "clock is correct" from "clock is wrong but
plausible-looking" without a working NTP client actually reporting `synchronized` status
— an earlier mitigation attempt (skip `exp:` checks if the clock year looks like
RouterOS's 1970 factory default) was removed because it only caught that one narrow case
and did nothing for the actual failure mode (a clock that's simply wrong after reboot,
not still on the factory default). The uptime-quota check is unaffected — RouterOS's own
`uptime` counter has no dependency on wall-clock time.

### Migrating from Mikhmon

Routers previously managed by Mikhmon carry its own comment conventions
(`up-<validitySecs>-<MM.DD.YY>-` for unused vouchers, bare `YYYY-MM-DD HH:MM:SS` for
stamped expiry) and Mikhmon-style on-login scripts, neither of which Pikro's `exp:`
parsing understands — such users show "Waiting" forever even once connected. Settings →
Migration (`SettingsMigration.vue`, backed by `Client.MigrateFromMikhmon()` in
`internal/router/migrate_mikhmon.go`) converts both formats to Pikro's own convention and
rewrites each profile's on-login script; safe to re-run.

## Versioning

The version string is injected at build time via `-ldflags "-X main.Version=v1.0.0"`.
The Makefile reads it from the latest git tag via `git describe --tags --always --dirty`.
The frontend fetches it from `GET /api/version` on startup.


-lock user as mikhmon
- The generated app should display "Open Pikro" and "Stop server" buttons with ports on which it will open clearly 
