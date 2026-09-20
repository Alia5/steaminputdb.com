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
	for _, tabs := range steamcef.UIModTabs {
		_, err := executor.ExecuteInAnyTab(ctx, tabs, &struct{}{})
		if err != nil {
			slog.Warn("cleanup failed for tab", "tab", tabs[0], "err", err)
			errs = append(errs, err)
		}
	}
	if len(errs) == len(steamcef.UIModTabs) {
		return errors.Join(errs...)
	}
	return nil
}
