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

Expiry is managed through two RouterOS mechanisms:

1. **`on-login` script on the hotspot user profile** — fires on first login, writes `exp:<unix_epoch>` into the user's comment. The validity window starts from first use, not from voucher creation.

2. **`pikro-cleanup` scheduler** — runs on a configurable interval, removes users whose `exp:` epoch has passed.

The script variant is selected based on the RouterOS version detected at install time:
- **v7.12+**: uses `[:tonsec [:timestamp]]` for Unix epoch comparison
- **v6 / early v7**: uses Mikhmon-style `dateint`/`timeint` inline functions with date string comparison

## Versioning

The version string is injected at build time via `-ldflags "-X main.Version=v1.0.0"`.
The Makefile reads it from the latest git tag via `git describe --tags --always --dirty`.
The frontend fetches it from `GET /api/version` on startup.
