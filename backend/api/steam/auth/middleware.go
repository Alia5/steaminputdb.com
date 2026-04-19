package auth

import (
	"log/slog"
	"net/http"

	"github.com/Alia5/steaminputdb.com/api/ctx"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ForceAuthMiddleware(a huma.API) func(c huma.Context, next func(huma.Context)) {
	return func(c huma.Context, next func(huma.Context)) {
		cookie, err := huma.ReadCookie(c, "token")
		if err != nil {
			err := huma.WriteErr(a, c, http.StatusUnauthorized, "missing token")
			if err != nil {
				slog.Error("failed to write error response", "error", err)
			}
			return
		}

		token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
			return []byte(config.Parsed.JWTSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
		if err != nil || !token.Valid {
			err := huma.WriteErr(a, c, http.StatusUnauthorized, "invalid token")
			if err != nil {
				slog.Error("failed to write error response", "error", err)
			}
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			err := huma.WriteErr(a, c, http.StatusUnauthorized, "invalid token claims")
			if err != nil {
				slog.Error("failed to write error response", "error", err)
			}
			return
		}

		steamID, ok := claims["sub"].(string)
		if !ok || steamID == "" {
			err := huma.WriteErr(a, c, http.StatusUnauthorized, "missing steamid")
			if err != nil {
				slog.Error("failed to write error response", "error", err)
			}
			return
		}

		c = huma.WithValue(c, ctx.KeySteamID, steamID)
		next(c)
	}
}

func ExtractSteamIDMiddleware(c huma.Context, next func(huma.Context)) {

	cookie, err := huma.ReadCookie(c, "token")
	if err != nil {
		next(c)
		return
	}

	token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (any, error) {
		return []byte(config.Parsed.JWTSecret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}))
	if err != nil || !token.Valid {
		slog.Debug("invalid token", "error", err)
		next(c)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Debug("invalid token claims")
		next(c)
		return
	}

	steamID, ok := claims["sub"].(string)
	if !ok || steamID == "" {
		slog.Debug("missing steamid in token claims")
		next(c)
		return
	}

	c = huma.WithValue(c, ctx.KeySteamID, steamID)
	next(c)
}
