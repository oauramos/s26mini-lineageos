package main

import (
	"os"
	"path/filepath"
	"time"
)

// flash is the Go version of scripts/flash-gsi.sh. The phone must be in adb or fastboot mode.
func flash(img string) {
	vbmeta := filepath.Join(cacheDir, "vbmeta_disabled.img")
	if err := os.WriteFile(vbmeta, disabledVBMeta(), 0o644); err != nil {
		die("%v", err)
	}

	if detect() == stateADB {
		say(t("Reiniciando o celular no modo fastboot...", "Rebooting the phone into fastboot..."))
		adb("reboot", "bootloader")
	}
	if !waitWithDots(t("Esperando o modo fastboot", "Waiting for fastboot"), 2*time.Minute, 2*time.Second,
		func() bool { return detect() == stateFastboot }) {
		die(t("O celular não entrou no modo fastboot. No Windows isso costuma ser falta de driver USB.",
			"The phone didn't show up in fastboot. On Windows this is usually a missing USB driver."))
	}

	if fastbootVar("unlocked") != "yes" {
		big(t(">>> NO CELULAR: aperte VOLUME + (volume para cima) para confirmar o desbloqueio <<<",
			">>> ON THE PHONE: press VOLUME UP to confirm the unlock <<<"))
		out, err := run(3*time.Minute, "fastboot", "flashing", "unlock")
		if err != nil || fastbootVar("unlocked") != "yes" {
			die(t("O desbloqueio não foi confirmado.\n%s", "The unlock wasn't confirmed.\n%s"), out)
		}
		ok(t("Bootloader desbloqueado", "Bootloader unlocked"))
	}

	say(t("Desativando a verificação de boot (AVB)...", "Disabling boot verification (AVB)..."))
	for _, p := range []string{"vbmeta", "vbmeta_system", "vbmeta_vendor"} {
		if out, err := fastboot("--disable-verity", "--disable-verification", "flash", p, vbmeta); err != nil {
			die("fastboot flash %s: %s", p, out)
		}
	}
	ok("AVB")

	say(t("Entrando no fastbootd...", "Entering fastbootd..."))
	fastboot("reboot", "fastboot")
	if !waitWithDots(t("Esperando o fastbootd", "Waiting for fastbootd"), 3*time.Minute, 3*time.Second, inFastbootd) {
		die(t("O fastbootd não apareceu.", "fastbootd did not come up."))
	}

	// The stock product partition holds Android 10 apps that break newer GSIs; its space goes to system.
	fastboot("delete-logical-partition", "product")

	big(t("Gravando o sistema. NÃO desconecte o cabo (uns 5 minutos)...",
		"Writing the system. DO NOT unplug the cable (about 5 minutes)..."))
	if out, err := run(30*time.Minute, "fastboot", "flash", "system", img); err != nil {
		die(t("Falha ao gravar o sistema:\n%s\nO celular continua em fastboot: rode o instalador de novo.",
			"Writing the system failed:\n%s\nThe phone is still in fastboot: run the installer again."), out)
	}
	ok(t("Sistema gravado", "System written"))

	say(t("Apagando os dados antigos...", "Wiping old data..."))
	run(5*time.Minute, "fastboot", "-w")
	fastboot("reboot")
	ok(t("Reiniciando. O primeiro boot demora alguns minutos.", "Rebooting. The first boot takes a few minutes."))
}
