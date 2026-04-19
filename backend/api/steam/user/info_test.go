package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Alia5/steaminputdb.com/api/steam/user"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/steamapi"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestSteamUserInfo(t *testing.T) {

	type testCase struct {
		name             string
		setupMock        func() *httptest.Server
		setupToken       func() string
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{
			name: "SUCCESS",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					resp := steamapi.PlayerSummaries{
						Response: steamapi.Response{
							Players: []steamapi.Players{
								{
									Steamid:                  "76561197997352479",
									Communityvisibilitystate: 3,
									Personaname:              "TestUser",
									Profileurl:               "https://steamcommunity.com/id/testuser/",
									Avatar:                   "https://avatars.steamstatic.com/test.jpg",
									Avatarmedium:             "https://avatars.steamstatic.com/test_medium.jpg",
									Avatarfull:               "https://avatars.steamstatic.com/test_full.jpg",
									Avatarhash:               "testhash",
									Lastlogoff:               1234567890,
									Primaryclanid:            "123456789",
									Timecreated:              1234567890,
									Loccountrycode:           "US",
								},
							},
						},
					}
					json.NewEncoder(w).Encode(resp)
				}))
			},
			setupToken: func() string {
				claims := jwt.MapClaims{
					"sub": "76561197997352479",
					"exp": time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(config.Parsed.JWTSecret))
				return tokenString
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "STEAM_API_NO_RESPONSE",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
			},
			setupToken: func() string {
				claims := jwt.MapClaims{
					"sub": "76561197997352479",
					"exp": time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(config.Parsed.JWTSecret))
				return tokenString
			},
			expectedStatus: http.StatusBadGateway,
		},
		{
			name: "STEAM_USER_NOT_FOUND",
			setupMock: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					resp := steamapi.PlayerSummaries{
						Response: steamapi.Response{
							Players: []steamapi.Players{},
						},
					}
					json.NewEncoder(w).Encode(resp)
				}))
			},
			setupToken: func() string {
				claims := jwt.MapClaims{
					"sub": "76561197997352479",
					"exp": time.Now().Add(time.Hour).Unix(),
				}
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
				tokenString, _ := token.SignedString([]byte(config.Parsed.JWTSecret))
				return tokenString
			},
			expectedStatus:   http.StatusNotFound,
			expectedResponse: `{"detail":"steam user not found", "status":404, "title":"Not Found"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			origSteam := steamapi.DefaultClient
			steamapi.DefaultClient = steamapi.NewClient("")
			t.Cleanup(func() {
				steamapi.DefaultClient = origSteam
			})

			config.Parsed = config.Config{
				API: config.API{
					PublicAddress: "localhost:8889",
				},
			}

			api, dal, err := sidbtest.MockAPI(t)
			assert.NoError(t, err)

			var steamAPIURL string
			var steamAPIServer *httptest.Server

			if tc.setupMock != nil {
				steamAPIServer = tc.setupMock()
				defer steamAPIServer.Close()
				steamAPIURL = steamAPIServer.URL
				steamapi.DefaultClient = steamapi.NewClientWithBaseURL("test-key", steamAPIURL)
			}

			user.RegisterRoutes(api, dal)

			var token string
			if tc.setupToken != nil {
				token = tc.setupToken()
			}

			var cookieHeader string
			if token != "" {
				cookieHeader = "Cookie: token=" + token
			}

			var resp *httptest.ResponseRecorder
			if cookieHeader != "" {
				resp = api.Get("/v1/steam/userinfo", cookieHeader)
			} else {
				resp = api.Get("/v1/steam/userinfo")
			}

			assert.Equal(t, tc.expectedStatus, resp.Code)
			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, resp.Body.String())
			}

		})
	}

}
