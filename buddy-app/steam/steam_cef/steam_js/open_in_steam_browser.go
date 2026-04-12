package steam_js

import (
	"context"
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

//go:embed templates/dist/open_in_steam_browser.js.tmpl
var openInSteamBrowserJSTmpl string

var openInSteamBrowserJS = template.Must(template.New("openInSteamBrowser").Delims("<<%", "%>>").Parse(openInSteamBrowserJSTmpl))

type templateParams struct {
	URL string
}

func newOpenInSteamBrowser(cfg *appconfig.Steam) steamcef.Executor[templateParams, struct{}] {
	return steamcef.NewExecutor[templateParams, struct{}](cfg, openInSteamBrowserJS)
}

func OpenInSteamBrowser(ctx context.Context, cfg *appconfig.Steam, url string) error {
	executor := newOpenInSteamBrowser(cfg)

	tabs, err := steam.GetCEFTabs(ctx, cfg)
	if err != nil {
		return err
	}
	var tab *steam.TabInfo
	for _, t := range tabs {
		if t.Title == "SharedJSContext" {
			tab = &t
			break
		}
	}
	if tab == nil {
		return ErrStoreTabNotFound
	}
	_, err = executor.ExecuteInTab(ctx, tab.Title, templateParams{
		URL: url,
	})
	if err != nil {
		return err
	}
	return nil
}
