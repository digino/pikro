<p align="center">
  <img src="assets/brand/pikro-logo-main.png" alt="Pikro" width="120" />
</p>

<h1 align="center">Pikro</h1>
<p align="center">A modern, open-source MikroTik hotspot manager</p>

<p align="center">
  Manage MikroTik RouterOS hotspots and vouchers from one zero-install app.<br/>
</p>

---

## What it does

- **Hotspot voucher management** — generate, print, enable/disable, and bulk-manage users and vouchers, no CLI required
- **Vouchers that actually expire** — installs a cleanup scheduler on the router, fixing the well-known RouterOS v7 expiry bug
- **Guided hotspot setup wizard** — provisions IP pool, DHCP, NAT, walled garden, DNS, login page, and cleanup in one flow
- **User & profile management** — create profiles with validity, price, bandwidth, and data-quota controls
- **Live monitoring** — active sessions, bandwidth, WAN IP, and system resources
- **Automatic router discovery** — finds MikroTik devices on your LAN via MNDP
- **Multi-platform** — single binary for macOS, Windows, and Linux (incl. Raspberry Pi)
- **Local-first** — runs on your machine, no telemetry, credentials never leave your computer

> A free, open-source alternative to Mikhmon — built for ISPs, café/hotel operators, and community WiFi admins who sell or hand out WiFi vouchers and want a simple, reliable tool in markets where MikroTik is dominant.

---

## Quick start

Download the binary for your platform from [Releases](../../releases) and run it:

On **macOS**, unzip and drag `Pikro.app` to Applications. A tray icon appears in the menu bar.

On **Windows**, double-click `pikro.exe`. A tray icon appears in the system tray and your browser opens automatically.

```bash
./pikro
# Pikro running at http://localhost:8080
```

---

## Docs

- [Architecture](docs/architecture.md)
- [Building from source](docs/building.md)
