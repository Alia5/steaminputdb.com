//go:build !windows

package install

import (
	"bytes"
	_ "embed"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

const (
	serviceFileName  = "steaminputdb-buddy.service"
	desktopEntryName = "steaminputdb-buddy.desktop"
)

//go:embed icon.png
var iconData []byte

var unitTmpl = template.Must(template.New("unit").Parse(`[Unit]
Description=SteamInputDB Buddy
After=graphical-session.target

[Service]
Type=simple
ExecStart={{.ExecStart}}
WorkingDirectory={{.WorkingDirectory}}
Restart=on-failure
RestartSec=5

[Install]
WantedBy=default.target
`))

var desktopEntryTmpl = template.Must(template.New("desktop").Parse(`[Desktop Entry]
Type=Application
Name=SteamInputDB Buddy
Exec={{.Exec}}
Icon={{.Icon}}
Terminal=false
Categories=Utility;Game;
Comment=SteamInputDB Buddy App
`))

func defaultInstallPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "bin", "steaminputdb-buddy")
}

func serviceFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "systemd", "user", serviceFileName), nil
}

func setAutoStart(installPath string, enabled bool) error {
	path, err := serviceFilePath()
	if err != nil {
		return err
	}

	if !enabled {
		_ = exec.Command("systemctl", "--user", "disable", "--now", serviceFileName).Run()
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
		_ = exec.Command("systemctl", "--user", "daemon-reload").Run()
		return nil
	}

	var buf bytes.Buffer
	if err := unitTmpl.Execute(&buf, map[string]string{
		"ExecStart":        strconv.Quote(installPath),
		"WorkingDirectory": filepath.Dir(installPath),
	}); err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return err
	}

	if err := exec.Command("systemctl", "--user", "daemon-reload").Run(); err != nil {
		return err
	}
	if err := exec.Command("systemctl", "--user", "enable", "--now", serviceFileName).Run(); err != nil {
		return err
	}

	return nil
}

func autoStart() (bool, error) {
	path, err := serviceFilePath()
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	out, err := exec.Command("systemctl", "--user", "is-enabled", serviceFileName).Output()
	if err != nil {
		return false, nil
	}
	return strings.Contains(strings.ToLower(string(out)), "enabled"), nil
}

func createShortcuts(installPath string, desktop, startMenu bool) error {
	if !desktop && !startMenu {
		return nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	iconDir := filepath.Join(home, ".local", "share", "icons", "hicolor", "256x256", "apps")
	if err := os.MkdirAll(iconDir, 0755); err != nil {
		return err
	}
	iconPath := filepath.Join(iconDir, "steaminputdb-buddy.png")
	if err := os.WriteFile(iconPath, iconData, 0644); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := desktopEntryTmpl.Execute(&buf, map[string]string{
		"Exec": installPath,
		"Icon": iconPath,
	}); err != nil {
		return err
	}
	entry := buf.Bytes()

	desktopPath := filepath.Join(home, "Desktop", desktopEntryName)
	applicationsPath := filepath.Join(home, ".local", "share", "applications", desktopEntryName)

	if desktop {
		if err := os.WriteFile(desktopPath, entry, 0755); err != nil {
			return err
		}
	} else {
		if err := os.Remove(desktopPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	if startMenu {
		if err := os.MkdirAll(filepath.Dir(applicationsPath), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(applicationsPath, entry, 0644); err != nil {
			return err
		}
	} else {
		if err := os.Remove(applicationsPath); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func uninstall() error {
	var errs []error

	_ = exec.Command("systemctl", "--user", "disable", "--now", serviceFileName).Run()
	path, err := serviceFilePath()
	if err != nil {
		errs = append(errs, err)
	} else {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}
	_ = exec.Command("systemctl", "--user", "daemon-reload").Run()

	home, err := os.UserHomeDir()
	if err != nil {
		errs = append(errs, err)
	} else {
		for _, p := range []string{
			filepath.Join(home, "Desktop", desktopEntryName),
			filepath.Join(home, ".local", "share", "applications", desktopEntryName),
			filepath.Join(home, ".local", "share", "icons", "hicolor", "256x256", "apps", "steaminputdb-buddy.png"),
		} {
			if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
				errs = append(errs, err)
			}
		}
		if err := os.Remove(defaultInstallPath()); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
