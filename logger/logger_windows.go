//go:build windows

package logger

import (
	"os"
	"syscall"
	"unsafe"
)

func enableWindowsANSI() {
	var mode uint32

	stdoutHandle := syscall.Handle(os.Stdout.Fd())
	kernel32 := syscall.NewLazyDLL("kernel32.dll")

	procGetConsoleMode := kernel32.NewProc("GetConsoleMode")
	procSetConsoleMode := kernel32.NewProc("SetConsoleMode")

	r, _, _ := procGetConsoleMode.Call(
		uintptr(stdoutHandle),
		uintptr(unsafe.Pointer(&mode)),
	)

	if r != 0 {
		mode |= 0x0004
		procSetConsoleMode.Call(
			uintptr(stdoutHandle),
			uintptr(mode),
		)
	}
}
