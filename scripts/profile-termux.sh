#!/usr/bin/env bash
# "Termux edition": a lean LineageOS for using the phone as an SSH/terminal device.
# Disables apps not needed for that (all reversible with `adb shell pm enable <pkg>`),
# applies the optimisations and installs the app bundle.
# Run after flash-gsi.sh and post-install.sh.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"

DISABLE=(
  # screensavers, easter egg, printing, live wallpapers, backup, audio effects
  com.android.dreams.basic com.android.dreams.phototable com.android.egg
  com.android.bips com.android.printspooler com.android.wallpaper.livepicker
  com.stevesoltys.seedvault org.lineageos.audiofx
  # media: music, gallery, recorder (the camera stays)
  org.lineageos.eleven org.lineageos.glimpse org.lineageos.recorder
  # calendar, and the LineageOS browser (replaced by Firefox/Fennec below)
  org.lineageos.etar org.lineageos.jelly
  # phone and SMS (remove these three lines' packages if you use a SIM for calls)
  com.android.dialer com.android.contacts com.android.messaging
)

"$HERE/optimize.sh" --disable-apps "${DISABLE[*]}"
"$HERE/install-apps.sh" --with-browser

# Firefox (Fennec F-Droid) becomes the default browser
adb shell cmd role add-role-holder --user 0 android.app.role.BROWSER org.mozilla.fennec_fdroid

echo
echo "Termux edition ready. Kept: Settings, Launcher, Camera, Clock, Calculator, Files, Keyboard + the app bundle + Firefox (Fennec F-Droid)."
