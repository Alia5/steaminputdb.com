package get_controllers

import (
	"context"
	"net/http"
	"reflect"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

type responseBody interface {
	getControllerResponse()
}
type raw []steam_js.ControllerInfo

func (r *raw) getControllerResponse() {}

type mapped []ControllerResponse

func (r *mapped) getControllerResponse() {}

type GetControllersResponse struct {
	Body responseBody
}

type requestBody struct {
	Raw bool `query:"raw,omitempty" default:"false"`
}

type ControllerResponse struct {
	Name     string                          `json:"name" example:"Owl Controller"`
	TypeID   steamtypes.ControllerTypeNumber `json:"typeId" example:"2"`
	Type     string                          `json:"type" example:"controller_steamcontroller_gordon"`
	TypeNice string                          `json:"typeNice" example:"Steam Controller (2015)"`
	Idx      int                             `json:"index"`
}

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	registry := a.OpenAPI().Components.Schemas

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/controllers",
			Tags:        []string{"steam"},
			Summary:     "Get Steam Controllers",
			Description: "Retrieve information about connected Steam Controllers",
			Responses: map[string]*huma.Response{
				"200": {
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: &huma.Schema{
								OneOf: []*huma.Schema{
									registry.Schema(reflect.TypeFor[[]ControllerResponse](), true, ""),
									registry.Schema(reflect.TypeFor[[]steam_js.ControllerInfo](), true, ""),
								},
							},
						},
					},
				},
			},
		},
		func(ctx context.Context, req *requestBody) (*GetControllersResponse, error) {
			res, err := steam_js.NewGetControllers(&cfg.Steam).Execute(ctx, nil)
			if err != nil {
				return nil, err
			}
			if req.Raw {
				return &GetControllersResponse{Body: (*raw)(&res)}, nil
			}

			var mappedRes mapped
			for _, c := range res {
				mappedRes = append(mappedRes, ControllerResponse{
					Name:     c.Name,
					TypeID:   c.ControllerType,
					Type:     string(steamtypes.EControllerTypeFromInt(c.ControllerType)),
					TypeNice: c.ControllerType.NiceName(),
					Idx:      c.ControllerIndex,
				})
			}
			return &GetControllersResponse{Body: &mappedRes}, nil
		},
	)
}
