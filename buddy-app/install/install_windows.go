//go:build windows

package install

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/sys/windows/registry"
)

const (
	runKeyPath   = `Software\Microsoft\Windows\CurrentVersion\Run`
	runKeyValue  = "SteamInputDBBuddy"
	shortcutName = "SteamInputDB Buddy"
	exeName      = "steaminputdb-buddy.exe"
)

func defaultInstallPath() string {
	return filepath.Join(os.Getenv("LOCALAPPDATA"), "SteamInputDB", exeName)
}

func setAutoStart(installPath string, enabled bool) error {
	if !enabled {
		k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
		if err != nil {
			if errors.Is(err, registry.ErrNotExist) {
				return nil
			}
			return err
		}
		defer k.Close()
		err = k.DeleteValue(runKeyValue)
		if errors.Is(err, registry.ErrNotExist) {
			return nil
		}
		return err
	}

	previousExe, err := currentAutorunExe()
	if err != nil {
		return err
	}

	k, _, err := registry.CreateKey(registry.CURRENT_USER, runKeyPath, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer k.Close()

	if err := k.SetStringValue(runKeyValue, `"`+installPath+`"`); err != nil {
		return err
	}

	if previousExe != "" {
		_ = killProcessesByExe(previousExe)
	}

	return nil
}

func autoStart() (bool, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	defer k.Close()

	_, _, err = k.GetStringValue(runKeyValue)
	if errors.Is(err, registry.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

func currentAutorunExe() (string, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.QUERY_VALUE)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", nil
		}
		return "", err
	}
	defer k.Close()

	val, _, err := k.GetStringValue(runKeyValue)
	if err != nil {
		if errors.Is(err, registry.ErrNotExist) {
			return "", nil
		}
		return "", err
	}

	trimmed := strings.TrimSpace(val)
	if trimmed == "" {
		return "", nil
	}

	if strings.HasPrefix(trimmed, `"`) {
		trimmed = strings.TrimPrefix(trimmed, `"`)
		if end := strings.Index(trimmed, `"`); end >= 0 {
			trimmed = trimmed[:end]
		}
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		return "", nil
	}

	return filepath.Clean(fields[0]), nil
}

func killProcessesByExe(target string) error {
	target = filepath.Clean(target)
	if target == "" {
		return nil
	}

	script := fmt.Sprintf(
		"$ErrorActionPreference='SilentlyContinue';$t='%s';Get-CimInstance Win32_Process | Where-Object { $_.ExecutablePath -eq $t } | Select-Object -ExpandProperty ProcessId",
		strings.ReplaceAll(target, "'", "''"),
	)
	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("process query failed: %w: %s", err, strings.TrimSpace(string(output)))
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	var pids []int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		pid, err := strconv.Atoi(line)
		if err == nil {
			pids = append(pids, pid)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	self := os.Getpid()
	for _, pid := range pids {
		if pid == self {
			continue
		}
		cmd := exec.Command("taskkill", "/PID", strconv.Itoa(pid), "/T", "/F")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("taskkill pid %d failed: %w: %s", pid, err, strings.TrimSpace(string(out)))
		}
	}

	return nil
}

func createShortcuts(installPath string, desktop, startMenu bool) error {
	if !desktop && !startMenu {
		return nil
	}

	for _, s := range []struct {
		folder string
		create bool
	}{
		{"Desktop", desktop},
		{"Programs", startMenu},
	} {
		dir, err := windowsKnownFolder(s.folder)
		if err != nil {
			return err
		}
		lnkPath := filepath.Join(dir, shortcutName+".lnk")
		if s.create {
			if err := createWindowsShortcut(lnkPath, installPath); err != nil {
				return err
			}
		} else {
			if err := os.Remove(lnkPath); err != nil && !os.IsNotExist(err) {
				return err
			}
		}
	}

	return nil
}

func windowsKnownFolder(name string) (string, error) {
	out, err := exec.Command("powershell", "-NoProfile", "-Command",
		fmt.Sprintf("[Environment]::GetFolderPath('%s')", name)).Output()
	if err != nil {
		return "", err
	}
	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", fmt.Errorf("%s folder path is empty", name)
	}
	return path, nil
}

func createWindowsShortcut(lnkPath, targetExe string) error {
	script := fmt.Sprintf(
		"$ws=New-Object -ComObject WScript.Shell;$s=$ws.CreateShortcut('%s');$s.TargetPath='%s';$s.WorkingDirectory='%s';$s.Save()",
		strings.ReplaceAll(lnkPath, "'", "''"),
		strings.ReplaceAll(targetExe, "'", "''"),
		strings.ReplaceAll(filepath.Dir(targetExe), "'", "''"),
	)
	out, err := exec.Command("powershell", "-NoProfile", "-Command", script).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}
	return nil
}

func uninstall() error {
	var errs []error

	k, err := registry.OpenKey(registry.CURRENT_USER, runKeyPath, registry.SET_VALUE)
	if err != nil {
		if !errors.Is(err, registry.ErrNotExist) {
			errs = append(errs, err)
		}
	} else {
		if err := k.DeleteValue(runKeyValue); err != nil && !errors.Is(err, registry.ErrNotExist) {
			errs = append(errs, err)
		}
		k.Close()
	}

	for _, folder := range []string{"Desktop", "Programs"} {
		dir, err := windowsKnownFolder(folder)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		if err := os.Remove(filepath.Join(dir, shortcutName+".lnk")); err != nil && !os.IsNotExist(err) {
			errs = append(errs, err)
		}
	}

	defPath := defaultInstallPath()
	_ = killProcessesByExe(defPath)

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-Command",
		fmt.Sprintf("Remove-Item -Force '%s'; Remove-Item '%s'",
			strings.ReplaceAll(defPath, "'", "''"),
			strings.ReplaceAll(filepath.Dir(defPath), "'", "''"),
		))
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000}
	if err := cmd.Start(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}
