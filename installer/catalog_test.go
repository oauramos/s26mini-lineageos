package main

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// The installer and scripts/install-apps.sh must pin the same APKs.
func TestAppsMatchInstallScript(t *testing.T) {
	script, err := os.ReadFile("../scripts/install-apps.sh")
	if err != nil {
		t.Fatal(err)
	}
	for _, a := range append(appBundle, auroraStore) {
		if !strings.Contains(string(script), a.Name) || !strings.Contains(string(script), a.SHA256) {
			t.Errorf("%s (%s) not pinned identically in scripts/install-apps.sh", a.Name, a.SHA256)
		}
	}
}

func TestVBMetaMatchesPythonTool(t *testing.T) {
	out := filepath.Join(t.TempDir(), "vbmeta.img")
	if err := exec.Command("python3", "../tools/make_vbmeta_disabled.py", out).Run(); err != nil {
		t.Skip("python3 not available:", err)
	}
	want, _ := os.ReadFile(out)
	if !bytes.Equal(disabledVBMeta(), want) {
		t.Fatal("vbmeta differs from tools/make_vbmeta_disabled.py")
	}
}
