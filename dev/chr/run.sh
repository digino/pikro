#!/usr/bin/env bash
# Boot a local RouterOS CHR (Cloud Hosted Router) VM for testing Pikro without
# touching a production router. See dev/chr/README.md.
#
#   dev/chr/run.sh v7    RouterOS 7, native ARM64 (Apple Silicon, hvf)
#   dev/chr/run.sh v6    RouterOS 6, x86 emulated (slow, for compatibility checks)
#
# The RouterOS console runs in this terminal. Quit QEMU with Ctrl-A then X.
set -euo pipefail

cd "$(dirname "$0")"
QEMU_SHARE="$(brew --prefix qemu)/share/qemu"

case "${1:-}" in
  v7)
    ZIP=images/v7/chr-7.25-arm64.img.zip
    IMG=images/v7/chr-7.25-arm64.img
    API_PORT=18728 WINBOX_PORT=18291
    ;;
  v6)
    ZIP=images/v6/chr-6.49.23.img.zip
    IMG=images/v6/chr-6.49.23.img
    API_PORT=28728 WINBOX_PORT=28291
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

# Only localhost can reach the forwarded ports.
NET="user,id=n0,hostfwd=tcp:127.0.0.1:${API_PORT}-:8728,hostfwd=tcp:127.0.0.1:${WINBOX_PORT}-:8291"

echo "RouterOS $1 — API 127.0.0.1:${API_PORT}, WinBox 127.0.0.1:${WINBOX_PORT} (Ctrl-A X quits)"

if [ "$1" = v7 ]; then
  # Not -accel hvf: RouterOS's ARM64 userland includes 32-bit ARM binaries,
  # which Apple Silicon cannot execute ("Kernel panic: No working init found").
  # An emulated Cortex-A72 runs them and still boots in ~15 s.
  exec qemu-system-aarch64 \
    -machine virt -accel tcg -cpu cortex-a72 -smp 2 -m 512 \
    -bios "$QEMU_SHARE/edk2-aarch64-code.fd" \
    -drive file="$IMG",if=none,format=raw,id=hd0 -device virtio-blk-device,drive=hd0 \
    -netdev "$NET" -device virtio-net-device,netdev=n0 \
    -nographic
else
  exec qemu-system-x86_64 \
    -machine q35 -accel tcg -smp 1 -m 256 \
    -drive file="$IMG",if=none,format=raw,id=hd0 -device virtio-blk-pci,drive=hd0 \
    -netdev "$NET" -device virtio-net-pci,netdev=n0 \
    -nographic
fi
