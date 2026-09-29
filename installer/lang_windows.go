package main

import "syscall"

// systemLang reads the Windows UI language (Portuguese primary language id is 0x16).
func systemLang() string {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetUserDefaultUILanguage")
	if r, _, _ := proc.Call(); r&0x3ff == 0x16 {
		return "pt"
	}
	return "en"
}
