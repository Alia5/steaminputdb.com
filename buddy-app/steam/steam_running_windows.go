package steam

import (
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func steamRunning() bool {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(snapshot)

	var proc windows.ProcessEntry32
	proc.Size = uint32(unsafe.Sizeof(proc))

	err = windows.Process32First(snapshot, &proc)
	for err == nil {
		exe := strings.ToLower(windows.UTF16ToString(proc.ExeFile[:]))
		if exe == "steam.exe" {
			return true
		}
		err = windows.Process32Next(snapshot, &proc)
	}

	return false
}
