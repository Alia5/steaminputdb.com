package steam_js

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

type UiModCleanupExecutor interface {
	steamcef.Executor[*struct{}, *struct{}]
}

//go:embed templates/dist/cleanup.js.tmpl
var cleanupTmpl string

func NewUiModCleanup(cfg *appconfig.Steam) UiModCleanupExecutor {
	return steamcef.NewExecutor[*struct{}, *struct{}](cfg, template.Must(template.New("cleanup").Delims("<<%", "%>>").Parse(cleanupTmpl)))
}

func UiModCleanup(ctx context.Context, cfg *appconfig.Steam) error {
	executor := NewUiModCleanup(cfg)
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
