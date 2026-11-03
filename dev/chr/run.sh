#!/usr/bin/env bash
# Boot a local RouterOS CHR (Cloud Hosted Router) VM for testing Pikro without
# touching a production router. See dev/chr/README.md.
#
#   dev/chr/run.sh v7    RouterOS 7, ARM64 (emulated CPU, see below)
#   dev/chr/run.sh v6    RouterOS 6, x86 emulated (slow, for compatibility checks)
#   dev/chr/run.sh v7 guest   same router, with a guest device on the hotspot
#                             (see "Test as a hotspot guest" in the README)
#
# The RouterOS console runs in this terminal. Quit QEMU with Ctrl-A then X.
set -euo pipefail

cd "$(dirname "$0")"
QEMU_SHARE="$(brew --prefix qemu)/share/qemu"

case "${1:-}" in
  v7)
    ZIP=images/v7/chr-7.25-arm64.img.zip
    IMG=images/v7/chr-7.25-arm64.img
    API_PORT=18728 WINBOX_PORT=18291 GUEST_PORT=18080
    ;;
  v6)
    ZIP=images/v6/chr-6.49.23.img.zip
    IMG=images/v6/chr-6.49.23.img
    API_PORT=28728 WINBOX_PORT=28291 GUEST_PORT=28080
    ;;
  *)
    echo "usage: $0 v7|v6" >&2
    exit 1
    ;;
esac

if [ ! -f "$ZIP" ]; then
  echo "Missing $ZIP — download it as described in dev/chr/README.md" >&2
  exit 1
fi
# The VM writes to the raw image in place, so router state survives restarts.
# Delete the .img to start again from a factory-fresh router.
[ -f "$IMG" ] || unzip -o -d "$(dirname "$IMG")" "$ZIP"

# ether1: uplink with the API/WinBox ports forwarded (only localhost can reach them).
NET="user,id=n0,hostfwd=tcp:127.0.0.1:${API_PORT}-:8728,hostfwd=tcp:127.0.0.1:${WINBOX_PORT}-:8291"
# ether2: an isolated LAN port for Pikro's hotspot setup to use.
LAN="user,id=n1,restrict=on"

echo "RouterOS $1 — API 127.0.0.1:${API_PORT}, WinBox 127.0.0.1:${WINBOX_PORT} (Ctrl-A X quits)"

if [ "${2:-}" = guest ]; then
  # QEMU's user network joins ether2 to the hotspot subnet as one guest device
  # (GUEST_CLIENT). 127.0.0.1:GUEST_PORT on this Mac reaches the hotspot's
  # web server through that device, so a browser here is a hotspot guest.
  # Override the defaults if the hotspot was set up on another subnet.
  GUEST_NET=${GUEST_NET:-192.168.88.0/24}
  GUEST_HOTSPOT=${GUEST_HOTSPOT:-192.168.88.1}
  GUEST_CLIENT=${GUEST_CLIENT:-192.168.88.254}
  LAN="user,id=n1,restrict=on,net=${GUEST_NET},host=${GUEST_CLIENT},hostfwd=tcp:127.0.0.1:${GUEST_PORT}-${GUEST_HOTSPOT}:80"
  echo "Guest mode — hotspot ${GUEST_HOTSPOT} at http://127.0.0.1:${GUEST_PORT}, guest device ${GUEST_CLIENT}"
fi

if [ "$1" = v7 ]; then
  # Not -accel hvf: RouterOS's ARM64 userland includes 32-bit ARM binaries,
  # which Apple Silicon cannot execute ("Kernel panic: No working init found").
  # An emulated Cortex-A72 runs them and still boots in ~15 s.
  exec qemu-system-aarch64 \
    -machine virt -accel tcg -cpu cortex-a72 -smp 2 -m 512 \
    -bios "$QEMU_SHARE/edk2-aarch64-code.fd" \
    -drive file="$IMG",if=none,format=raw,id=hd0 -device virtio-blk-device,drive=hd0 \
    -netdev "$NET" -device virtio-net-device,netdev=n0 \
    -netdev "$LAN" -device virtio-net-device,netdev=n1 \
    -nographic
else
  exec qemu-system-x86_64 \
    -machine q35 -accel tcg -smp 1 -m 256 \
    -drive file="$IMG",if=none,format=raw,id=hd0 -device virtio-blk-pci,drive=hd0 \
    -netdev "$NET" -device virtio-net-pci,netdev=n0 \
    -netdev "$LAN" -device virtio-net-pci,netdev=n1 \
    -nographic
fi
