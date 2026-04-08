package steam

import (
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam/apply_config"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam/get_apps"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam/get_controllers"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam/open_configurator"
	"github.com/Alia5/steaminputdb.com/buddy-app/api/steam/status"
	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"

	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(a huma.API, dal *db.DAL, cfg *appconfig.Config) {
	status.RegisterRoutes(a, dal, cfg)
	get_apps.RegisterRoutes(a, dal, cfg)
	get_controllers.RegisterRoutes(a, dal, cfg)
	open_configurator.RegisterRoutes(a, dal, cfg)
	apply_config.RegisterRoutes(a, dal, cfg)

	// DO NOT REGISTER THIS IN PROD BUILDS! ONLY USE FOR DEBUGGING!!!!
	// execute_js.RegisterRoutes(a, dal, cfg)
}
