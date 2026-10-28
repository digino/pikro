# Local RouterOS test VMs (CHR)

Test Pikro against a real RouterOS without touching a production router.
CHR (Cloud Hosted Router) is MikroTik's official VM image of RouterOS. The free
licence caps throughput at 1 Mbps per interface, which does not matter here.

| Target | RouterOS | API | WinBox | Speed on Apple Silicon |
|---|---|---|---|---|
| `make chr7` | 7.25 (main target) | `127.0.0.1:18728` | `127.0.0.1:18291` | fast enough (ARM64, emulated CPU, boots in ~15 s) |
| `make chr6` | 6.49.23 (compatibility) | `127.0.0.1:28728` | `127.0.0.1:28291` | fine (x86 emulated, boots in ~20 s) |

The console runs in the terminal; quit with **Ctrl-A then X**.

## One-time setup

```sh
brew install qemu
mkdir -p dev/chr/images/v7 dev/chr/images/v6
curl -fL -o dev/chr/images/v7/chr-7.25-arm64.img.zip \
  https://download.mikrotik.com/routeros/7.25/chr-7.25-arm64.img.zip
curl -fL -o dev/chr/images/v6/chr-6.49.23.img.zip \
  https://download.mikrotik.com/routeros/6.49.23/chr-6.49.23.img.zip
```

Check each zip against MikroTik's published SHA256 (same URL plus `.sha256`).
`dev/chr/images/` is gitignored.

## Use it from Pikro

1. `make chr7` and wait for the `CHR Login:` prompt.
2. In Pikro: Routers, then **Add by IP address**: `127.0.0.1`, port `18728`,
   user `admin`, empty password (unless you set one at the console).
3. Run Pikro's hotspot setup on it using **ether2** (an isolated LAN port;
   ether1 is the uplink carrying the API), then test freely.

Network scan (MNDP) will not find the VM: it sits behind QEMU's NAT. Add it by IP.

## Test as a hotspot guest

`make chr7-guest` (or `chr6-guest`) boots the same router with a guest device
on the hotspot, so the login pages can be tested end to end without Wi-Fi:

1. Stop the normal VM first (Ctrl-A X): both use the same disk image.
2. `make chr7-guest`. Everything sent to `http://127.0.0.1:18080`
   (`28080` for v6) reaches the hotspot as the guest `192.168.88.254`.
3. A browser must use the hotspot's DNS name, because the hotspot redirects
   to it. In Chrome: `--host-resolver-rules="MAP chr7.spot:80 127.0.0.1:18080"`,
   then open `http://chr7.spot/login`.
4. Quit with Ctrl-A X and go back to `make chr7`.

The defaults match a hotspot on `192.168.88.0/24` at `192.168.88.1`. For another
subnet set `GUEST_NET`, `GUEST_HOTSPOT` and `GUEST_CLIENT`, for example
`GUEST_NET=10.5.50.0/24 GUEST_HOTSPOT=10.5.50.1 GUEST_CLIENT=10.5.50.254 make chr7-guest`.
All requests come from one guest device, so this tests one guest at a time.

## State and reset

The VM writes to the unpacked `.img`, so router state survives restarts.
Delete `dev/chr/images/v7/*.img` (or `v6`) to start from a factory-fresh router;
the next run unpacks it again from the zip.

## Limits

- No Wi-Fi clients: for logins, "Online now" and login/logout log lines use
  guest mode (one guest device, above).
- To test a different RouterOS version, download its CHR zip and update the
  file names in `run.sh`.
