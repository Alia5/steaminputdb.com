// Package tray provides the cross-platform system tray integration.
package tray

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"fyne.io/systray"
	"github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"
	"github.com/Alia5/steaminputdb.com/buddy-app/version"
)

const defaultBuddyUIHost = "https://steaminputdb.com/"
const devBuddyUIHost = "http://localhost:5173/"

//go:embed icon.ico
var trayIconICO []byte

//go:embed icon.png
var trayIconPNG []byte

func Run(ctx context.Context, dal *db.DAL, cfg *config.Steam, shutdown func()) {
	go systray.Run(func() {
		runtime.LockOSThread()

		if runtime.GOOS == "windows" {
			systray.SetIcon(trayIconICO)
		} else {
			systray.SetIcon(trayIconPNG)
		}
		systray.SetTooltip("SteamInputDB Buddy")

		isAutoRun, err := install.AutoStart()
		if err != nil {
			slog.Error("Failed to get auto-start status", "error", err)
		}

		infoStr := fmt.Sprintf("SteamInputDB Buddy - %s", version.Version)
		versionItem := systray.AddMenuItem(infoStr, infoStr)
		versionItem.Disable()

		systray.AddSeparator()
		showUIItem := systray.AddMenuItem("Show UI", "Open SteamInputDB Buddy UI")
		systray.AddSeparator()
		autoStartItem := systray.AddMenuItemCheckbox("Run at startup", "", isAutoRun)
		systray.AddSeparator()

		exitItem := systray.AddMenuItem("Quit", "Exit SteamInputDB Buddy")

		uiURL := defaultBuddyUIHost
		if os.Getenv("DEV") != "" {
			uiURL = devBuddyUIHost
		}
		uiURL += "buddy-app/"
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case <-showUIItem.ClickedCh:
					openBrowser(dal, cfg, uiURL)
				case <-autoStartItem.ClickedCh:
					enabled := toggleAutoRun()
					if enabled {
						autoStartItem.Check()
					} else {
						autoStartItem.Uncheck()
					}
				case <-exitItem.ClickedCh:
					systray.Quit()
					shutdown()
					return
				}
			}
		}()
	}, func() {})
}

func openBrowser(dal *db.DAL, cfg *config.Steam, url string) {

	settings, err := dal.Settings.Get(context.Background())
	if err != nil {
		slog.Error("Failed to get settings", "error", err)
	}

	if settings.DesktopUseSteamBrowser != nil && *settings.DesktopUseSteamBrowser {
		cancelCtx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
		if steam.CEFRemoteDebugReachable(cancelCtx, cfg) {
			cancel()
			cancelCtx, cancel = context.WithTimeout(context.Background(), 500*time.Millisecond)
			err := steam_js.OpenInSteamBrowser(cancelCtx, cfg, url)
			cancel()
			if err != nil {
				slog.Error("Failed to open in Steam browser", "error", err)
			} else {
				return
			}
		}
	}

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	err = cmd.Start()
	if err != nil {
		slog.Error("failed to open browser", "url", url, "err", err)
	}
}

func toggleAutoRun() bool {
	enabled, err := install.AutoStart()
	if err != nil {
		slog.Error("Failed to get auto-start status", "error", err)
		return enabled
	}
	exe, err := os.Executable()
	if err != nil {
		slog.Error("Failed to get auto-start status", "error", err)
		return enabled
	}
	selfPath, err := filepath.EvalSymlinks(exe)
	if err != nil {
		slog.Error("Failed to evaluate symlinks", "error", err)
		return enabled
	}
	err = install.SetAutoStart(selfPath, !enabled)
	if err != nil {
		slog.Error("Failed to toggle auto-start", "error", err)
	} else {
		slog.Info("Auto-start toggled successfully", "enabled", !enabled)
	}
	return !enabled
}
