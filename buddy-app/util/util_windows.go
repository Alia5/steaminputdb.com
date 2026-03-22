//go:build windows

package util

import (
	"log/slog"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	user32               = windows.NewLazySystemDLL("user32.dll")
	procGetConsoleWindow = kernel32.NewProc("GetConsoleWindow")
	procShowWindow       = user32.NewProc("ShowWindow")
	procFreeConsole      = kernel32.NewProc("FreeConsole")
)

func IsRunFromGUI() bool {
	hwnd, _, _ := procGetConsoleWindow.Call()
	hasConsole := hwnd != 0

	parentName := getParentProcessName()
	isCliParent := isCliProcess(parentName)

	slog.Debug("Parent Process Info", "parentName", parentName, "hasConsole", hasConsole, "isCliParent", isCliParent)

	if !hasConsole {
		return true
	}

	if isCliParent {
		return false
	}

	return strings.EqualFold(parentName, "explorer.exe")
}

func HideConsoleWindow() {
	hwnd, _, _ := procGetConsoleWindow.Call()
	if hwnd == 0 {
		slog.Debug("HideConsoleWindow: no console window found")
		return
	}

	_, _, _ = procShowWindow.Call(hwnd, windows.SW_HIDE)
	_, _, _ = procFreeConsole.Call()
}

func getParentProcessName() string {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return ""
	}
	defer windows.CloseHandle(snapshot)

	var pe windows.ProcessEntry32
	pe.Size = uint32(unsafe.Sizeof(pe))

	currentPID := uint32(os.Getpid())
	var parentPID uint32

	if err := windows.Process32First(snapshot, &pe); err != nil {
		return ""
	}

	for {
		if pe.ProcessID == currentPID {
			parentPID = pe.ParentProcessID
			break
		}
		if err := windows.Process32Next(snapshot, &pe); err != nil {
			return ""
		}
	}

	if parentPID == 0 {
		return ""
	}

	if err := windows.Process32First(snapshot, &pe); err != nil {
		return ""
	}

	for {
		if pe.ProcessID == parentPID {
			return windows.UTF16ToString(pe.ExeFile[:])
		}
		if err := windows.Process32Next(snapshot, &pe); err != nil {
			break
		}
	}

	return ""
}

func isCliProcess(name string) bool {
	cliProcesses := []string{
		"cmd.exe",
		"powershell.exe",
		"pwsh.exe",
		"wt.exe",
		"conhost.exe",
		"windowsterminal.exe",
	}

	nameLower := strings.ToLower(name)
	for _, cli := range cliProcesses {
		if nameLower == cli {
			return true
		}
	}
	return false
}
