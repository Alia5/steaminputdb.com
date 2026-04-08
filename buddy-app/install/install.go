package install

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"
)

func DefaultInstallPath() string {
	return defaultInstallPath()
}

func SetAutoStart(installPath string, enabled bool) error {
	return setAutoStart(installPath, enabled)
}

func AutoStart() (bool, error) {
	return autoStart()
}

func EnableCefRemoteDebug(cfg *config.Steam) (fileCreated bool, err error) {
	filePresent, err := steam.CEFRemoteDebugEnableFilePresent(cfg)
	if err != nil {
		return false, err
	}
	if filePresent {
		return false, nil
	}
	steamDir := cfg.InstallDir
	if steamDir == "" {
		steamDir, err = steam.ExecuteableDir()
		if err != nil {
			return false, err
		}
	}

	filePath := filepath.Join(steamDir, ".cef-enable-remote-debugging")
	if err = os.WriteFile(filePath, nil, 0644); err != nil {
		if !errors.Is(err, os.ErrPermission) {
			return false, err
		}
		slog.Info("direct file creation failed, attempting elevated creation", "error", err)
		if err = createCefFileElevated(filePath); err != nil {
			return false, err
		}
	}
	return true, nil
}

func createCefFileElevated(filePath string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("elevated file creation not supported on %s", runtime.GOOS)
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command",
		fmt.Sprintf("Start-Process '%s' -ArgumentList '--create-cef-file','\"%s\"' -Verb RunAs",
			strings.ReplaceAll(exe, "'", "''"),
			strings.ReplaceAll(filePath, "'", "''"),
		))
	return cmd.Run()
}

func CreateCefFile(path string) error {
	return os.WriteFile(path, nil, 0644)
}

func RestartSteam() {
	openURL("steam://exit")
	time.Sleep(5 * time.Second)
	openURL("steam://open/main")
	time.Sleep(5 * time.Second)
}

func openURL(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd.exe", "/C", "start", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		slog.Error("failed to open steam URL", "url", url, "error", err)
	}
}

func CreateShortcuts(installPath string, desktop, startMenu bool) error {
	return createShortcuts(installPath, desktop, startMenu)
}

func Uninstall() error {
	return uninstall()
}

func CopyExecutable(src string, dst string) error {
	slog.Info("Copying executable to default location", "src", src, "dst", dst)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	if err == nil {
		slog.Info("Successfully copied executable to default location", "src", src, "dst", dst)
	}
	return err
}
