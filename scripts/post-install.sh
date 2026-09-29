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
APP=/data/data/me.phh.treble.app
PREFS=$APP/shared_prefs/me.phh.treble.app_preferences.xml
adb shell am force-stop me.phh.treble.app
if adb shell "grep -q key_misc_bluetooth $PREFS" 2>/dev/null; then
  adb shell "sed -i 's#<string name=\"key_misc_bluetooth\">[^<]*</string>#<string name=\"key_misc_bluetooth\">mediatek</string>#' $PREFS"
elif adb shell "[ -f $PREFS ]"; then
  adb shell "sed -i 's#</map>#    <string name=\"key_misc_bluetooth\">mediatek</string>\n</map>#' $PREFS"
else
  # Fresh install: TrebleApp hasn't created its prefs yet, so create them owned by the app.
  adb shell "u=\$(stat -c %u $APP) && mkdir -p $APP/shared_prefs && printf '<?xml version=\"1.0\" encoding=\"utf-8\" standalone=\"yes\" ?>\n<map>\n    <string name=\"key_misc_bluetooth\">mediatek</string>\n</map>\n' > $PREFS && chown -R \$u:\$u $APP/shared_prefs && chmod 771 $APP/shared_prefs && chmod 660 $PREFS && restorecon -R $APP/shared_prefs"
fi
adb shell "grep -q '\"key_misc_bluetooth\">mediatek<' $PREFS" \
  || echo "Could not save TrebleApp's setting. Open Phh Treble Settings > Misc > Bluetooth workarounds > Mediatek by hand."
adb shell setprop persist.sys.bt.unsupported.commands 182

adb shell cmd bluetooth_manager disable >/dev/null 2>&1 || true
sleep 2
adb shell cmd bluetooth_manager enable >/dev/null 2>&1 || true
sleep 10

state="$(adb shell dumpsys bluetooth_manager | grep -m1 -E '^\s*state:' | awk '{print $2}' | tr -d '\r')"
echo "Bluetooth state: ${state:-unknown}"

# Front camera: the lens sits inside the panel, ~7.2 mm from the top, and hides content.
# A static overlay (overlay/ in this repo, 60x72 px notch) declares it as a display cutout so the
# status bar grows around it and apps start below it. It goes into /system, which only the
# "vndklite" images let you write; runtime (fabricated) overlays are lost on every reboot.
HERE="$(cd "$(dirname "$0")/.." && pwd)"
adb shell cmd overlay disable com.android.internal.display.cutout.emulation.tall >/dev/null 2>&1 || true
if adb shell "mount -o rw,remount / && echo ok" | grep -q ok; then
  # The image has ~1 MB free, less than ext4's runtime metadata reserve (up to 16 MB), so plain
  # writes fail with ENOSPC. Lift the reserve for this one small file, then put it back.
  RESV="/sys/fs/ext4/\$(basename \$(mount | grep ' / ' | cut -d' ' -f1))/reserved_clusters"
  OLD_RESV="$(adb shell "cat $RESV" | tr -d '\r')"
  adb shell "echo 0 > $RESV"
  adb push "$HERE/installer/S26miniCutout.apk" /system/product/overlay/S26miniCutout.apk >/dev/null
  adb shell "chmod 644 /system/product/overlay/S26miniCutout.apk && restorecon /system/product/overlay/S26miniCutout.apk; echo $OLD_RESV > $RESV; sync; mount -o ro,remount /"
  echo "Front camera cutout installed. Rebooting to turn it on."
  adb reboot
  adb wait-for-device
  until [[ "$(adb shell getprop sys.boot_completed | tr -d '\r')" == "1" ]]; do sleep 2; done
  sleep 10
  adb shell cmd overlay list android | grep -q '\[x\] io.github.oauramos.s26mini.cutout' \
    && echo "Front camera cutout active." || echo "WARNING: the cutout overlay is not active."
else
  echo "WARNING: /system is read-only (use a *-vndklite* image). The front camera area stays unprotected."
fi
echo "Done. Remember to turn USB debugging off."
