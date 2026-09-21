package appinfo

import (
	"context"
	"log/slog"
	"time"

	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/Alia5/steaminputdb.com/steam/client"
	"github.com/danielgtaylor/huma/v2"
)

const (
	dbMaxAge         = 24 * time.Hour
	reconnectTimeout = 10 * time.Second
)

func RegisterRoute(a huma.API, dal db.DAL, opts ...bool) {

	registry := a.OpenAPI().Components.Schemas

	var useMemCache bool
	if len(opts) > 0 {
		useMemCache = opts[0]
	} else {
		useMemCache = true
	}
	cache := memcache.New(30*time.Minute, 1000)

	sc := client.New()
	bgCtx := context.Background()
	scDetails := client.LoginDetails{Anonymous: true, Language: "english"}
	if err := sc.Connect(bgCtx); err != nil {
		slog.Error("steam client connect failed", "error", err)
	} else if err := sc.Login(bgCtx, scDetails); err != nil {
		slog.Error("steam client login failed", "error", err)
	}
	sc.EnableAutoReconnect(scDetails, reconnectTimeout)
	registerGetAppInfo(a, dal, registry, sc, useMemCache, cache)
	registerPatchSteamInputDBInfos(a, dal, registry, sc, useMemCache, cache)

}
