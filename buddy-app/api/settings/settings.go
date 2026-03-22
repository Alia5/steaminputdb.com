package settings

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"time"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/Alia5/steaminputdb.com/buddy-app/db/models"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	uimods "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js/ui_mods"
	"github.com/danielgtaylor/huma/v2"
)

type settingsBody struct {
	AddDesktopUIEntries    *bool     `json:"addDesktopUIEntries,omitempty"`
	AddBigPictureUIEntries *bool     `json:"addBigPictureUIEntries,omitempty"`
	SteamWaitTimeout       *duration `json:"steamWaitTimeout,omitempty" format:"duration"`
	DesktopUseSteamBrowser *bool     `json:"desktopUseSteamBrowser,omitempty"`
	AutoStart              *bool     `json:"autoStart,omitempty"`
}

type getResponse struct {
	Body settingsBody
}

type putRequest struct {
	Body settingsBody
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	registry := a.OpenAPI().Components.Schemas

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/settings",
			Tags:        []string{"settings"},
			Summary:     "Get Settings",
			Description: "Get the current application settings",
		},
		func(ctx context.Context, _ *struct{}) (*getResponse, error) {
			s, err := dal.Settings.Get(ctx)
			if err != nil {
				return nil, err
			}

			autoStart, err := install.AutoStart()
			if err != nil {
				return nil, err
			}

			return &getResponse{
				Body: settingsBody{
					AddDesktopUIEntries:    s.AddDesktopUIEntries,
					AddBigPictureUIEntries: s.AddBigPictureUIEntries,
					DesktopUseSteamBrowser: s.DesktopUseSteamBrowser,
					SteamWaitTimeout:       (*duration)(s.SteamWaitTimeout),
					AutoStart:              &autoStart,
				},
			}, nil
		},
	)

	inputSchema := registry.Schema(reflect.TypeFor[settingsBody](), false, "")
	slog.Debug("Input schema for settingsBody", "schema", inputSchema)
	inputSchema.Properties["steamWaitTimeout"].Type = "string"
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPut,
			Path:        "/v1/settings",
			Tags:        []string{"settings"},
			Summary:     "Update Settings",
			Description: "Update application settings",
			RequestBody: &huma.RequestBody{
				Description: "Settings to update",
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: inputSchema,
					},
				},
			},
		},
		func(ctx context.Context, req *putRequest) (*getResponse, error) {
			s := &models.Settings{
				AddDesktopUIEntries:    req.Body.AddDesktopUIEntries,
				AddBigPictureUIEntries: req.Body.AddBigPictureUIEntries,
				DesktopUseSteamBrowser: req.Body.DesktopUseSteamBrowser,
				SteamWaitTimeout:       (*time.Duration)(req.Body.SteamWaitTimeout),
			}

			autoStart, err := install.AutoStart()
			if err != nil {
				return nil, err
			}
			if req.Body.AutoStart != nil && *req.Body.AutoStart != autoStart {
				exe, err := os.Executable()
				if err != nil {
					slog.Error("Failed to get auto-start status", "error", err)
					return nil, err
				}
				selfPath, err := filepath.EvalSymlinks(exe)
				if err != nil {
					slog.Error("Failed to evaluate symlinks", "error", err)
					return nil, err
				}
				err = install.SetAutoStart(selfPath, *req.Body.AutoStart)
				if err != nil {
					slog.Error("Failed to set auto-start status", "error", err)
					return nil, err
				}
			}
			updated, err := dal.Settings.Update(ctx, s)
			if err != nil {
				return nil, err
			}
			allSettings, err := dal.Settings.Get(ctx)
			if err != nil {
				return nil, err
			}

			uimods.Cleanup(ctx, &cfg.Steam)
			uimods.InjectUiMods(
				ctx,
				cfg,
				*allSettings.DesktopUseSteamBrowser,
				*allSettings.AddDesktopUIEntries,
				*allSettings.AddBigPictureUIEntries,
				true,
			)

			return &getResponse{
				Body: settingsBody{
					AddDesktopUIEntries:    updated.AddDesktopUIEntries,
					AddBigPictureUIEntries: updated.AddBigPictureUIEntries,
					DesktopUseSteamBrowser: updated.DesktopUseSteamBrowser,
					SteamWaitTimeout:       (*duration)(updated.SteamWaitTimeout),
					AutoStart:              req.Body.AutoStart,
				},
			}, nil
		},
	)

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/settings/reset",
			Tags:        []string{"settings"},
			Summary:     "Reset Settings",
			Description: "Reset application settings to defaults",
		},
		func(ctx context.Context, _ *struct{}) (*getResponse, error) {
			s, err := dal.Settings.ResetToDefault(ctx)
			if err != nil {
				return nil, err
			}
			autoStart, err := install.AutoStart()
			if err != nil {
				return nil, err
			}
			uimods.Cleanup(ctx, &cfg.Steam)
			uimods.InjectUiMods(
				ctx,
				cfg,
				*s.DesktopUseSteamBrowser,
				*s.AddDesktopUIEntries,
				*s.AddBigPictureUIEntries,
				true,
			)
			return &getResponse{
				Body: settingsBody{
					AddDesktopUIEntries:    s.AddDesktopUIEntries,
					AddBigPictureUIEntries: s.AddBigPictureUIEntries,
					DesktopUseSteamBrowser: s.DesktopUseSteamBrowser,
					SteamWaitTimeout:       (*duration)(s.SteamWaitTimeout),
					AutoStart:              &autoStart,
				},
			}, nil
		},
	)
}

type duration time.Duration

func (d duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *duration) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	val, err := time.ParseDuration(str)
	*d = duration(val)
	return err
}

func (d duration) String() string {
	return time.Duration(d).String()
}
