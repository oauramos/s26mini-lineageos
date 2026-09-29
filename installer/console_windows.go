package main

import (
	"os"
	"syscall"
	"unsafe"
)

// Old Windows consoles need virtual-terminal mode switched on to render ANSI colours.
func init() {
	const enableVirtualTerminalProcessing = 0x0004
	k32 := syscall.NewLazyDLL("kernel32.dll")
	getMode := k32.NewProc("GetConsoleMode")
	setMode := k32.NewProc("SetConsoleMode")
	h := syscall.Handle(os.Stdout.Fd())
	var mode uint32
	if r, _, _ := getMode.Call(uintptr(h), uintptr(unsafe.Pointer(&mode))); r != 0 {
		setMode.Call(uintptr(h), uintptr(mode|enableVirtualTerminalProcessing))
	}
}
