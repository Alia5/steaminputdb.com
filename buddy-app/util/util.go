//go:build !windows

package util

func IsRunFromGUI() bool {
	return false
}

func HideConsoleWindow() {
	// No-op on non-Windows platforms
}
