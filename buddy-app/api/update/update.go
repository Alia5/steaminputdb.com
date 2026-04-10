package update

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

const downloadURLBase = "https://github.com/Alia5/steaminputdb.com/releases/download"

type updateRequestBody struct {
	Version string `json:"version"`
}

var _ huma.ResolverWithPath = (*updateRequestBody)(nil)

func (b *updateRequestBody) Resolve(_ huma.Context, prefix *huma.PathBuffer) []error {
	if strings.Contains(b.Version, "/") {
		return []error{&huma.ErrorDetail{
			Message:  "version must not contain slashes",
			Location: prefix.With("version"),
			Value:    b.Version,
		}}
	}
	return nil
}

type updateRequest struct {
	Body updateRequestBody
}

func RegisterRoutes(a huma.API) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/update",
			Tags:        []string{"update"},
			Summary:     "Update",
			Description: "Update SteamInputDB-Buddy to a specified version",
		},
		func(ctx context.Context, req *updateRequest) (*struct{}, error) {
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

			asset := fmt.Sprintf("steaminputdb-buddy-%s-%s", runtime.GOOS, runtime.GOARCH)
			if runtime.GOOS == "windows" {
				asset += ".exe"
			}
			downloadURL := fmt.Sprintf("%s/%s/%s", downloadURLBase, req.Body.Version, asset)

			slog.Info("Downloading update", "version", req.Body.Version, "url", downloadURL)

			// URL is constructed from validated version + constants
			resp, err := http.Get(downloadURL) // nolint:gosec
			if err != nil {
				slog.Error("failed to download update", "error", err)
				return nil, huma.Error502BadGateway("failed to download update", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				slog.Error("unexpected status code downloading update", "status", resp.StatusCode)
				return nil, huma.Error502BadGateway(
					fmt.Sprintf("failed to download update: HTTP %d", resp.StatusCode),
				)
			}

			dir := filepath.Dir(selfPath)

			tmpFile, err := os.CreateTemp(dir, "steaminputdb-buddy-update-*")
			if err != nil {
				slog.Error("failed to create temp file for update", "error", err)
				return nil, err
			}
			tmpPath := tmpFile.Name()

			_, err = io.Copy(tmpFile, resp.Body)
			tmpFile.Close()
			if err != nil {
				os.Remove(tmpPath)
				slog.Error("failed to write update to temp file", "error", err)
				return nil, err
			}

			if err = os.Chmod(tmpPath, 0755); err != nil {
				os.Remove(tmpPath)
				slog.Error("failed to chmod temp file", "error", err)
				return nil, err
			}

			oldPath := selfPath + ".old"
			if err = os.Rename(selfPath, oldPath); err != nil {
				os.Remove(tmpPath)
				slog.Error("failed to move old executable out of the way", "error", err)
				return nil, err
			}
			if err = os.Rename(tmpPath, selfPath); err != nil {
				_ = os.Rename(oldPath, selfPath)
				os.Remove(tmpPath)
				slog.Error("failed to place new executable", "error", err)
				return nil, err
			}
			os.Remove(oldPath)

			slog.Info("Update applied, restarting...", "version", req.Body.Version)

			go func() {
				time.Sleep(1 * time.Second)
				cmd := exec.Command(selfPath)
				if err := cmd.Start(); err != nil {
					slog.Error("failed to start updated executable", "error", err)
				}
				os.Exit(0)
			}()

			return nil, nil
		},
	)
}
