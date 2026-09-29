package main

import (
	"strings"
	"time"
)

// installApps downloads, verifies and installs APKs (the Go version of scripts/install-apps.sh).
func installApps(apps []download) {
	for _, a := range apps {
		path := fetch(a)
		out, err := run(5*time.Minute, "adb", "install", "-r", path)
		if err != nil || !strings.Contains(out, "Success") {
			warn(t("Não consegui instalar %s: %s", "Could not install %s: %s"), a.Name, lastLine(out))
			continue
		}
		ok(t("%s instalado", "%s installed"), appName(a.Name))
	}
}

func appName(file string) string {
	names := map[string]string{
		"org.fdroid.fdroid":           "F-Droid",
		"com.termux":                  "Termux",
		"org.localsend.localsend_app": "LocalSend",
		"me.zhanghai.android.files":   "Material Files",
		"org.mozilla.fennec_fdroid":   "Firefox (Fennec F-Droid)",
		"org.primftpd":                "Primitive FTPd",
		"com.aurora.store":            "Aurora Store",
	}
	pkg := file[:strings.LastIndex(file, "_")]
	if n, found := names[pkg]; found {
		return n
	}
	return pkg
}

func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return lines[len(lines)-1]
}
