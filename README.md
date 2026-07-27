# Pikro — a modern, open-source MikroTik hotspot manager

**Manage MikroTik RouterOS hotspots and vouchers from one zero-install app.**
A single binary serves a clean web UI and talks to your router over the native
RouterOS API (port 8728/8729 — the same protocol WinBox uses). No PHP, no web
server to configure, no database. Download, run, done.

> A free, open-source alternative to Mikhmon — built for café, hotel, and community
> WiFi operators who sell hotspot vouchers. **Vouchers that actually expire correctly
> on RouterOS v7**, without touching the router CLI.

---

## Features

- 🎫 **Hotspot voucher management** — generate, print, enable/disable, and bulk-manage
  hotspot users and vouchers from a clean UI, no RouterOS CLI required.
- ⏱️ **Vouchers that actually expire** — a cleanup scheduler installed on the router
  enforces expiry reliably, fixing the well-known Mikhmon-on-RouterOS-v7 issue where
  expired vouchers keep working.
- 🪄 **Guided hotspot setup wizard** — one flow provisions IP pool, DHCP, NAT, walled
  garden, DNS, login page, and the cleanup scheduler. Full teardown too.
- 👤 **User & profile management** — create and edit hotspot profiles with validity,
  price, bandwidth (rate-limit), and data-quota controls.
- 🎨 **Custom login page & voucher templates** — brand the captive portal and printed
  vouchers with your business name, colours, price, and validity.
- 📡 **Live monitoring** — active sessions, per-interface traffic, bandwidth, WAN IP,
  and system resource dashboard.
- 🔎 **Automatic router discovery** — finds MikroTik devices on your LAN via MNDP.
- 🖥️ **Multi-platform** — single ~9 MB binary for macOS, Windows, and Linux (incl.
  Raspberry Pi / ARM).
- 🔒 **Local-first & private** — runs on your machine, no cloud, no telemetry;
  credentials never leave your computer.

---

## Quick start (end users)

Nothing to install. Download the binary for your platform from
[Releases](../../releases), then run it:

```bash
./pikro
# Pikro running at http://localhost:8080  (your browser opens automatically)
```

On Windows, double-click `pikro.exe` (keep the console window open — closing it stops
the server). See [Windows behavior](#windows-behavior) for details.

---

## Who it's for

Small ISPs, café/hotel operators, and community-network admins in markets where
MikroTik is dominant — anyone who sells or hands out WiFi hotspot vouchers and wants
a simple, reliable tool instead of the router CLI or an aging PHP panel.

---

## Building from source

### Requirements

| Tool | Version | Install |
|------|---------|---------|
| Go   | 1.22+   | [golang.org/dl](https://golang.org/dl) |
| Node | 18+     | [nodejs.org](https://nodejs.org) |
| npm  | 9+      | bundled with Node |

---

## Running in development

Two processes run side by side — Go backend and Vite dev server.
Vite proxies `/api` requests to the Go backend automatically.

### One command (recommended)

```bash
make dev
```

Starts both the Go backend (`:8080`) and the Vite dev server (`:5173`) together.
`Ctrl+C` kills both cleanly.

### Two terminals (alternative)

```bash
# Terminal 1 — Go backend
go run .

# Terminal 2 — Vite dev server with hot-reload
cd web && npm run dev
```

Open `http://localhost:5173` in your browser.
The Go backend logs all actions including MNDP discovery to stdout.

---

## Building for production

```bash
# 1. Build the Vue app
cd web && npm run build

# 2. Compile the Go binary (embeds web/dist/ automatically)
cd .. && go build -o pikro .
```

Or use the Makefile shortcut:

```bash
make build
```

The resulting `pikro` binary (~9 MB) contains everything.
Run it and your browser opens automatically:

```bash
./pikro
# Pikro running at http://localhost:8080
```

---

## Cross-platform release builds

```bash
make release
```

Produces binaries in `./dist/`:

| File | Platform |
|------|----------|
| `pikro-mac-arm64` | macOS Apple Silicon |
| `pikro-mac-intel` | macOS Intel |
| `pikro.exe` | Windows 64-bit |
| `pikro-linux-amd64` | Linux 64-bit (x86) |
| `pikro-linux-arm64` | Linux 64-bit (ARM, e.g. Raspberry Pi) |

No cross-compilation toolchain needed — Go handles it natively.

### Platform requirements

| Platform | Minimum version |
|----------|----------------|
| Windows  | **Windows 10** (64-bit) |
| macOS    | macOS 11 Big Sur or later |
| Linux    | Kernel 2.6.32+ (any modern distro) |

> **Windows 7 is not supported.** The binary is built with Go 1.21+ which dropped
> Windows 7/8 support. The app requires Windows 10 or later.

### Windows behavior

On Windows, double-clicking `pikro.exe` opens a console window showing the
server address, then launches your browser automatically. **Keep the console window
open** — closing it stops the server. To run it in the background, use:

```bat
start /b pikro.exe
```

---

## Configuration

Router profiles are stored at `~/.pikro.json`:

```json
{
  "routers": [
    {
      "id": "abc-123",
      "name": "Home router",
      "host": "192.168.88.1",
      "port": 8728,
      "username": "admin",
      "password": "secret",
      "useTls": false
    }
  ]
}
```

The file is created automatically when you add your first router.

### A note on credential storage

Router passwords are stored in plain text in `~/.pikro.json`, which Pikro writes
with owner-only permissions (`0600`). This is a deliberate, local-first design:

- Pikro runs on **your own machine** and talks to routers on **your own network** —
  the threat model is "another user on this same computer," which `0600` addresses.
- RouterOS credentials are not a new secret Pikro creates; they are the same
  username/password you already type into WinBox/WebFig.
- Pikro **never sends passwords to the browser** — the web UI receives router
  profiles with the password field stripped.

If your machine is shared or higher-risk, treat `~/.pikro.json` like an SSH key.

---

## Makefile reference

| Command | Description |
|---------|-------------|
| `make dev` | Start Go backend + Vite dev server together (Ctrl+C kills both) |
| `make backend` | Start only the Go backend (when Vite is already running) |
| `make build` | Build Vue then Go binary |
| `make release` | Cross-compile for all platforms |
| `make clean` | Remove binary and dist files |