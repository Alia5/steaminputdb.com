package execute_js

import (
	"context"
	"log/slog"
	"net/http"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type requestBody struct {
	RawBody []byte
}
type responseBody struct {
	Body *string
}

// Endpoint for executing arbitrary JavaScript in the Steam Client.
// DO NOT REGISTER THIS IN PRODUCTION BUILDS!
func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/steam/execute_js",
			Tags:        []string{"steam"},
			Summary:     "Execute JavaScript in Steam Client",
			Description: "Execute JavaScript code in the context of the Steam Client",
		},
		func(ctx context.Context, req *requestBody) (*responseBody, error) {
			slog.Info("ExecuteJS request", "body", string(req.RawBody))
			tpl, err := template.New("").Parse(string(req.RawBody))
			if err != nil {
				return nil, err
			}
			res, err := steamcef.NewExecutor[*struct{}, *string](
				&cfg.Steam,
				tpl,
			).Execute(
				ctx,
				nil,
			)
			if err != nil {
				slog.Error("ExecuteJS error", "error", err)
				return nil, err
			}
			if res == nil {
				return nil, nil
			}

			slog.Info("ExecuteJS response", "response", *res)

			return &responseBody{Body: res}, nil
		},
	)
}
