package uimods

import (
	"context"
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

var cleanupTabs = []string{
	"Steam",
	"SharedJSContext",
}

type CleanupExecutor interface {
	steamcef.Executor[*struct{}, *struct{}]
}

//go:embed templates/cleanup.js.tmpl
var cleanupTmpl string

func NewCleanup(cfg *appconfig.Steam) CleanupExecutor {
	return steamcef.NewExecutor[*struct{}, *struct{}](cfg, template.Must(template.New("cleanup").Parse(cleanupTmpl)))
}

func Cleanup(ctx context.Context, cfg *appconfig.Steam) error {
	executor := NewCleanup(cfg)
	for _, tab := range cleanupTabs {
		_, err := executor.ExecuteInTab(ctx, tab, &struct{}{})
		if err != nil {
			return err
		}
	}
	return nil
}
