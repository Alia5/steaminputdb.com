//go:build !windows

package steam

import (
	"errors"
	"os"
	"path/filepath"
)

func steamPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", errors.New("failed to determine user home directory")
	}

	return filepath.Join(home, ".steam", "steam"), nil
}
