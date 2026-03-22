package uimods

import (
	"context"
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

type AddSteamInputDbButtonDesktopArgs struct {
	Override        bool `json:"override"`
	UseSteamBrowser bool `json:"useSteamBrowser"`
	//
	Logo string `json:"logo"`
}

type AddSteamInputDBButtonDesktop interface {
	steamcef.Executor[*AddSteamInputDbButtonDesktopArgs, *struct{}]
}

//go:embed templates/lib_sidb_button_desktop.js.tmpl
var addSteamInputDBButtonDesktopTmpl string

var addSteamInputDBButtonDesktopJS = template.Must(template.New("addSteamInputDBButtonDesktop").Parse(addSteamInputDBButtonDesktopTmpl))

func NewAddSteamInputDBButtonDesktop(cfg *appconfig.Steam) AddSteamInputDBButtonDesktop {
	return steamcef.NewExecutor[*AddSteamInputDbButtonDesktopArgs, *struct{}](cfg, addSteamInputDBButtonDesktopJS)
}

func AddSteamInputDBButton_Desktop(ctx context.Context, cfg *appconfig.Steam, override bool, useSteamBrowser bool) error {
	executor := NewAddSteamInputDBButtonDesktop(cfg)
	_, err := executor.ExecuteInTab(ctx, "Steam", &AddSteamInputDbButtonDesktopArgs{
		Override:        override,
		UseSteamBrowser: useSteamBrowser,
		//
		Logo: logoStr,
	})
	if err != nil {
		return err
	}
	return nil
}
