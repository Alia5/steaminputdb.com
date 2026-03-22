package apply_config

import (
	"context"
	"net/http"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type requestBody struct {
	Body              steam_js.ApplyConfigArgs
	OpenConfigurator  bool `query:"openConfigurator" default:"false"`
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/steam/apply_config",
			Tags:        []string{"steam"},
			Summary:     "Apply Controller Config",
			Description: "Apply a workshop controller configuration for a given app and controller",
		},
		func(ctx context.Context, req *requestBody) (*struct{}, error) {
			_, err := steam_js.NewApplyConfig(&cfg.Steam).Execute(ctx, &req.Body)
			if err != nil {
				return nil, err
			}
			if req.OpenConfigurator {
				_, err := steam_js.NewOpenConfigurator(&cfg.Steam).Execute(ctx, &steam_js.OpenConfiguratorArgs{
					AppID: req.Body.AppID,
				})
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		},
	)
}
