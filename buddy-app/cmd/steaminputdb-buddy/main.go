package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	backendapi "github.com/Alia5/steaminputdb.com/api"
	buddyapi "github.com/Alia5/steaminputdb.com/buddy-app/api"
	installCmd "github.com/Alia5/steaminputdb.com/buddy-app/cmd/install"
	"github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	"github.com/Alia5/steaminputdb.com/buddy-app/middleware"
	uimods "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef/steam_js/ui_mods"
	buddytray "github.com/Alia5/steaminputdb.com/buddy-app/tray"
	"github.com/Alia5/steaminputdb.com/buddy-app/version"
	"github.com/Alia5/steaminputdb.com/logging"
	backendMiddleware "github.com/Alia5/steaminputdb.com/middleware"
	"github.com/Alia5/steaminputdb.com/routes"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/rs/cors"
)

func main() {
	userCfg := findUserConfig(os.Args[1:])
	jsonPaths, yamlPaths, tomlPaths := configCandidatePaths(userCfg)

	var cfg config.Config
	ctx := kong.Parse(&cfg,
		kong.Name("buddy-app"),
		kong.Description(fmt.Sprintf("SteamInputDB Buddy - v%s", version.Version)),
		kong.UsageOnError(),
		kong.Configuration(kong.JSON, jsonPaths...),
		kong.Configuration(kongyaml.Loader, yamlPaths...),
		kong.Configuration(kongtoml.Loader, tomlPaths...),
	)

	logging.SetupDefault(cfg.LogLevel)

	// Yeah yeah this is a misuse of kong.
	// FUCKIGN WHATEVER!
	if ctx.Selected().Name == "uninstall" {
		slog.Info("Running uninstall...")
		if err := install.Uninstall(); err != nil {
			slog.Error("Uninstall failed", "error", err)
			os.Exit(1)
		}
		return
	}
	if ctx.Selected().Name == "install" {
		err := installCmd.Install(cfg.Install, &cfg.Steam)
		if err != nil {
			slog.Error("install failed", "error", err)
			os.Exit(1)
		}
		return
	}

	dal, err := db.Init()
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
		return
	}

	schemaPrefix := "#/components/schemas/"
	schemasPath := "/schemas"

	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	apiClickable := "http://" + cfg.API.ListenAddress
	if strings.HasPrefix(cfg.API.ListenAddress, ":") {
		apiClickable = "http://localhost" + cfg.API.ListenAddress
	}
	docAPISrvs := []*huma.Server{{
		URL:         apiClickable,
		Description: "Local API",
	}}

	apiMux := http.NewServeMux()
	hAPI := humago.New(apiMux, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:   "SteamInputDB Buddy API",
				Version: version.Version,
			},
			Components: &huma.Components{
				Schemas: registry,
			},
			Servers: docAPISrvs,
		},
		OpenAPIPath:   "/openapi",
		SchemasPath:   schemasPath,
		Formats:       huma.DefaultFormats,
		DefaultFormat: "application/json",
		CreateHooks: []func(huma.Config) huma.Config{
			func(c huma.Config) huma.Config {
				linkTransformer := huma.NewSchemaLinkTransformer(schemaPrefix, c.SchemasPath)
				c.OnAddOperation = append(c.OnAddOperation, linkTransformer.OnAddOperation)
				c.Transformers = append(c.Transformers, linkTransformer.Transform)
				return c
			},
		},
		Transformers: []huma.Transformer{
			func(c huma.Context, _ string, v any) (any, error) {
				if err, is := v.(error); is {
					if sw, ok := c.BodyWriter().(*backendapi.StatusWriter); ok {
						sw.Error = err
					}
				}
				return v, nil
			},
		},
	})

	hAPI.Adapter().Handle(&huma.Operation{
		Method: http.MethodGet,
		Path:   "/docs",
	}, func(ctx huma.Context) {
		ctx.SetHeader("Content-Type", "text/html")
		_, _ = ctx.BodyWriter().Write([]byte(`<!doctype html>
			<html>
			<head>
				<title>SteamInputDB Buddy API</title>
				<meta name="referrer" content="same-origin" />
				<meta charset="utf-8" />
				<meta
				name="viewport"
				content="width=device-width, initial-scale=1" />
			</head>
			<body>
				<script
				id="api-reference"
				data-url="/openapi.yaml"></script>
				<script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
			</body>
			</html>`,
		))
	})

	buddyapi.RegisterAPI(hAPI, dal, &cfg)

	apiSrv := http.Server{
		Addr: cfg.API.ListenAddress,
		Handler: backendMiddleware.With(
			apiMux,
			logging.Middleware,
			cors.New(cors.Options{
				AllowedOrigins:   strings.Split(cfg.API.CORSOrigins, ","),
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			}).Handler,
			middleware.MeasureRequest,
			routes.UnregisteredMiddleware(hAPI),
		),
	}
	if os.Getenv("DEV") == "1" {
		yml, err := hAPI.OpenAPI().YAML()
		if err != nil {
			slog.Error("failed to generate OpenAPI YAML", "err", err)
		}

		err = os.WriteFile("../openapi-buddy.yaml", yml, 0644)
		if err != nil {
			slog.Error("failed to write OpenAPI YAML to file", "err", err)
		} else {
			slog.Info("wrote OpenAPI YAML to ../openapi-buddy.yaml")
		}
	}

	serverURL := "http://" + apiSrv.Addr
	if strings.HasPrefix(apiSrv.Addr, ":") {
		serverURL = "http://localhost" + apiSrv.Addr
	}
	slog.Info("Starting Server", "addr", apiSrv.Addr, "url", serverURL)
	slog.Info("Docs on", "addr", apiSrv.Addr, "url", serverURL+"/docs")

	sigCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.TrayDisplay {
		go buddytray.Run(dal, &cfg.Steam, stop)
	}

	uimods.Init(&cfg, dal)

	go func() {
		if err := apiSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "addr", apiSrv.Addr, "err", err)
			stop()
		}
	}()

	<-sigCtx.Done()
	slog.Info("Shutting down...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := apiSrv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown error", "err", err)
	}

	if err := uimods.Cleanup(shutdownCtx, &cfg.Steam); err != nil {
		slog.Error("UI mods cleanup error", "err", err)
	}
}

func findUserConfig(args []string) string {
	for i, a := range args {
		if strings.HasPrefix(a, "--config=") {
			return a[len("--config="):]
		}
		if a == "--config" && i+1 < len(args) {
			return args[i+1]
		}
	}
	return os.Getenv("BUDDY_CONFIG")
}

func configCandidatePaths(userPath string) (jsonPaths, yamlPaths, tomlPaths []string) {
	add := func(slice *[]string, p string) { *slice = append(*slice, p) }

	if userPath != "" {
		switch ext := filepath.Ext(userPath); ext {
		case ".json":
			add(&jsonPaths, userPath)
		case ".yaml", ".yml":
			add(&yamlPaths, userPath)
		case ".toml":
			add(&tomlPaths, userPath)
		default:
			add(&jsonPaths, userPath)
		}
	}

	wd, _ := os.Getwd()
	for _, base := range []string{"github.com/Alia5/steaminputdb.com/buddy-app", "buddy-app", "config"} {
		add(&jsonPaths, filepath.Join(wd, base+".json"))
		add(&yamlPaths, filepath.Join(wd, base+".yaml"))
		add(&yamlPaths, filepath.Join(wd, base+".yml"))
		add(&tomlPaths, filepath.Join(wd, base+".toml"))
	}

	return
}
