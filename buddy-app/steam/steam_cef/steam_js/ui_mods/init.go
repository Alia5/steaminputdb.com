package uimods

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/Alia5/steaminputdb.com/buddy-app/db/models"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam"
	"github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js"
)

var initPollInterval = 2 * time.Second
var refreshPollInterval = 15 * time.Second
var defaultWaitTimeout = 60 * time.Second
var injectTimeout = 1 * time.Second

var cleanupTabs = []string{
	"Steam",
	"Steam Big Picture Mode",
	"SharedJSContext",
}

func Init(cfg *appconfig.Config, dal *db.DAL) {
	ctx := context.Background()
	settings, err := dal.Settings.Get(ctx)
	if err != nil {
		slog.Error("Failed to get settings", "error", err)
	}

	useSteamBrowser,
		addDesktopUIEntries,
		addBPMUIEntries,
		waitTimeout := applyDefaults(settings)

	initTicker := time.NewTicker(initPollInterval)
	var initWaitStop *time.Timer
	if waitTimeout > 0 {
		initWaitStop = time.AfterFunc(waitTimeout, func() {
			slog.Warn("Steam CEF remote debug not reachable after waiting", "waitTimeout", waitTimeout)
		})
	}
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for range ticker.C {
			if !steam.CEFRemoteDebugReachable(ctx, &cfg.Steam) {
				continue
			}
			injectCtx, cancelFunc := context.WithTimeout(context.Background(), injectTimeout)
			err := InjectUiMods(
				injectCtx, cfg,
				useSteamBrowser, addDesktopUIEntries, addBPMUIEntries,
				true,
			)
			cancelFunc()
			if err != nil {
				slog.Error("Failed to inject UI mods", "error", err)
				continue
			}

			slog.Info("Successfully injected UI mods")
			if initWaitStop != nil {
				initWaitStop.Stop()
			}
			startPolling(cfg, dal)
			return
		}
	}(initTicker)
}

func startPolling(cfg *appconfig.Config, dal *db.DAL) {
	knownTabIDs := []string{}
	refreshTicker := time.NewTicker(refreshPollInterval)
	go func(ticker *time.Ticker) {
		defer ticker.Stop()
		for range ticker.C {
			ctx := context.Background()
			tabCtx, cancel := context.WithTimeout(ctx, injectTimeout)
			tabs, err := steam.GetCEFTabs(tabCtx, &cfg.Steam)
			cancel()
			if err != nil {
				slog.Error("Failed to get CEF tabs", "error", err)
				continue
			}
			currentIDs := []string{}
			for _, tab := range tabs {
				if slices.Contains(cleanupTabs, tab.Title) {
					currentIDs = append(currentIDs, tab.ID)
				}
			}
			slices.Sort(currentIDs)
			slog.Debug("startPolling: tick", "currentIDs", currentIDs, "knownIDs", knownTabIDs)
			if slices.Equal(currentIDs, knownTabIDs) {
				continue
			}
			slog.Debug("startPolling: tabs changed, injecting")
			knownTabIDs = currentIDs
			if err = refreshUIMods(ctx, cfg, dal, false); err != nil {
				slog.Error("Failed to refresh UI mods", "error", err)
			}
		}
	}(refreshTicker)
}

func refreshUIMods(ctx context.Context, cfg *appconfig.Config, dal *db.DAL, settingsChanged bool) error {
	settings, err := dal.Settings.Get(ctx)
	if err != nil {
		slog.Error("Failed to get settings", "error", err)
		return err
	}

	useSteamBrowser,
		addDesktopUIEntries,
		addBPMUIEntries,
		_ := applyDefaults(settings)

	ctx, cancelFunc := context.WithTimeout(context.Background(), injectTimeout)
	defer cancelFunc()
	return InjectUiMods(
		ctx,
		cfg,
		useSteamBrowser,
		addDesktopUIEntries,
		addBPMUIEntries,
		settingsChanged,
	)
}

func InjectUiMods(
	ctx context.Context,
	cfg *appconfig.Config,
	useSteamBrowser bool,
	addDesktopUIEntries bool,
	addBPMUIEntries bool,
	override bool,
) error {

	var desktopInjectErr error
	if addDesktopUIEntries {
		desktopInjectErr = steam_js.AddSteamInputDBButton_Desktop(ctx, &cfg.Steam, override, useSteamBrowser)
	}

	var bpmInjectErr error
	if addBPMUIEntries {
		bpmInjectErr = steam_js.AddSteamInputDBButton_BPM(ctx, &cfg.Steam, override)
	}
	if desktopInjectErr != nil && bpmInjectErr != nil {
		return fmt.Errorf("desktop inject error: %w; BPM inject error: %w", desktopInjectErr, bpmInjectErr)
	}
	return nil
}

func applyDefaults(settings *models.Settings) (
	useSteamBrowser bool,
	addDesktopUIEntries bool,
	addBPMUIEntries bool,
	waitTimeout time.Duration,
) {
	useSteamBrowser = false
	if settings.DesktopUseSteamBrowser != nil {
		useSteamBrowser = *settings.DesktopUseSteamBrowser
	}
	addDesktopUIEntries = true
	if settings.AddDesktopUIEntries != nil {
		addDesktopUIEntries = *settings.AddDesktopUIEntries
	}
	addBPMUIEntries = true
	if settings.AddBigPictureUIEntries != nil {
		addBPMUIEntries = *settings.AddBigPictureUIEntries
	}
	waitTimeout = defaultWaitTimeout
	if settings.SteamWaitTimeout != nil {
		waitTimeout = *settings.SteamWaitTimeout
	}
	return
}
