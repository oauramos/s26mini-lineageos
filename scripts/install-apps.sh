#!/usr/bin/env bash
# Install the app bundle from the official F-Droid repo, pinned by SHA-256.
# Any byte difference aborts the install. Updates come later through the F-Droid app itself.
#
# Bundle: F-Droid, Termux, LocalSend (arm64), Material Files (file manager with SFTP/SMB/FTP/WebDAV).
# --with-browser also installs Fennec F-Droid (Firefox built from Mozilla's source by F-Droid, arm64).
#
# Usage: scripts/install-apps.sh [--with-browser]
set -euo pipefail

CACHE="${XDG_CACHE_HOME:-$HOME/.cache}/s26mini-lineageos/apks"
mkdir -p "$CACHE"

# file  sha256 (verified against f-droid.org/repo/index-v2.json)
APPS=(
  "org.fdroid.fdroid_1023052.apk            985f5181d48bb6bafd54083a048b391271e0ab28385881cc41294fb01a222762"
  "com.termux_1002.apk                      e6265a57eb5ca363808488e3b01955958bed93bc0c8a0d281849b363b11027ec"
  "org.localsend.localsend_app_643.apk      82ec3568fba2aa5295b9aae8b76f701d7a4703d86b9f8bad749472038fbaeab3"
  "me.zhanghai.android.files_40.apk         2fe900bf43d725b655008d438f5ca46d0da0f301e3784403b5b96c3fc2df6e8a"
)
if [[ "${1:-}" == "--with-browser" ]]; then
  APPS+=("org.mozilla.fennec_fdroid_1560020.apk    27f2951376ca1085e0933066c902fdfe4d260916d381e233475c5fd6964f5745")
fi

sha256() { if command -v sha256sum >/dev/null; then sha256sum "$1" | cut -d' ' -f1; else shasum -a 256 "$1" | cut -d' ' -f1; fi; }

adb wait-for-device

for entry in "${APPS[@]}"; do
  read -r file want <<<"$entry"
  dest="$CACHE/$file"
  if [[ ! -f "$dest" || "$(sha256 "$dest")" != "$want" ]]; then
    echo "Downloading $file"
    # Older versions move to /archive once F-Droid publishes newer ones
    curl -fsSL -o "$dest" "https://f-droid.org/repo/$file" \
      || curl -fsSL -o "$dest" "https://f-droid.org/archive/$file"
  fi
  got="$(sha256 "$dest")"
  if [[ "$got" != "$want" ]]; then
    echo "SHA-256 mismatch for $file"; echo "  expected $want"; echo "  got      $got"
    rm -f "$dest"; exit 1
  fi
  printf '%-40s ' "$file"
  adb install -r "$dest" | tail -1
done

echo "Done. Open F-Droid once to let it refresh the repo and offer updates."
