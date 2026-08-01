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

Download the binary for your platform from [Releases](../../releases) and run it. 

On macOS and Windows, Pikro runs as a tray application. Starting it places an
icon in the menu bar (macOS) or system tray (Windows). Click **Open Pikro** to
open the app in your browser, or **Quit** to stop the server.

A firewall prompt may appear on first run — click **Allow** so Pikro can reach
your router over the local network.

```bash
./pikro
# Pikro running at http://localhost:8080
```

---

## Docs

- [Building from source](docs/building.md)
