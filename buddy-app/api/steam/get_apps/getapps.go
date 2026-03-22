package get_apps

import (
	"context"
	"net/http"
	"reflect"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type responseBody interface {
	getAppsResponse()
}

type raw []steam_js.AppInfo

func (r *raw) getAppsResponse() {}

type mapped []AppResponse

func (r *mapped) getAppsResponse() {}

type GetAppsResponse struct {
	Body responseBody
}

type requestBody struct {
	NonSteamOnly  bool `query:"nonSteamOnly" default:"false"`
	InstalledOnly bool `query:"installedOnly" default:"false"`
	Raw           bool `query:"raw" default:"false"`
}

type AppResponse struct {
	AppID      uint32 `json:"appid" example:"250900"`
	Name       string `json:"name" example:"The Binding of Isaac: Rebirth"`
	Installed  bool   `json:"installed"`
	IsNonSteam bool   `json:"isNonSteam"`
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	registry := a.OpenAPI().Components.Schemas

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/apps",
			Tags:        []string{"steam"},
			Summary:     "Get Steam Apps",
			Description: "Retrieve information about Steam apps",
			Responses: map[string]*huma.Response{
				"200": {
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: &huma.Schema{
								OneOf: []*huma.Schema{
									registry.Schema(reflect.TypeFor[[]AppResponse](), true, ""),
									registry.Schema(reflect.TypeFor[[]steam_js.AppInfo](), true, ""),
								},
							},
						},
					},
				},
			},
		},
		func(ctx context.Context, req *requestBody) (*GetAppsResponse, error) {
			res, err := steam_js.NewGetApps(&cfg.Steam).Execute(ctx, &steam_js.GetAppsArgs{
				NonSteamOnly:  req.NonSteamOnly,
				InstalledOnly: req.InstalledOnly,
			})
			if err != nil {
				return nil, err
			}
			if req.Raw {
				return &GetAppsResponse{Body: (*raw)(&res)}, nil
			}

			var mappedRes mapped
			for _, a := range res {
				installed := false
				for _, pcd := range a.PerClientData {
					if pcd.Installed != nil && *pcd.Installed {
						installed = true
						break
					}
				}
				mappedRes = append(mappedRes, AppResponse{
					AppID:      a.AppID,
					Name:       a.DisplayName,
					Installed:  installed,
					IsNonSteam: a.AppType == 0x40000000,
				})
			}
			return &GetAppsResponse{Body: &mappedRes}, nil
		},
	)
}
