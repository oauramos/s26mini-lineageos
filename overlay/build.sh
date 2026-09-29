#!/usr/bin/env bash
# Build S26miniCutout.apk, the static camera-cutout overlay, into installer/ (embedded there).
# Needs Android build-tools (aapt2, apksigner), platform android-34 (android.jar) and a JDK.
#   SDK=~/.cache/s26mini-lineageos/sdk overlay/build.sh
# The signing key only has to exist (partition overlays aren't checked against the target's
# certificate); a throwaway one is generated if missing and never committed.
set -euo pipefail
HERE="$(cd "$(dirname "$0")" && pwd)"
SDK="${SDK:-$HOME/.cache/s26mini-lineageos/sdk}"
AAPT2="$SDK/android-14/aapt2"
APKSIGNER="$SDK/android-14/apksigner"
ANDROID_JAR="$SDK/android-34/android.jar"
OUT="$HERE/../installer/S26miniCutout.apk"
KEY="$HERE/.build/key.jks"

rm -rf "$HERE/.build/res" && mkdir -p "$HERE/.build/res"
"$AAPT2" compile --dir "$HERE/res" -o "$HERE/.build/res"
"$AAPT2" link -o "$HERE/.build/unsigned.apk" -I "$ANDROID_JAR" \
  --manifest "$HERE/AndroidManifest.xml" "$HERE"/.build/res/*.flat

[[ -f "$KEY" ]] || keytool -genkeypair -keystore "$KEY" -storepass s26mini -keypass s26mini \
  -alias overlay -keyalg RSA -keysize 2048 -validity 36500 -dname "CN=s26mini-lineageos" >/dev/null
"$APKSIGNER" sign --v4-signing-enabled false --ks "$KEY" --ks-pass pass:s26mini --out "$OUT" "$HERE/.build/unsigned.apk"
echo "built $OUT"
