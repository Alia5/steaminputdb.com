package config

import "time"

type Config struct {
	ConfigPath  string  `help:"Path to configuration file (json|yaml|toml)" default:"" name:"config" env:"CONFIG"`
	SteamAPIKey string  `help:"Steam API Key" env:"STEAM_API_KEY"`
	LogLevel    string  `help:"Logging level" default:"info" enum:"debug,info,warning,error" env:"LOG_LEVEL"`
	API         API     `embed:"" prefix:"api-"`
	Metrics     Metrics `embed:"" prefix:"metrics-"`
	DB          DB      `embed:"" prefix:"db-"`
	JWTSecret   string  `help:"JWT secret for signing tokens" default:"testenv" env:"JWT_SECRET"`
}

type API struct {
	ListenAddress string `help:"API server listen address" default:":8889" env:"API_LISTEN_ADDRESS"`
	PublicAddress string `help:"Public API address for OpenAPI page" default:"localhost:8889" env:"API_PUBLIC_ADDRESS"`
	CorsOrigins   string `help:"CORS allowed origins" default:"steaminputdb.com" env:"API_CORS_ORIGINS"`
}

type Metrics struct {
	ListenAddress string `help:"Metrics server listen address" default:":8899" env:"METRICS_LISTEN_ADDRESS"`
}

type DB struct {
	DatabaseURL     string        `help:"Database connection URL" default:"" env:"DATABASE_URL"`
	MaxOpenConns    int           `help:"Database max open connections" default:"25" env:"DB_MAX_OPEN_CONNS"`
	MaxIdleConns    int           `help:"Database max idle connections" default:"10" env:"DB_MAX_IDLE_CONNS"`
	MaxConnLifetime time.Duration `help:"Database max connection lifetime" default:"5m" env:"DB_MAX_CONN_LIFETIME"`
	MaxConnIdleTime time.Duration `help:"Database max idle time" default:"5m" env:"DB_MAX_IDLE_TIME"`

	SlowThreshold          time.Duration `help:"Log queries slower than this as warnings" default:"200ms" env:"DB_SLOW_THRESHOLD"`
	SkipDefaultTransaction bool          `help:"Skip the implicit transaction around single writes" default:"true" env:"DB_SKIP_DEFAULT_TRANSACTION"`
	PrepareStatement       bool          `help:"Cache prepared statements" default:"true" env:"DB_PREPARE_STATEMENT"`
}

var Parsed Config
