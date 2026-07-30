# Building Pikro from source

## Requirements

| Tool | Version | Install |
|------|---------|---------|
| Go   | 1.22+   | [golang.org/dl](https://golang.org/dl) |
| Node | 18+     | [nodejs.org](https://nodejs.org) |
| npm  | 9+      | bundled with Node |

---

## Running in development

Two processes run side by side — Go backend and Vite dev server.
Vite proxies `/api` requests to the Go backend automatically.

```bash
make dev
```

Starts both the Go backend (`:8080`) and the Vite dev server (`:5173`) together.
`Ctrl+C` kills both cleanly.

Alternatively in two terminals:

```bash
# Terminal 1 — Go backend
go run .

# Terminal 2 — Vite dev server with hot-reload
cd web && npm run dev
```

Open `http://localhost:5173` in your browser.

---

## Building for production

```bash
make build
```

Builds the Vue app then compiles the Go binary with the frontend embedded.
The resulting `pikro` binary (~9 MB) contains everything and needs no separate
web server.

---

## Cross-platform release builds

```bash
make release
```

Produces binaries in `./dist/` via Go's native cross-compilation — no external
toolchain needed:

| File | Platform |
|------|----------|
| `pikro-mac-arm64` | macOS Apple Silicon |
| `pikro-mac-intel` | macOS Intel |
| `pikro.exe` | Windows 64-bit |
| `pikro-linux-amd64` | Linux 64-bit (x86) |
| `pikro-linux-arm64` | Linux 64-bit (ARM, e.g. Raspberry Pi) |

### Platform requirements

| Platform | Minimum version |
|----------|----------------|
| Windows  | Windows 10 (64-bit) |
| macOS    | macOS 11 Big Sur or later |
| Linux    | Kernel 2.6.32+ (any modern distro) |

> **Windows 7/8 are not supported.** Go 1.21+ dropped support for those versions.

---

## Makefile reference

| Command | Description |
|---------|-------------|
| `make dev` | Start Go backend + Vite dev server together (Ctrl+C kills both) |
| `make backend` | Start only the Go backend |
| `make build` | Build Vue then Go binary |
| `make release` | Cross-compile for all platforms |
| `make clean` | Remove binary and dist files |

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

### Credential storage

Passwords are stored in plain text in `~/.pikro.json`, written with owner-only
permissions (`0600`). Pikro runs on your own machine and never sends passwords
to the browser — the web UI receives router profiles with the password field
stripped. If your machine is shared, treat `~/.pikro.json` like an SSH key.

---

## macOS / Windows tray behavior

On macOS and Windows, Pikro runs as a tray application. Starting it places an
icon in the menu bar (macOS) or system tray (Windows). Click **Open Pikro** to
open the app in your browser, or **Quit** to stop the server.

A firewall prompt may appear on first run — click **Allow** so Pikro can reach
your router over the local network.
