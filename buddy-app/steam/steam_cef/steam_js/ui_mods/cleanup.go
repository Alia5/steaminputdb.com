package uimods

import (
	"context"
	_ "embed"
	"errors"
	"log/slog"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

var cleanupTabs = []string{
	"Steam",
	"Steam Big Picture Mode",
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
	var errs []error
	for _, tab := range cleanupTabs {
		_, err := executor.ExecuteInTab(ctx, tab, &struct{}{})
		if err != nil {
			slog.Warn("cleanup failed for tab", "tab", tab, "err", err)
			errs = append(errs, err)
		}
	}
	if len(errs) == len(cleanupTabs) {
		return errors.Join(errs...)
	}
	return nil
}
