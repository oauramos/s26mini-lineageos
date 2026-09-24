#!/usr/bin/env bash
# Replace the backdoored stock system of the fake "S26 ULTRA Mini" (MT6739, board d39g_4m_bml)
# with a Treble GSI. ERASES ALL DATA on the phone.
#
# Usage: scripts/flash-gsi.sh <system.img>
# Prereqs: adb + fastboot in PATH, phone booted with USB debugging on and
#          Developer options -> "OEM unlocking" enabled.
set -euo pipefail

EXPECTED_BOARD="d39g_4m_bml_s26ultra_mini_pt"
GSI="${1:?usage: $0 <system.img>}"
HERE="$(cd "$(dirname "$0")/.." && pwd)"
VBMETA="$(mktemp -t vbmeta_disabled).img"

[[ -f "$GSI" ]] || { echo "GSI image not found: $GSI"; exit 1; }
[[ "$GSI" != *.gz && "$GSI" != *.xz ]] || { echo "Decompress the GSI first (.img, not .gz/.xz)"; exit 1; }

confirm() { read -r -p "$1 [type YES] " a; [[ "$a" == "YES" ]] || { echo "aborted"; exit 1; }; }

wait_fastboot() {
  for _ in $(seq 1 60); do fastboot devices | grep -q . && return 0; sleep 2; done
  echo "device not seen in fastboot"; exit 1
}

wait_fastbootd() {
  for _ in $(seq 1 60); do
    fastboot getvar is-userspace 2>&1 | grep -q 'is-userspace: yes' && return 0; sleep 2
  done
  echo "fastbootd did not come up"; exit 1
}

python3 "$HERE/tools/make_vbmeta_disabled.py" "$VBMETA" >/dev/null

if adb get-state >/dev/null 2>&1; then
  board="$(adb shell getprop ro.product.vendor.device | tr -d '\r')"
  echo "Detected vendor device: $board"
  if [[ "$board" != "$EXPECTED_BOARD" ]]; then
    echo "WARNING: this script was only tested on $EXPECTED_BOARD."
    confirm "Continue on an untested board?"
  fi
  [[ "$(adb shell getprop sys.oem_unlock_allowed | tr -d '\r')" == "1" ]] || {
    echo "Enable Developer options -> OEM unlocking first."; exit 1; }
  adb reboot bootloader
fi

wait_fastboot
if ! fastboot getvar unlocked 2>&1 | grep -q 'unlocked: yes'; then
  confirm "Bootloader is locked. Unlocking WIPES the phone. Continue?"
  echo ">>> On the phone, press VOLUME UP to confirm the unlock."
  fastboot flashing unlock
fi

echo "== Disabling AVB verification"
for p in vbmeta vbmeta_system vbmeta_vendor; do
  fastboot --disable-verity --disable-verification flash "$p" "$VBMETA"
done

echo "== Entering fastbootd"
fastboot reboot fastboot
wait_fastbootd

echo "== Removing stock product partition (Android 10 apps that break newer GSIs)"
fastboot delete-logical-partition product || true

echo "== Flashing $GSI"
fastboot flash system "$GSI"

echo "== Wiping userdata"
fastboot -w || true

fastboot reboot
echo "Done. First boot can take a few minutes."
