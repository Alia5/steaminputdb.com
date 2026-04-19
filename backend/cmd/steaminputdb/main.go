package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "embed"

	"github.com/Alia5/steaminputdb.com/api"
	"github.com/Alia5/steaminputdb.com/api/steam"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/Alia5/steaminputdb.com/logging"
	"github.com/Alia5/steaminputdb.com/metrics"
	"github.com/Alia5/steaminputdb.com/middleware"
	"github.com/Alia5/steaminputdb.com/routes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/Alia5/steaminputdb.com/version"
	"github.com/alecthomas/kong"
	kongtoml "github.com/alecthomas/kong-toml"
	kongyaml "github.com/alecthomas/kong-yaml"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

const baseURL = "https://api.steaminputdb.com"

//go:embed apidesc.md
var apiDescription string

func main() {

	userCfg := findUserConfig(os.Args[1:])
	jsonPaths, yamlPaths, tomlPaths := configCandidatePaths(userCfg)

	var cfg config.Config
	_ = kong.Parse(&cfg,
		kong.Name("steaminputdb"),
		kong.Description(fmt.Sprintf("SteamInputDB Server - v%s", version.Version)),
		kong.UsageOnError(),
		// Load configuration from JSON/YAML/TOML in priority order; flags/env override config values.
		kong.Configuration(kong.JSON, jsonPaths...),
		kong.Configuration(kongyaml.Loader, yamlPaths...),
		kong.Configuration(kongtoml.Loader, tomlPaths...),
	)

	logging.SetupDefault(cfg.LogLevel)

	if config.Parsed.JWTSecret == "" || config.Parsed.JWTSecret == "testenv" {
		randomBytes := make([]byte, 32)
		_, err := rand.Read(randomBytes)
		if err != nil {
			slog.Error("Failed to generate random JWT secret", "error", err)
			os.Exit(1)
		}
		cfg.JWTSecret = fmt.Sprintf("%x", randomBytes)
	}
	config.Parsed = cfg

	dal, err := db.Init(cfg.DB)
	if err != nil {
		slog.Error("Failed to initialize database", "error", err)
		os.Exit(1)
		return
	}

	steamapi.DefaultClient = steamapi.NewClient(cfg.SteamAPIKey)

	metricsMux := http.NewServeMux()
	metricsSrv := http.Server{
		Addr: cfg.Metrics.ListenAddress,
		Handler: middleware.With(
			metricsMux,
			logging.Middleware,
			metrics.Middleware,
			// routes.UnregisteredMiddleware,
		),
	}
	metricsMux.Handle("GET /metrics", promhttp.Handler())

	schemaPrefix := "#/components/schemas/"
	schemasPath := "/schemas"

	registry := huma.NewMapRegistry(schemaPrefix, huma.DefaultSchemaNamer)

	docAPISrvs := []*huma.Server{}

	if cfg.API.PublicAddress != "" {
		docAPISrvs = append(docAPISrvs, &huma.Server{
			URL:         cfg.API.PublicAddress,
			Description: "SteamInputDB Live API",
		})
	} else {
		docAPISrvs = append(docAPISrvs, &huma.Server{
			URL:         baseURL,
			Description: "SteamInputDB Live API",
		})
	}

	apiClickable := formatClickableAddr(cfg.API.ListenAddress)
	if strings.Contains(apiClickable, "localhost") && os.Getenv("DEV") == "1" {
		docAPISrvs = append([]*huma.Server{{
			URL:         apiClickable,
			Description: "Local API",
		}}, docAPISrvs...)
	}

	apiMux := http.NewServeMux()
	hAPI := humago.New(apiMux, huma.Config{
		OpenAPI: &huma.OpenAPI{
			OpenAPI: "3.1.0",
			Info: &huma.Info{
				Title:       "SteamInputDB API",
				Description: apiDescription,
				License: &huma.License{
					Name:       "GNU Affero General Public License v3.0",
					URL:        "https://www.gnu.org/licenses/agpl-3.0.en.html",
					Identifier: "AGPL-3.0",
				},
				Version: version.Version,
			},
			Components: &huma.Components{
				Schemas: registry,
				SecuritySchemes: map[string]*huma.SecurityScheme{
					// makes doc pages work this way..., flow is set later
					"Steam OpenID": {
						Type:             "oauth2",
						OpenIDConnectURL: "https://steamcommunity.com/openid/.well-known/openid-configuration",
						Scheme:           "OAuth",
						Name:             "Steam Auth",
						In:               "query",
						Flows: &huma.OAuthFlows{
							Implicit: &huma.OAuthFlow{},
						},
					},
				},
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
					if sw, ok := c.BodyWriter().(*api.StatusWriter); ok {
						sw.Error = err
					}
				}
				return v, nil
			},
		},
	})
	o := hAPI.OpenAPI()
	flow := o.Components.SecuritySchemes["Steam OpenID"].Flows.Implicit

	authAddr := cfg.API.PublicAddress
	if authAddr == "" {
		authAddr = "http://localhost:8889"
	}

	flow.AuthorizationURL = steam.OpenIDAuthorizationURL(fmt.Sprintf("%s/v1/steam/login", authAddr))
	flow.TokenURL = fmt.Sprintf("%s/v1/steam/login", authAddr)

	apiSrv := http.Server{
		Addr: cfg.API.ListenAddress,
		Handler: middleware.With(
			apiMux,
			logging.Middleware,
			cors.New(cors.Options{
				AllowedOrigins:   strings.Split(cfg.API.CorsOrigins, ","),
				AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders:   []string{"*"},
				AllowCredentials: true,
			}).Handler,
			metrics.Middleware,
			routes.UnregisteredMiddleware(hAPI),
		),
	}

	hAPI.Adapter().Handle(&huma.Operation{
		Method: http.MethodGet,
		Path:   "/docs",
	}, func(ctx huma.Context) {
		ctx.SetHeader("Content-Type", "text/html")
		_, _ = ctx.BodyWriter().Write([]byte(`<!doctype html>
			<html>
			<head>
				<title>SteamInputDB API</title>
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

	api.RegisterAPI(hAPI, dal)

	// use kong for parsing, ignore humas config parser
	cli := humacli.New(func(h humacli.Hooks, _ *struct{}) {
		closed := false
		h.OnStart(func() {

			if os.Getenv("DEV") == "1" {
				yml, err := hAPI.OpenAPI().YAML()
				if err != nil {
					slog.Error("failed to generate OpenAPI YAML", "err", err)
				}

				err = os.WriteFile("../openapi.yaml", yml, 0644)
				if err != nil {
					slog.Error("failed to write OpenAPI YAML to file", "err", err)
				} else {
					slog.Info("wrote OpenAPI YAML to ../openapi.yaml")
				}
			}

			var wg sync.WaitGroup
			servers := []*http.Server{
				&metricsSrv,
				&apiSrv,
			}
			for _, srv := range servers {
				wg.Add(1)
				go func(s *http.Server) {
					defer wg.Done()
					slog.Info("Starting Server", "addr", s.Addr, "url", formatClickableAddr(s.Addr))
					if err := s.ListenAndServe(); err != nil && !closed {
						slog.Error("server error", "addr", s.Addr, "err", err)
						os.Exit(1)
					}
				}(srv)
			}
			wg.Wait()
		})

		h.OnStop(func() {
			closed = true
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err = metricsSrv.Shutdown(ctx)
			if err != nil {
				slog.Error("error shutting down metrics server", "err", err)
			}
			err = apiSrv.Shutdown(ctx)
			if err != nil {
				slog.Error("error shutting down API server", "err", err)
			}
		})

	})

	cli.Run()
	os.Exit(0)
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
	return os.Getenv("CONFIG")
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
	for _, base := range []string{"github.com/Alia5/steaminputdb.com", "steaminputdb", "config"} {
		add(&jsonPaths, filepath.Join(wd, base+".json"))
		add(&yamlPaths, filepath.Join(wd, base+".yaml"))
		add(&yamlPaths, filepath.Join(wd, base+".yml"))
		add(&tomlPaths, filepath.Join(wd, base+".toml"))
	}

	return
}

func formatClickableAddr(addr string) string {
	host := addr
	if strings.HasPrefix(addr, ":") {
		host = "localhost" + addr
	} else if after, ok := strings.CutPrefix(addr, "0.0.0.0:"); ok {
		host = "localhost:" + after
	} else if after, ok := strings.CutPrefix(addr, "[::]:"); ok {
		host = "localhost:" + after
	}
	return "http://" + host
}
