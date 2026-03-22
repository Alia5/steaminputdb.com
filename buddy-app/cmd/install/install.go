package install

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
)

var errNoGoRun = errors.New("cannot run install from 'go run'")

func Install(c *config.Install, cfg *config.Steam) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	if strings.Contains(exe, "go-build") {
		return errNoGoRun
	}
	selfPath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		return err
	}

	installPath := selfPath
	if c.InPlace != nil && *c.InPlace {
		slog.Info("Installing in place, skipping copy to default location", "path", selfPath)
	} else {
		installPath = install.DefaultInstallPath()
		if err = install.CopyExecutable(selfPath, installPath); err != nil {
			return err
		}
	}
	if c.AutoStart != nil {
		if err = install.SetAutoStart(installPath, *c.AutoStart); err != nil {
			slog.Error("Failed to set auto-start", "error", err)
			return err
		} else {
			slog.Info("Auto-start set successfully", "enabled", *c.AutoStart)
		}
	}

	createDesktopShortcut := c.DesktopShortcut == nil || *c.DesktopShortcut
	createStartMenuShortcut := c.StartMenuShortcut == nil || *c.StartMenuShortcut
	err = install.CreateShortcuts(installPath, createDesktopShortcut, createStartMenuShortcut)
	if err != nil {
		slog.Error("Failed to create shortcuts", "error", err)
		return err
	} else if createDesktopShortcut || createStartMenuShortcut {
		slog.Info("Shortcuts created successfully", "desktopShortcut", createDesktopShortcut, "startMenuShortcut", createStartMenuShortcut)
	}

	if c.EnableSteamCEFRemoteDebug != nil && *c.EnableSteamCEFRemoteDebug {
		created, err := install.EnableCefRemoteDebug(cfg)
		if err != nil {
			slog.Error("Failed to enable Steam CEF remote debug", "error", err)
			return err
		}

		if created {
			slog.Info("Steam CEF remote debug enabled successfully")
		} else {
			slog.Info("Steam CEF remote debug already enabled")
		}
	}

	if selfPath != installPath {
		slog.Info("Installed successfully to default location", "path", installPath)
	}
	cmd := exec.Command(installPath)
	if err = cmd.Start(); err != nil {
		slog.Error("Failed to start application after installation", "error", err)
		return err
	}
	if c.ShowUI != nil && *c.ShowUI {
		slog.Info("Application started successfully, showing UI as requested")
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "windows":
			cmd = exec.Command("cmd.exe", "/C", "start", "https://steaminputdb.com/buddy-app")
		default:
			cmd = exec.Command("xdg-open", "https://steaminputdb.com/buddy-app")
		}
		if err := cmd.Start(); err != nil {
			slog.Error("failed to open URL", "url", "https://steaminputdb.com/buddy-app", "error", err)
		}
	}

	return nil
}
