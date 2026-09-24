#!/usr/bin/env bash
# Device-specific fixes for Android 14+ TrebleDroid-based GSIs on the fake "S26 ULTRA Mini" (MT6739).
# Run once after the first boot, with USB debugging enabled. Needs a userdebug GSI (adb root).
set -euo pipefail

adb wait-for-device
until [[ "$(adb shell getprop sys.boot_completed | tr -d '\r')" == "1" ]]; do sleep 2; done
adb root >/dev/null
sleep 3
adb wait-for-device

# Bluetooth: the Android 14 stack sends HCI READ_DEFAULT_ERRONEOUS_DATA_REPORTING (0x0C5A),
# which the MT6739 controller advertises but rejects -> the BT process aborts in a loop.
# TrebleApp's "Bluetooth workaround = mediatek" masks it (persist.sys.bt.unsupported.commands=182).
# Setting only the property is not enough: TrebleApp rewrites it on every boot from its own setting.
PREFS=/data/data/me.phh.treble.app/shared_prefs/me.phh.treble.app_preferences.xml
adb shell am force-stop me.phh.treble.app
if adb shell "[ -f $PREFS ]"; then
  if ! adb shell "grep -q key_misc_bluetooth $PREFS"; then
    adb shell "sed -i 's#</map>#    <string name=\"key_misc_bluetooth\">mediatek</string>\n</map>#' $PREFS"
  else
    adb shell "sed -i 's#<string name=\"key_misc_bluetooth\">[^<]*</string>#<string name=\"key_misc_bluetooth\">mediatek</string>#' $PREFS"
  fi
else
  echo "TrebleApp prefs not found. Open Phh Treble Settings > Misc > Bluetooth workarounds > Mediatek by hand."
fi
adb shell setprop persist.sys.bt.unsupported.commands 182

adb shell cmd bluetooth_manager disable >/dev/null 2>&1 || true
sleep 2
adb shell cmd bluetooth_manager enable >/dev/null 2>&1 || true
sleep 10

state="$(adb shell dumpsys bluetooth_manager | grep -m1 -E '^\s*state:' | awk '{print $2}' | tr -d '\r')"
echo "Bluetooth state: ${state:-unknown}"
echo "Done. Remember to turn USB debugging off."
