package main

import (
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// deviceFixes is the Go version of scripts/post-install.sh.
func deviceFixes() {
	// Bluetooth: the Android 14 stack sends an HCI command the MT6739 controller rejects, so the
	// BT process aborts in a loop. TrebleApp's "mediatek" workaround masks it, and TrebleApp
	// rewrites the property on every boot from its own setting, so set both.
	// On a fresh install TrebleApp hasn't created its prefs yet, so create them owned by the app.
	const app = "/data/data/me.phh.treble.app"
	const prefs = app + "/shared_prefs/me.phh.treble.app_preferences.xml"
	shell("am force-stop me.phh.treble.app")
	switch {
	case shell("grep -q key_misc_bluetooth "+prefs+" && echo yes") == "yes":
		shell(`sed -i 's#<string name="key_misc_bluetooth">[^<]*</string>#<string name="key_misc_bluetooth">mediatek</string>#' ` + prefs)
	case shell("[ -f "+prefs+" ] && echo yes") == "yes":
		shell(`sed -i 's#</map>#    <string name="key_misc_bluetooth">mediatek</string>\n</map>#' ` + prefs)
	default:
		shell(`u=$(stat -c %u ` + app + `) && mkdir -p ` + app + `/shared_prefs && ` +
			`printf '<?xml version="1.0" encoding="utf-8" standalone="yes" ?>\n<map>\n    <string name="key_misc_bluetooth">mediatek</string>\n</map>\n' > ` + prefs + ` && ` +
			`chown -R $u:$u ` + app + `/shared_prefs && chmod 771 ` + app + `/shared_prefs && chmod 660 ` + prefs + ` && ` +
			`restorecon -R ` + app + `/shared_prefs`)
	}
	if shell("grep -q '\"key_misc_bluetooth\">mediatek<' "+prefs+" && echo yes") != "yes" {
		warn(t("Não consegui gravar o ajuste do Bluetooth no TrebleApp. Faça à mão: Phh Treble Settings > Misc > Bluetooth workarounds > Mediatek.",
			"Could not save the Bluetooth setting in TrebleApp. Set it by hand: Phh Treble Settings > Misc > Bluetooth workarounds > Mediatek."))
	}
	shell("setprop persist.sys.bt.unsupported.commands 182")
	shell("cmd bluetooth_manager disable")
	time.Sleep(2 * time.Second)
	shell("cmd bluetooth_manager enable")
	ok(t("Bluetooth corrigido", "Bluetooth fixed"))
}

//go:embed S26miniCutout.apk
var cutoutOverlay []byte // built from overlay/ by overlay/build.sh

const (
	cutoutPkg     = "io.github.oauramos.s26mini.cutout"
	cutoutOnPhone = "/system/product/overlay/S26miniCutout.apk"
)

// installCutoutOverlay puts the static camera-cutout overlay into /system. The front lens sits
// ~7.2 mm into the panel; the overlay declares a 60x72 px notch so the status bar grows around
// it and apps start below it. Runtime (fabricated) overlays don't survive a reboot on these
// GSIs, hence a real overlay APK, which needs the writable "vndklite" images.
// Returns false if /system can't be written. Takes effect after a reboot.
func installCutoutOverlay() bool {
	local := filepath.Join(cacheDir, "S26miniCutout.apk")
	if err := os.WriteFile(local, cutoutOverlay, 0o644); err != nil {
		return false
	}
	if shell("mount -o rw,remount / && echo yes") != "yes" {
		return false
	}
	// The image has ~1 MB free, less than ext4's runtime reserve for metadata (up to 16 MB),
	// so plain writes fail with ENOSPC. Lift the reserve for this one small file.
	resv := "/sys/fs/ext4/$(basename $(mount | grep ' / ' | cut -d' ' -f1))/reserved_clusters"
	old := shell("cat " + resv)
	shell("echo 0 > " + resv)
	defer shell("echo " + old + " > " + resv + "; sync; mount -o ro,remount /")
	if out, err := adb("push", local, cutoutOnPhone); err != nil {
		warn("%s", lastLine(out))
		return false
	}
	shell("chmod 644 " + cutoutOnPhone + " && restorecon " + cutoutOnPhone)
	shell("cmd overlay disable com.android.internal.display.cutout.emulation.tall")
	return true
}

// cutoutActive reports whether the phone currently reserves the 72 px camera area.
func cutoutActive() bool {
	return strings.Contains(shell("cmd overlay list android"), "[x] "+cutoutPkg) &&
		strings.Contains(shell("dumpsys window"), "type=statusBars frame=[0,0][384,72]")
}

// backupCalibration saves the preloader and the IMEI/RF calibration partitions to the computer.
// Flashing never touches them, but they are the only way back from a hard brick or a lost IMEI.
func backupCalibration() string {
	home, _ := os.UserHomeDir()
	dir := filepath.Join(home, "S26mini-backup-"+strings.ReplaceAll(prop("ro.serialno"), "/", "_"))
	os.MkdirAll(dir, 0o700)
	parts := map[string]string{
		"preloader_boot0": "/dev/block/mmcblk0boot0",
		"nvram":           "/dev/block/by-name/nvram",
		"nvdata":          "/dev/block/by-name/nvdata",
		"nvcfg":           "/dev/block/by-name/nvcfg",
		"proinfo":         "/dev/block/by-name/proinfo",
	}
	saved := 0
	for name, dev := range parts {
		dest := filepath.Join(dir, name+".img")
		if st, err := os.Stat(dest); err == nil && st.Size() > 0 {
			saved++
			continue
		}
		if err := adbToFile(dest, "dd if="+dev+" bs=1M 2>/dev/null"); err != nil {
			os.Remove(dest)
			warn(t("Não consegui copiar %s", "Could not copy %s"), name)
			continue
		}
		saved++
	}
	if saved == len(parts) {
		ok(t("Backup salvo em %s", "Backup saved to %s"), dir)
	}
	return dir
}
