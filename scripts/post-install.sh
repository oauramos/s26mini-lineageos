#!/usr/bin/env bash
# Device-specific fixes for Android 14+ GSIs on the fake "S26 ULTRA Mini" (MT6739).
# Run once after the first boot, with USB debugging enabled.
set -euo pipefail

adb wait-for-device
until [[ "$(adb shell getprop sys.boot_completed | tr -d '\r')" == "1" ]]; do sleep 2; done

# Bluetooth: the Android 14 stack sends HCI READ_DEFAULT_ERRONEOUS_DATA_REPORTING (0x0C5A),
# which the MT6739 controller advertises but rejects -> the BT process aborts in a loop.
# Mask it out of the supported-commands bitmap (octet 18, bit 2 -> "182").
# Needs a TrebleDroid-based GSI (TrebleDroid, AndyYan LineageOS "td").
adb shell setprop persist.sys.bt.unsupported.commands 182

adb shell cmd bluetooth_manager disable >/dev/null 2>&1 || true
sleep 2
adb shell cmd bluetooth_manager enable >/dev/null 2>&1 || true
sleep 10

state="$(adb shell dumpsys bluetooth_manager | grep -m1 -E '^\s*state:' | awk '{print $2}' | tr -d '\r')"
echo "Bluetooth state: ${state:-unknown}"
echo "Done. Remember to turn USB debugging off."
