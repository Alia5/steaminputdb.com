package steam_js

import (
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

type OpenConfiguratorArgs struct {
	AppID uint32 `json:"appId" required:"true"`
}

type OpenConfiguratorExecutor interface {
	steamcef.Executor[*OpenConfiguratorArgs, *struct{}]
}

//go:embed templates/open_configurator.js.tmpl
var openConfiguratorJSTmpl string

var openConfiguratorJS = template.Must(template.New("openConfigurator").Parse(openConfiguratorJSTmpl))

func NewOpenConfigurator(cfg *appconfig.Steam) OpenConfiguratorExecutor {
	return steamcef.NewExecutor[*OpenConfiguratorArgs, *struct{}](cfg, openConfiguratorJS)
}
