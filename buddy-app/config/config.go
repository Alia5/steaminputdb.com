package config

type Config struct {
	ConfigPath  string `help:"Path to configuration file (json|yaml|toml)" name:"config" env:"STEAMINPUTDB_BUDDY_CONFIG"`
	LogLevel    string `help:"Logging level" default:"info" enum:"debug,info,warning,error" env:"LOG_LEVEL"`
	TrayDisplay bool   `help:"Show tray icon" default:"true" env:"STEAMINPUTDB_BUDDY_TRAY_DISPLAY"`
	API         API    `embed:"" prefix:"api-"`
	Steam       Steam  `embed:"" prefix:"steam-"`

	Run       *struct{} `cmd:"" default:"true"`
	Install   *Install  `cmd:"install" help:"Install the application" `
	Uninstall *struct{} `cmd:"uninstall" help:"Uninstall the application"`
}

type API struct {
	ListenAddress string `help:"API server listen address" default:"localhost:5119" env:"STEAMINPUTDB_BUDDY_LISTEN_ADDRESS"`
	CORSOrigins   string `help:"CORS allowed origins" default:"https://steaminputdb.com,https://www.steaminputdb.com" env:"STEAMINPUTDB_BUDDY_CORS_ORIGINS"`
}

type Steam struct {
	InstallDir         string `help:"Steam installation directory" default:"" env:"STEAMINPUTDB_BUDDY_STEAM_INSTALL_DIR"`
	CEFRemoteDebugPort uint16 `help:"Steam CEF remote debug port" default:"8080" env:"STEAMINPUTDB_BUDDY_CEF_REMOTE_DEBUG_PORT"`
}

type Install struct {
	InPlace                   *bool `help:"Install in place (do not copy to default install location)" default:"false" json:"inPlace"`
	AutoStart                 *bool `help:"Enable auto-start on login" default:"true" json:"autoStart"`
	EnableSteamCEFRemoteDebug *bool `help:"Enable Steam CEF remote debugging (CEF Remote debugging is required!)" default:"true" json:"enableSteamCEFRemoteDebug"`
	DesktopShortcut           *bool `help:"Create desktop shortcut" default:"true" json:"desktopShortcut"`
	StartMenuShortcut         *bool `help:"Create start menu shortcut" default:"true" json:"startMenuShortcut"`
	ShowUI                    *bool `help:"Show UI after installation" default:"true" json:"showUI"`
}
