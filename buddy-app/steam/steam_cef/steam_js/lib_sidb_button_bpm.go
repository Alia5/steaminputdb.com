package steam_js

import (
	"context"
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

type AddSteamInputDbButtonBPMArgs struct {
	Override bool `json:"override"`
}

type AddSteamInputDBButtonBPM interface {
	steamcef.Executor[*AddSteamInputDbButtonBPMArgs, *struct{}]
}

//go:embed templates/dist/sidb_button_BPM.js.tmpl
var addSteamInputDBButtonBPMTmpl string

var addSteamInputDBButtonBPMJS = template.Must(template.New("addSteamInputDBButtonBPM").Delims("<<%", "%>>").Parse(addSteamInputDBButtonBPMTmpl))

func NewAddSteamInputDBButtonBPM(cfg *appconfig.Steam) AddSteamInputDBButtonBPM {
	return steamcef.NewExecutor[*AddSteamInputDbButtonBPMArgs, *struct{}](cfg, addSteamInputDBButtonBPMJS)
}

func AddSteamInputDBButton_BPM(ctx context.Context, cfg *appconfig.Steam, override bool) error {
	executor := NewAddSteamInputDBButtonBPM(cfg)
	_, err := executor.ExecuteInTab(ctx, "Steam Big Picture Mode", &AddSteamInputDbButtonBPMArgs{
		Override: override,
	})
	if err != nil {
		return err
	}
	return nil
}
