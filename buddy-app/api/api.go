package api

import (
	"github.com/Alia5/steaminputdb.com/buddy-app/api/install"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/ping"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/settings"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam"
	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterAPI(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	ping.RegisterRoutes(a)
	settings.RegisterRoutes(a, dal, cfg)
	steam.RegisterRoutes(a, dal, cfg)
	install.RegisterRoutes(a, dal, cfg)
}
