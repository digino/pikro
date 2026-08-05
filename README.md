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

- **Hotspot voucher management** — generate, print, enable/disable, edit, and bulk-manage users and vouchers, no CLI required
- **Vouchers that actually expire** — installs a cleanup scheduler on the router, fixing the well-known RouterOS v7 expiry bug
- **Guided hotspot setup wizard** — provisions IP pool, DHCP, NAT, walled garden, DNS, login page, and cleanup in one flow
- **Customizable login pages** — pick from several ready-made hotspot login page designs, set a title/subtitle, preview live, and upload straight to the router
- **User & profile management** — create profiles with validity, price, bandwidth, and data-quota controls
- **Live monitoring** — active sessions, bandwidth, WAN IP, and system resources, with one-click disconnect for active sessions
- **Automatic router discovery** — finds MikroTik devices on your LAN via MNDP
- **Multi-platform** — single binary for macOS, Windows, and Linux (incl. Raspberry Pi)
- **Local-first** — runs on your machine, no telemetry, credentials never leave your computer

> A free, open-source alternative to Mikhmon — built for ISPs, café/hotel operators, and community WiFi admins who sell or hand out WiFi vouchers and want a simple, reliable tool in markets where MikroTik is dominant.

---

## Quick start

Download the binary for your platform from [Releases](../../releases) and run it.

On macOS and Windows, Pikro runs as a tray application. Starting it places an
icon in the menu bar (macOS) or system tray (Windows). Click **Open Pikro** to
open the app in your browser, or **Quit** to stop the server. On Linux, Pikro
runs in the terminal and opens your default browser automatically — stop it
with Ctrl+C.

A firewall prompt may appear on first run — click **Allow** so Pikro can reach
your router over the local network.

```bash
./pikro
# Pikro running at http://localhost:8080
```

> **Note:** Release binaries aren't code-signed yet (see
> [Releases](../../releases) — each build carries a
> [SLSA](https://slsa.dev) provenance attestation proving it was built from
> this repo, but isn't Apple/Microsoft-signed). On first run:
> - **macOS:** right-click the app → **Open** (instead of double-clicking) to
>   get past the "unidentified developer" warning.
> - **Windows:** if SmartScreen appears, click **More info** → **Run anyway**.

---

## Docs

- [Building from source](docs/building.md)
