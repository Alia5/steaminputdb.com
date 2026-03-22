package steam_js

import (
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

type ApplyConfigArgs struct {
	AppID           uint32 `json:"appId" required:"true"`
	ControllerIndex int    `json:"controllerIndex" required:"true"`
	WorkshopItemID  string `json:"workshopItemId" required:"true"`
}

type ApplyConfigExecutor interface {
	steamcef.Executor[*ApplyConfigArgs, *struct{}]
}

//go:embed templates/apply_config.js.tmpl
var applyConfigJSTmpl string

var applyConfigJS = template.Must(template.New("applyConfig").Parse(applyConfigJSTmpl))

func NewApplyConfig(cfg *appconfig.Steam) ApplyConfigExecutor {
	return steamcef.NewExecutor[*ApplyConfigArgs, *struct{}](cfg, applyConfigJS)
}
