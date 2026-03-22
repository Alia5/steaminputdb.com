package open_configurator

import (
	"context"
	"net/http"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type requestBody struct {
	Body steam_js.OpenConfiguratorArgs
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/steam/open_configurator",
			Tags:        []string{"steam"},
			Summary:     "Open Controller Configurator",
			Description: "Open the Steam Controller Configurator for a given app",
		},
		func(ctx context.Context, req *requestBody) (*struct{}, error) {
			_, err := steam_js.NewOpenConfigurator(&cfg.Steam).Execute(ctx, &req.Body)
			if err != nil {
				return nil, err
			}
			return nil, nil
		},
	)
}
