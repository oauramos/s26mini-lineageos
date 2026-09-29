package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var toolsDir string

// run executes adb or fastboot and returns combined stdout+stderr (fastboot writes getvar to stderr).
func run(timeout time.Duration, tool string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, filepath.Join(toolsDir, exe(tool)), args...)
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &out
	err := cmd.Run()
	return strings.TrimSpace(strings.ReplaceAll(out.String(), "\r", "")), err
}

// adbToFile streams the raw stdout of a phone command into a file (binary-safe, unlike run).
func adbToFile(dest, cmd string) error {
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	c := exec.CommandContext(ctx, filepath.Join(toolsDir, exe("adb")), "exec-out", cmd)
	c.Stdout = f
	if err := c.Run(); err != nil {
		return err
	}
	if st, err := f.Stat(); err != nil || st.Size() == 0 {
		return fmt.Errorf("empty")
	}
	return nil
}

func adb(args ...string) (string, error)      { return run(2*time.Minute, "adb", args...) }
func fastboot(args ...string) (string, error) { return run(2*time.Minute, "fastboot", args...) }

// shell runs a command on the phone and returns its output.
func shell(cmd string) string {
	out, _ := adb("shell", cmd)
	return out
}

func prop(name string) string { return shell("getprop " + name) }

type phoneState int

const (
	stateNone         phoneState = iota
	stateUnauthorized            // adb sees it, user hasn't accepted the RSA prompt
	stateADB                     // Android running, debugging allowed
	stateFastboot                // bootloader fastboot
	stateMany                    // more than one phone
)

func detect() phoneState {
	out, _ := run(15*time.Second, "adb", "devices")
	var adbStates []string
	for _, l := range strings.Split(out, "\n")[1:] {
		if f := strings.Fields(l); len(f) >= 2 {
			adbStates = append(adbStates, f[1])
		}
	}
	fb, _ := run(15*time.Second, "fastboot", "devices")
	fbCount := 0
	for _, l := range strings.Split(fb, "\n") {
		if strings.Contains(l, "fastboot") {
			fbCount++
		}
	}
	switch {
	case len(adbStates)+fbCount > 1:
		return stateMany
	case fbCount == 1:
		return stateFastboot
	case len(adbStates) == 1 && adbStates[0] == "device":
		return stateADB
	case len(adbStates) == 1 && adbStates[0] == "unauthorized":
		return stateUnauthorized
	}
	return stateNone
}

func bootCompleted() bool {
	return detect() == stateADB && prop("sys.boot_completed") == "1"
}

func fastbootVar(name string) string {
	out, _ := fastboot("getvar", name)
	for _, l := range strings.Split(out, "\n") {
		if strings.HasPrefix(l, name+":") {
			return strings.TrimSpace(strings.TrimPrefix(l, name+":"))
		}
	}
	return ""
}

func inFastbootd() bool {
	return detect() == stateFastboot && fastbootVar("is-userspace") == "yes"
}

// adbRoot restarts adbd as root; LineageOS userdebug builds allow it.
func adbRoot() bool {
	adb("root")
	time.Sleep(3 * time.Second)
	waitWithDots(t("Reconectando", "Reconnecting"), time.Minute, 2*time.Second, func() bool { return detect() == stateADB })
	return shell("id -u") == "0"
}
