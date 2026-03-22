package steam

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func steamRunning() bool {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid := entry.Name()
		if pid == "" {
			continue
		}

		if _, err := strconv.Atoi(pid); err != nil {
			continue
		}

		comm, err := os.ReadFile(filepath.Join("/proc", pid, "comm"))
		if err != nil {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(string(comm)))
		if name == "steam" {
			return true
		}
	}

	return false
}
