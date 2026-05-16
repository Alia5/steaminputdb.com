package steam

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
)

type TabInfo struct {
	Description          string `json:"description"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
}

func GetCEFTabs(ctx context.Context, cfg *appconfig.Steam) ([]TabInfo, error) {
	if cfg == nil || cfg.CEFRemoteDebugPort == 0 {
		return nil, errors.New("CEF remote debug port not configured")
	}

	if cfg.CEFRemoteDebugPort == 8080 {
		// find if user is running millenium and detect cef debug port
		port := GetCefDebugPort(ctx)
		if port != 8080 {
			slog.Info("Likely running millenium...", "CEF Remote port", port)
			cfg.CEFRemoteDebugPort = port
		}
	}

	url := fmt.Sprintf("http://localhost:%d/json", cfg.CEFRemoteDebugPort)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from CEF remote debug endpoint: %d", resp.StatusCode)
	}

	var tabs []TabInfo
	if err := json.NewDecoder(resp.Body).Decode(&tabs); err != nil {
		return nil, err
	}

	return tabs, nil
}

func GetCefDebugPort(ctx context.Context) uint16 {
	if runtime.GOOS == "linux" {
		return getCefDebugPortLinux()
	}

	if runtime.GOOS == "windows" {
		return getCefDebugPortWindows(ctx)
	}

	slog.Error("CEF debug port detection is unsupported on this OS", "os", runtime.GOOS)
	return 8080
}

func getCefDebugPortLinux() uint16 {
	const defaultPort uint16 = 8080
	const prefix = "--remote-debugging-port="

	entries, err := os.ReadDir("/proc")
	if err != nil {
		slog.Error("Failed to inspect /proc for steamwebhelper", "error", err)
		return defaultPort
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		pid := entry.Name()
		_, parseErr := strconv.Atoi(pid)
		if parseErr != nil {
			continue
		}

		commPath := filepath.Join("/proc", pid, "comm")
		comm, readErr := os.ReadFile(commPath)
		if readErr != nil {
			continue
		}

		name := strings.ToLower(strings.TrimSpace(string(comm)))
		if name != "steamwebhelper" {
			continue
		}

		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		cmdline, readErr := os.ReadFile(cmdlinePath)
		if readErr != nil || len(cmdline) == 0 {
			continue
		}

		args := strings.FieldsSeq(strings.ReplaceAll(string(cmdline), "\x00", " "))
		for arg := range args {
			if !strings.HasPrefix(arg, prefix) {
				continue
			}

			portStr := strings.TrimPrefix(arg, prefix)
			port, parseErr := strconv.Atoi(portStr)
			if parseErr != nil || port <= 0 || port > 65535 {
				continue
			}

			return uint16(port)
		}
	}

	return defaultPort
}

func CEFRemoteDebugReachable(ctx context.Context, cfg *appconfig.Steam) bool {

	tabs, err := GetCEFTabs(ctx, cfg)
	if err != nil {
		slog.Error("Error getting CEF tabs", "error", err)
		return false
	}

	for _, tab := range tabs {
		if strings.HasPrefix(tab.URL, "https://steamloopback.host") {
			return true
		}
	}

	return false
}

func CEFRemoteDebugEnableFilePresent(cfg *appconfig.Steam) (bool, error) {
	steamDir := cfg.InstallDir
	var err error
	if steamDir == "" {
		steamDir, err = ExecuteableDir()
		if err != nil {
			slog.Error("Could not determine Steam Path", "error", err)
			return false, err
		}
	}

	_, err = os.Stat(filepath.Join(steamDir, ".cef-enable-remote-debugging"))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	slog.Error("Error checking for .cef-enable-remote-debugging file", "error", err)
	return false, err
}

func ExecuteableDir() (string, error) {
	return steamPath()
}

func ClientRunning() bool {
	return steamRunning()
}
