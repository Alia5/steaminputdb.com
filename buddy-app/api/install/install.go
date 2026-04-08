package install

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	"github.com/danielgtaylor/huma/v2"
)

type InstallRequestBody struct {
	InstallFile               *string `json:"installFile"`
	DefaultInstallDir         *bool   `json:"defaultInstallDir"`
	AutoStart                 *bool   `json:"autoStart"`
	EnableSteamCEFRemoteDebug *bool   `json:"enableSteamCEF"`
	DesktopShortcut           *bool   `json:"desktopShortcut"`
	StartMenuShortcut         *bool   `json:"startMenuShortcut"`
}

type UninstallRequestBody struct {
}

type InstallRequest struct {
	Body InstallRequestBody
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/install",
			Tags:        []string{"install"},
			Summary:     "Install",
			Description: "Install SteamInputDB-Buddy",
		},
		func(ctx context.Context, req *InstallRequest) (*struct{}, error) {

			exe, err := os.Executable()
			if err != nil {
				slog.Error("couldn't detect self executable", "error", err)
				return nil, err
			}

			selfPath, err := filepath.EvalSymlinks(exe)
			if err != nil {
				slog.Error("couldn't evaluate symlinks for executable path", "path", exe, "error", err)
				return nil, err
			}
			installPath := selfPath
			if req.Body.DefaultInstallDir != nil && *req.Body.DefaultInstallDir {
				installPath = install.DefaultInstallPath()
			} else if req.Body.InstallFile != nil {
				installPath = *req.Body.InstallFile
			}

			if installPath != selfPath {
				slog.Info("Copying executable", "from", selfPath, "to", installPath)
				if err = install.CopyExecutable(selfPath, installPath); err != nil {
					slog.Error("failed to copy executable to install path", "source", selfPath, "destination", installPath, "error", err)
					return nil, err
				}
			}
			if req.Body.AutoStart != nil && *req.Body.AutoStart {
				slog.Info("Setting auto-start")
				if err = install.SetAutoStart(installPath, *req.Body.AutoStart); err != nil {
					slog.Error("Failed to set auto-start", "error", err)
					return nil, err
				}
			}

			createDesktopShortcut := req.Body.DesktopShortcut == nil || *req.Body.DesktopShortcut
			createStartMenuShortcut := req.Body.StartMenuShortcut == nil || *req.Body.StartMenuShortcut
			err = install.CreateShortcuts(installPath, createDesktopShortcut, createStartMenuShortcut)
			if err != nil {
				slog.Error("Failed to create shortcuts", "error", err)
				return nil, err
			}

			if req.Body.EnableSteamCEFRemoteDebug != nil && *req.Body.EnableSteamCEFRemoteDebug {
				slog.Info("Enabling Steam CEF remote debugging...")
				created, err := install.EnableCefRemoteDebug(&cfg.Steam)
				if err != nil {
					slog.Error("Failed to enable Steam CEF remote debug", "error", err)
					return nil, err
				}

				if created {
					slog.Info("Steam CEF remote debug enabled successfully, restarting Steam...")
					install.RestartSteam()
					// hack! give steam some time to restart before the buddy app tries to connect to it again
					time.Sleep(10 * time.Second)
				} else {
					slog.Info("Steam CEF remote debug already enabled")
				}
			}

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
			if selfPath != installPath {
				go func() {
					time.Sleep(1 * time.Second)
					slog.Info("Installed successfully to default location", "path", installPath)
					cmd := exec.Command(installPath)
					if err = cmd.Start(); err != nil {
						slog.Error("Failed to start application after installation", "error", err)
					}
					os.Exit(0)
				}()
			}

			return nil, nil
		},
	)
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/uninstall",
			Tags:        []string{"install"},
			Summary:     "Uninstall",
			Description: "Uninstall SteamInputDB-Buddy",
		},
		func(ctx context.Context, _ *struct{}) (*struct{}, error) {
			err := install.Uninstall()
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
}
