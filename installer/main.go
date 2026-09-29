// Command s26mini-installer turns the fake "S26 ULTRA Mini" (MT6739) into a clean LineageOS 21
// phone in a few guided steps. It does what scripts/flash-gsi.sh, post-install.sh and
// install-apps.sh do, and downloads adb/fastboot itself, so the user needs nothing installed.
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

var version = "dev"

func main() {
	langFlag := flag.String("lang", "", "pt or en (default: from the system language)")
	imagePath := flag.String("image", "", "advanced: flash this .img instead of downloading one")
	skipFlash := flag.Bool("skip-flash", false, "advanced: only apply fixes and apps to a phone already running LineageOS")
	keepADB := flag.Bool("keep-adb", false, "advanced: leave USB debugging on at the end")
	flag.Parse()

	lang = pickLang(*langFlag)
	if *langFlag == "" {
		def := 0
		if lang == "en" {
			def = 1
		}
		if choose("Idioma / Language", []string{"Português", "English"}, def) == 1 {
			lang = "en"
		} else {
			lang = "pt"
		}
	}
	stepTotal = 7

	fmt.Printf("\n%s%s  S26 Mini · LineageOS  %s%s\n", bold, cyan, version, reset)
	say(t("Este programa troca o sistema de fábrica (com backdoor) por um LineageOS limpo e atualizado.",
		"This program replaces the stock system (which has a backdoor) with a clean, up-to-date LineageOS."))
	big(t("⚠  ISTO APAGA TUDO QUE ESTÁ NO CELULAR.", "⚠  THIS ERASES EVERYTHING ON THE PHONE."))
	say(t("Fotos, contatos e apps somem. Faça backup do que for importante antes.",
		"Photos, contacts and apps will be gone. Back up anything important first."))
	say(t("Use um cabo %sUSB-A → USB-C%s. Cabo USB-C ↔ USB-C %snão funciona%s neste celular.",
		"Use a %sUSB-A → USB-C%s cable. USB-C ↔ USB-C cables %sdon't work%s with this phone."), bold, reset, bold, reset)

	initCache()

	// 1. Choices
	step(t("Escolhas", "Choices"))
	google := false // with -skip-flash, read from the phone once it's connected
	if !*skipFlash {
		google = choose(t("Qual versão você quer?", "Which version do you want?"), []string{
			t("Sem Google  — mais privado e leve (recomendado)", "Without Google — more private and lighter (recommended)"),
			t("Com Google Play — Play Store e serviços do Google", "With Google Play — Play Store and Google services"),
		}, 0) == 1
	}
	bundle := yesNo(t("Instalar o pacote de apps? (F-Droid, Termux, LocalSend, Arquivos, Firefox, servidor FTP)",
		"Install the app bundle? (F-Droid, Termux, LocalSend, Files, Firefox, FTP server)"), !google)
	aurora := false
	if !google {
		aurora = yesNo(t("Instalar a Aurora Store? (baixa apps da Play Store sem conta Google)",
			"Install Aurora Store? (gets Play Store apps without a Google account)"), false)
	}

	// 2. Downloads (long, unattended: done before the phone steps)
	step(t("Baixando o necessário (pode levar alguns minutos)", "Downloading what's needed (may take a few minutes)"))
	toolsDir = setupPlatformTools()
	ok("adb / fastboot")
	img := *imagePath
	if img == "" && !*skipFlash {
		d := imageNoGoogle
		if google {
			d = imageGoogle
		}
		img = gunzipImage(fetch(d), d)
	}
	var apps []download
	if bundle {
		apps = append(apps, appBundle...)
	}
	if aurora {
		apps = append(apps, auroraStore)
	}
	for _, a := range apps {
		fetch(a)
	}

	// 3. Connect the phone
	step(t("Conectando o celular", "Connecting the phone"))
	connectPhone(*skipFlash)
	if *skipFlash {
		google = strings.Contains(shell("pm list packages com.android.vending"), "com.android.vending")
	}

	if !*skipFlash {
		// 4. Flash
		step(t("Instalando o LineageOS", "Installing LineageOS"))
		big(t("Última chance: a partir daqui o celular é APAGADO.", "Last chance: from here on the phone is ERASED."))
		if !confirmWord(t("APAGAR", "ERASE")) {
			die(t("Cancelado. Nada foi alterado no celular.", "Cancelled. Nothing was changed on the phone."))
		}
		flash(img)

		// 5. First boot
		step(t("Primeiro boot", "First boot"))
		waitFirstBoot()
	} else {
		stepNo += 2
	}

	// 6. Fixes and backup
	step(t("Ajustes do aparelho e backup", "Device fixes and backup"))
	if !adbRoot() {
		die(t("Sem acesso root pelo adb. Em Opções do desenvolvedor, ative \"Depuração com root\".",
			"No root over adb. In Developer options, turn on \"Rooted debugging\"."))
	}
	deviceFixes()
	backupDir := backupCalibration()

	// 7. Apps and finish
	step(t("Apps e finalização", "Apps and finishing"))
	waitWithDots(t("Esperando o sistema reiniciar a interface", "Waiting for the system UI to restart"), time.Minute, 2*time.Second, bootCompleted)
	if len(apps) > 0 {
		installApps(apps)
		if bundle {
			shell("cmd role add-role-holder --user 0 android.app.role.BROWSER " + firefoxPkg)
		}
	} else {
		say(t("Nenhum app extra escolhido.", "No extra apps chosen."))
	}
	gsfID := ""
	if google {
		gsfID = shell(`sqlite3 /data/data/com.google.android.gsf/databases/gservices.db "select value from main where name='android_id';" 2>/dev/null`)
	}
	// This build accepts adb without asking (ro.adb.secure=0): leave it off when we're done.
	if !*keepADB {
		shell("settings put global adb_enabled 0")
		ok(t("Depuração USB desligada (segurança)", "USB debugging turned off (security)"))
	}

	finish(google, gsfID, backupDir)
}

func pickLang(flagValue string) string {
	if flagValue == "en" || flagValue == "pt" {
		return flagValue
	}
	for _, v := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if s := os.Getenv(v); s != "" {
			if strings.HasPrefix(strings.ToLower(s), "pt") {
				return "pt"
			}
			return "en"
		}
	}
	return systemLang()
}

// connectPhone guides the user until the phone is reachable and checks it is the right model.
func connectPhone(skipFlash bool) {
	shown := phoneState(-1)
	for {
		st := detect()
		if st != shown {
			shown = st
			switch st {
			case stateNone:
				say(t("Ligue o celular e conecte no computador. No celular:", "Turn the phone on and plug it in. On the phone:"))
				say(t("  1. Configurações → Sobre o telefone → toque 7 vezes em %sNúmero da versão%s",
					"  1. Settings → About phone → tap %sBuild number%s 7 times"), bold, reset)
				say(t("  2. Configurações → Sistema → Opções do desenvolvedor:", "  2. Settings → System → Developer options:"))
				if !skipFlash {
					say(t("       ative %sDesbloqueio de OEM%s", "       turn on %sOEM unlocking%s"), bold, reset)
				}
				say(t("       ative %sDepuração USB%s", "       turn on %sUSB debugging%s"), bold, reset)
				say(t("  Nada aconteceu? Troque o cabo (USB-A → USB-C) ou a porta USB.", "  Nothing happens? Try another cable (USB-A → USB-C) or USB port."))
			case stateUnauthorized:
				big(t(">>> NO CELULAR: toque em \"Permitir\" na pergunta sobre depuração USB <<<",
					">>> ON THE PHONE: tap \"Allow\" on the USB debugging prompt <<<"))
				say(t("  (marque \"Sempre permitir deste computador\")", "  (tick \"Always allow from this computer\")"))
			case stateMany:
				warn(t("Mais de um celular conectado. Deixe só o S26 Mini.", "More than one phone connected. Leave only the S26 Mini."))
			}
		}
		switch st {
		case stateADB:
			board := prop("ro.product.vendor.device")
			if board != expectedBoard {
				warn(t("Modelo detectado: %q. Este instalador só foi testado no %s.",
					"Detected model: %q. This installer was only tested on %s."), board, expectedBoard)
				if !confirmWord(t("CONTINUAR", "CONTINUE")) {
					die(t("Cancelado.", "Cancelled."))
				}
			} else {
				ok(t("Celular reconhecido (%s)", "Phone recognised (%s)"), board)
			}
			if !skipFlash && prop("ro.boot.verifiedbootstate") != "orange" && prop("sys.oem_unlock_allowed") != "1" {
				warn(t("Falta ativar \"Desbloqueio de OEM\" nas Opções do desenvolvedor.",
					"\"OEM unlocking\" is still off in Developer options."))
				pause()
				continue
			}
			return
		case stateFastboot:
			if skipFlash {
				die(t("O celular está no modo fastboot. Ligue-o normalmente.", "The phone is in fastboot mode. Boot it normally."))
			}
			warn(t("O celular já está no modo fastboot; não dá para conferir o modelo.",
				"The phone is already in fastboot mode; can't check the model."))
			return
		}
		time.Sleep(2 * time.Second)
	}
}

// waitFirstBoot waits for Android to come up. These userdebug builds enable adb on first boot;
// if that ever changes, the user is asked to switch USB debugging on again.
func waitFirstBoot() {
	if waitWithDots(t("Esperando o LineageOS iniciar (até 10 min)", "Waiting for LineageOS to start (up to 10 min)"),
		10*time.Minute, 5*time.Second, bootCompleted) {
		return
	}
	say(t("O celular não respondeu sozinho. Se ele já mostra a tela de boas-vindas:",
		"The phone didn't answer on its own. If it already shows the welcome screen:"))
	say(t("  pule a configuração, depois ative de novo a %sDepuração USB%s (mesmos passos de antes).",
		"  skip the setup, then turn %sUSB debugging%s on again (same steps as before)."), bold, reset)
	connectPhone(true)
	if !waitWithDots(t("Esperando o sistema", "Waiting for the system"), 5*time.Minute, 5*time.Second, bootCompleted) {
		die(t("O celular não terminou de iniciar.", "The phone did not finish booting."))
	}
}

func finish(google bool, gsfID, backupDir string) {
	fmt.Printf("\n%s%s  ✔ %s%s\n", bold, green, t("Pronto! Seu celular está com o LineageOS.", "Done! Your phone is running LineageOS."), reset)
	say(t("Pode desconectar o cabo. No celular, siga a tela de boas-vindas.", "You can unplug the cable. On the phone, follow the welcome screen."))
	say(t("O aviso \"orange state\" em todo boot é normal com o bootloader desbloqueado.",
		"The \"orange state\" warning on every boot is normal with an unlocked bootloader."))
	fmt.Println()
	say(t("%sGuarde a pasta de backup%s (IMEI e calibração do rádio, não compartilhe):",
		"%sKeep the backup folder%s (IMEI and radio calibration, don't share it):"), bold, reset)
	say("  %s", backupDir)
	if google {
		fmt.Println()
		say(t("%sPara a Play Store funcionar%s, registre o aparelho no Google:", "%sTo make the Play Store work%s, register the phone with Google:"), bold, reset)
		say(t("  1. Entre com uma conta Google %ssecundária%s em https://www.google.com/android/uncertified",
			"  1. Sign in with a %ssecondary%s Google account at https://www.google.com/android/uncertified"), bold, reset)
		if gsfID != "" {
			say(t("  2. Cole este número: %s%s%s", "  2. Paste this number: %s%s%s"), bold, gsfID, reset)
		} else {
			say(t("  2. O número (GSF ID) aparece no app \"Device ID\" ou rodando o instalador de novo com -skip-flash.",
				"  2. The number (GSF ID) shows in a \"Device ID\" app, or run the installer again with -skip-flash."))
		}
		say(t("  3. Espere uns minutos e reinicie o celular.", "  3. Wait a few minutes and restart the phone."))
	}
	pause()
}
