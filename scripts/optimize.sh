#!/usr/bin/env bash
# Optional tweaks for this 3 GB / 384x854 phone. All reversible (see docs/device-fixes.md).
# Usage: scripts/optimize.sh [--disable-apps "pkg1 pkg2 ..."]
set -euo pipefail

adb wait-for-device

# Snappier UI without removing anything
for k in window_animation_scale transition_animation_scale animator_duration_scale; do
  adb shell settings put global "$k" 0.5
done

# 190 dpi is oversized on a 384x854 panel; 160 fits more content (nice for terminals)
adb shell wm density 160

# The setup wizard keeps a process around after setup is done
if [[ "$(adb shell settings get secure user_setup_complete | tr -d '\r')" == "1" ]]; then
  adb shell pm disable-user --user 0 org.lineageos.setupwizard >/dev/null 2>&1 || true
fi

# Optional: disable apps you don't use (one pm call per package)
if [[ "${1:-}" == "--disable-apps" && -n "${2:-}" ]]; then
  for p in $2; do
    adb shell pm disable-user --user 0 "$p" | tail -1
  done
fi

# Ahead-of-time compile installed apps (takes a couple of minutes)
adb shell cmd package bg-dexopt-job

adb shell grep MemAvailable /proc/meminfo
