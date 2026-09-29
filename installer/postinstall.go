package main

import (
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

	// Front camera: the lens sits ~7.2 mm into the panel. Declare it as a 60x72 px cutout
	// (x=0 is the screen centre) so the status bar grows around it and apps start below it.
	const cutout = "M -30,0 H 30 V 62 C 30,68 26,72 20,72 H -20 C -26,72 -30,68 -30,62 Z"
	shell("cmd overlay disable com.android.internal.display.cutout.emulation.tall")
	shell("cmd overlay fabricate --target android --name FrontCameraCutout android:string/config_mainBuiltInDisplayCutout 0x03 '" + cutout + "'")
	shell("cmd overlay fabricate --target android --name FrontCameraCutoutRect android:string/config_mainBuiltInDisplayCutoutRectApproximation 0x03 '" + cutout + "'")
	shell("cmd overlay fabricate --target android --name FrontCameraCutoutFill android:bool/config_fillMainBuiltInDisplayCutout 0x12 0xffffffff")
	for _, n := range []string{"FrontCameraCutout", "FrontCameraCutoutRect", "FrontCameraCutoutFill"} {
		shell("cmd overlay enable com.android.shell:" + n)
	}
	shell("pkill -f com.android.systemui")
	ok(t("Área da câmera frontal protegida", "Front camera area reserved"))
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
