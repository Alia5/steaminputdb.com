// Package login provides Steam OpenID authentication and JWT generation.
package login

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/api/steam/user"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/golang-jwt/jwt/v5"
)

const steamLoginURL = "https://steamcommunity.com/openid/login"

type OpenIDRequest struct {
	Body *OpenIDBody
	OpenIDBody
}

type OpenIDBody struct {
	NS            string `json:"ns" query:"openid.ns"`
	Mode          string `json:"mode" enum:"id_res" query:"openid.mode"`
	OpEndpoint    string `json:"op_endpoint" query:"openid.op_endpoint"`
	ClaimedID     string `json:"claimed_id" query:"openid.claimed_id" minLength:"1"`
	Identity      string `json:"identity" query:"openid.identity"`
	ReturnTo      string `json:"return_to" query:"openid.return_to"`
	ResponseNonce string `json:"response_nonce" query:"openid.response_nonce"`
	AssocHandle   string `json:"assoc_handle" query:"openid.assoc_handle"`
	Signed        string `json:"signed" query:"openid.signed"`
	Sig           string `json:"sig" query:"openid.sig"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	user.PlayerInfo
}

type LoginResponse struct {
	user.PlayerInfo
}

type Response struct {
	Body        LoginResponse
	TokenCookie http.Cookie `header:"Set-Cookie"`
	Status      int
	URL         string `header:"Location"`
}

const jwtValidity = time.Hour * 24 * 30

func RegisterWithURL(a huma.API, steamURL string) {
	registerRoutes(a, steamURL)
}

func RegisterRoutes(a huma.API) {
	registerRoutes(a, steamLoginURL)
}

// TODO: create private endpoint for docs callback

func registerRoutes(a huma.API, loginURL string) {
	handler := handler(loginURL)
	huma.Register(
		a,
		huma.Operation{
			Method:  http.MethodPost,
			Path:    "/v1/steam/login",
			Tags:    []string{"steam", "auth"},
			Summary: "Log in with Steam",
			Description: `Authenticate user via Steam OpenID and return JWT token  
			Wrapper endpoint for SSR frontend`,
			Errors: []int{
				http.StatusUnauthorized, http.StatusBadGateway,
			},
		},
		handler,
	)
	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/login",
			Tags:        []string{"steam", "auth"},
			Summary:     "Log in with Steam",
			Description: `Authenticate user via Steam OpenID and return JWT token in Cookie`,
			Errors: []int{
				http.StatusUnauthorized, http.StatusBadGateway,
			},
		},
		handler,
	)
}

func handler(loginURL string) func(ctx context.Context, req *OpenIDRequest) (*Response, error) {
	return func(ctx context.Context, req *OpenIDRequest) (*Response, error) {
		params := url.Values{}

		params.Set("openid.mode", "check_authentication")

		if req.Body != nil && req.Body.NS != "" {
			params.Set("openid.ns", req.Body.NS)
			params.Set("openid.op_endpoint", req.Body.OpEndpoint)
			params.Set("openid.claimed_id", req.Body.ClaimedID)
			params.Set("openid.identity", req.Body.Identity)
			params.Set("openid.return_to", req.Body.ReturnTo)
			params.Set("openid.response_nonce", req.Body.ResponseNonce)
			params.Set("openid.assoc_handle", req.Body.AssocHandle)
			params.Set("openid.signed", req.Body.Signed)
			params.Set("openid.sig", req.Body.Sig)
		} else {
			params.Set("openid.ns", req.NS)
			params.Set("openid.op_endpoint", req.OpEndpoint)
			params.Set("openid.claimed_id", req.ClaimedID)
			params.Set("openid.identity", req.Identity)
			params.Set("openid.return_to", req.ReturnTo)
			params.Set("openid.response_nonce", req.ResponseNonce)
			params.Set("openid.assoc_handle", req.AssocHandle)
			params.Set("openid.signed", req.Signed)
			params.Set("openid.sig", req.Sig)
		}

		var claimedID string
		if req.Body != nil {
			claimedID = req.Body.ClaimedID
		} else {
			claimedID = req.ClaimedID
		}

		parts := strings.Split(claimedID, "/")
		if len(parts) == 0 {
			return nil, huma.Error401Unauthorized("invalid claimed_id")
		}
		steamID := parts[len(parts)-1]

		steamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, loginURL, strings.NewReader(params.Encode()))
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to create request")
		}
		steamReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(steamReq)
		if err != nil {
			if errors.Is(err, steamapi.ErrRequest) {
				return nil, huma.Error502BadGateway("failed to get steam user info: %v", err)
			}
			return nil, huma.Error401Unauthorized("failed to verify with steam")
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				slog.Error("failed to close steam response body", "error", err)
			}
		}()

		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to read response")
		}

		if !strings.Contains(string(respBody), "is_valid:true") {
			return nil, huma.Error401Unauthorized("steam verification failed")
		}

		playerSummaries, err := steamapi.DefaultClient.GetPlayerSummaries(ctx, steamID)
		if err != nil {
			if errors.Is(err, steamapi.ErrRequest) {
				return nil, huma.Error502BadGateway("failed to get steam user info: %v", err)
			}
			return nil, err
		}

		if len(playerSummaries.Response.Players) == 0 {
			return nil, huma.Error404NotFound("steam user not found")
		}

		playerInfo := user.PlayerInfo{
			CommunityVisibilityState: playerSummaries.Response.Players[0].Communityvisibilitystate,
			PersonaName:              playerSummaries.Response.Players[0].Personaname,
			ProfileURL:               playerSummaries.Response.Players[0].Profileurl,
			Avatar:                   playerSummaries.Response.Players[0].Avatar,
			AvatarMedium:             playerSummaries.Response.Players[0].Avatarmedium,
			AvatarFull:               playerSummaries.Response.Players[0].Avatarfull,
			AvatarHash:               playerSummaries.Response.Players[0].Avatarhash,
			LastLogOff:               time.Unix(int64(playerSummaries.Response.Players[0].Lastlogoff), 0),
			PrimaryClanID:            playerSummaries.Response.Players[0].Primaryclanid,
			TimeCreated:              time.Unix(int64(playerSummaries.Response.Players[0].Timecreated), 0),
			LocCountryCode:           playerSummaries.Response.Players[0].Loccountrycode,
		}
		claims := JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject:   steamID,
				Issuer:    "steaminputdb",
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtValidity)),
			},
			PlayerInfo: playerInfo,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(config.Parsed.JWTSecret))
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to generate token")
		}

		HTTPSOnly := os.Getenv("DEV") != "1"

		domain := config.Parsed.API.PublicAddress
		domain = strings.Split(domain, ":")[0]
		domain = strings.TrimPrefix(domain, "http://")
		domain = strings.TrimPrefix(domain, "https://")

		cookie := http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   int(jwtValidity.Seconds()),
			Path:     "/",
			Domain:   fmt.Sprintf(".%s", domain),
			Secure:   HTTPSOnly,
		}

		if req.Body == nil {
			return &Response{
				Status:      http.StatusFound,
				URL:         "http://localhost:8889/v1/ping#access_token=" + tokenString,
				TokenCookie: cookie,
			}, nil
		}

		return &Response{
			Status: http.StatusCreated,
			Body: LoginResponse{
				PlayerInfo: playerInfo,
			},
			TokenCookie: cookie,
		}, nil
	}
}
