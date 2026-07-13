# Pikro

A zero-install web app for managing MikroTik RouterOS devices.
Run a single binary — your browser opens automatically.

---

## Requirements

### End users
Nothing. Download the binary and run it.

### Developers
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
Passwords are stored in plain text — the file is readable only by the
current user (`chmod 600`).

---

## Makefile reference

| Command | Description |
|---------|-------------|
| `make dev` | Start Go backend + Vite dev server together (Ctrl+C kills both) |
| `make backend` | Start only the Go backend (when Vite is already running) |
| `make build` | Build Vue then Go binary |
| `make release` | Cross-compile for all platforms |
| `make clean` | Remove binary and dist files |