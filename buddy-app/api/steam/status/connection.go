package status

import (
	"context"
	"net/http"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type responseBody struct {
	Body CEFStatusResponse
}
type CEFStatusResponse struct {
	SteamRunning              bool   `json:"steamRunning"`
	SteamPath                 string `json:"steamPath"`
	CEFDebugEnableFilePresent bool   `json:"cefDebugEnableFilePresent"`
	CEFRemoteDebugReachable   bool   `json:"cefRemoteDebugReachable"`
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/status",
			Tags:        []string{"steam"},
			Summary:     "Steam Client connection status",
			Description: "Steam Client information and connection status to CEF Remote Debugging",
		},
		func(ctx context.Context, _ *struct{}) (*responseBody, error) {
			steamPath := cfg.Steam.InstallDir
			var err error
			if steamPath == "" {
				steamPath, err = steam.ExecuteableDir()
				if err != nil {
					return nil, huma.Error412PreconditionFailed("Steam Path not found", err)
				}
			}

			debugFilePresent, err := steam.CEFRemoteDebugEnableFilePresent(&cfg.Steam)
			if err != nil {
				return nil, huma.Error412PreconditionFailed("", err)
			}

			steamRunning := steam.ClientRunning()
			debugReachable := false
			if steamRunning || debugFilePresent {
				debugReachable = steam.CEFRemoteDebugReachable(ctx, &cfg.Steam)
			}
			if !steamRunning && debugReachable {
				steamRunning = true
			}
			if debugReachable {
				debugFilePresent = true
			}

			return &responseBody{
				Body: CEFStatusResponse{
					SteamRunning:              steamRunning,
					SteamPath:                 steamPath,
					CEFDebugEnableFilePresent: debugFilePresent,
					CEFRemoteDebugReachable:   debugReachable,
				},
			}, nil
		},
	)
}
