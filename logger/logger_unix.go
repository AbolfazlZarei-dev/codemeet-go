//go:build !windows

package logger

func enableWindowsANSI() {
	// در لینوکس و مک نیازی به فعال‌سازی دستی ANSI Colors نیست
}
