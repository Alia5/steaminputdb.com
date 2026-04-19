package login_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Alia5/steaminputdb.com/api/steam/login"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/steamapi"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
	"github.com/stretchr/testify/assert"
)

func TestSteamLogin(t *testing.T) {

	type testCase struct {
		name             string
		setupMock        func() (*httptest.Server, *httptest.Server)
		requestBody      map[string]any
		expectedStatus   int
		expectedResponse string
	}

	testCases := []testCase{
		{
			name: "SUCCESS",
			setupMock: func() (*httptest.Server, *httptest.Server) {
				openIDServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ns:http://specs.openid.net/auth/2.0\nis_valid:true\n"))
				}))
				steamAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					resp := steamapi.PlayerSummaries{
						Response: steamapi.Response{
							Players: []steamapi.Players{
								{
									Steamid:                  "13371337696942069",
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
				return openIDServer, steamAPIServer
			},
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "id_res",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "https://steamcommunity.com/openid/id/13371337696942069",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "5z/7Gwh1MhMbJ8JeU2pPNSit9l8=",
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "INVALID_MODE",
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "invalid_mode",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "https://steamcommunity.com/openid/id/13371337696942069",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "5z/7Gwh1MhMbJ8JeU2pPNSit9l8=",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `{
				"title": "Unprocessable Entity",
				"status": 422,
				"detail": "validation failed",
				"errors": [
					{
						"message": "expected value to be one of \"id_res\"",
						"location": "body.mode",
						"value": "invalid_mode"
					}
				]
			}`,
		},
		{
			name: "INVALID_SIGNATURE",
			setupMock: func() (*httptest.Server, *httptest.Server) {
				openIDServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ns:http://specs.openid.net/auth/2.0\nis_valid:false\n"))
				}))
				return openIDServer, nil
			},
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "id_res",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "https://steamcommunity.com/openid/id/13371337696942069",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "invalid_signature",
			},
			expectedStatus:   http.StatusUnauthorized,
			expectedResponse: `{"detail":"steam verification failed", "status":401, "title":"Unauthorized"}`,
		},
		{
			name: "EMPTY_CLAIMED_ID",
			setupMock: func() (*httptest.Server, *httptest.Server) {
				openIDServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ns:http://specs.openid.net/auth/2.0\nis_valid:true\n"))
				}))
				return openIDServer, nil
			},
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "id_res",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "5z/7Gwh1MhMbJ8JeU2pPNSit9l8=",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `{
				"title": "Unprocessable Entity",
				"status": 422,
				"detail": "validation failed",
				"errors": [
					{
						"message": "expected length >= 1",
						"location": "body.claimed_id",
						"value": ""
					}
				]
			}`,
		},
		{
			name: "STEAM_API_NO_RESPONSE",
			setupMock: func() (*httptest.Server, *httptest.Server) {
				openIDServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ns:http://specs.openid.net/auth/2.0\nis_valid:true\n"))
				}))
				steamAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusServiceUnavailable)
				}))
				return openIDServer, steamAPIServer
			},
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "id_res",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "https://steamcommunity.com/openid/id/13371337696942069",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "5z/7Gwh1MhMbJ8JeU2pPNSit9l8=",
			},
			expectedStatus: http.StatusBadGateway,
		},
		{
			name: "STEAM_API_GARBAGE_DATA",
			setupMock: func() (*httptest.Server, *httptest.Server) {
				openIDServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("ns:http://specs.openid.net/auth/2.0\nis_valid:true\n"))
				}))
				steamAPIServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"response":{"players":[]}}`))
				}))
				return openIDServer, steamAPIServer
			},
			requestBody: map[string]any{
				"ns":             "http://specs.openid.net/auth/2.0",
				"mode":           "id_res",
				"op_endpoint":    "https://steamcommunity.com/openid/login",
				"claimed_id":     "https://steamcommunity.com/openid/id/13371337696942069",
				"identity":       "https://steamcommunity.com/openid/id/13371337696942069",
				"return_to":      "http://localhost:5173/login/callback",
				"response_nonce": "2026-01-31T23:46:40ZfbncPsCvbGlcAnm7O7AOGbvecyg=",
				"assoc_handle":   "1234567890",
				"signed":         "signed,op_endpoint,claimed_id,identity,return_to,response_nonce,assoc_handle",
				"sig":            "5z/7Gwh1MhMbJ8JeU2pPNSit9l8=",
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

			var openIDURL, steamAPIURL string
			var openIDServer, steamAPIServer *httptest.Server

			if tc.setupMock != nil {
				openIDServer, steamAPIServer = tc.setupMock()
				if openIDServer != nil {
					defer openIDServer.Close()
					openIDURL = openIDServer.URL
				}
				if steamAPIServer != nil {
					defer steamAPIServer.Close()
					steamAPIURL = steamAPIServer.URL
					steamapi.DefaultClient = steamapi.NewClientWithBaseURL("test-key", steamAPIURL)
				}
			}

			login.RegisterWithURL(api, dal, openIDURL)

			resp := api.Post("/v1/steam/login", tc.requestBody)
			assert.Equal(t, tc.expectedStatus, resp.Code)

			if tc.expectedResponse != "" {
				assert.JSONEq(t, tc.expectedResponse, resp.Body.String())
			}

		})
	}

}
