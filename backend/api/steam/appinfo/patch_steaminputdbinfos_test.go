package appinfo_test

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Alia5/steaminputdb.com/api/steam/appinfo"
	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/Alia5/steaminputdb.com/db/models"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPatchSteamInputDBInfos(t *testing.T) {
	type testCase struct {
		name           string
		path           string
		body           string
		expectedStatus int
		expectedBody   string
		setup          func(t *testing.T, api humatest.TestAPI, dal db.DAL)
	}

	const path = "/v1/steam/appinfo/999999/steaminputdbinfos"

	claims := jwt.MapClaims{
		"sub": "76561197997352479",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(config.Parsed.JWTSecret))
	cookieHeader := "Cookie: token=" + tokenString

	fullBody := `{
		"controller_support_rating": 5,
		"controller_support_notes": "notes",
		"mixed_input": {
			"type": 3,
			"glyph_flicker": true,
			"notes": "mixed notes",
			"mixed_input_mod_urls": ["https://example.com/a", "", "  ", "https://example.com/b"]
		},
		"glyphs": {
			"autodetect": true,
			"manual_select": true,
			"notes": "glyph notes",
			"controllers": [
				{"controller_type": "controller_ps5", "notes": "ps5"},
				{"controller_type": "controller_ps5", "notes": "duplicate"},
				{"controller_type": "controller_xboxone"}
			]
		},
		"steaminputapi_support": {
			"glyphs": [1, 2],
			"camera_support": 3,
			"pixels_per_360": "1234",
			"support_tags": [2, 3],
			"notes": "siapi notes"
		},
		"hw_features": [
			{"feature": 2, "controller_family": 0, "notes": "rumble"},
			{"feature": 2, "controller_family": 0, "notes": "duplicate"},
			{"feature": 3, "controller_family": 2}
		]
	}`

	testCases := []testCase{
		{
			name:           "USER_NOT_FOUND",
			path:           path,
			body:           fullBody,
			expectedStatus: http.StatusForbidden,
			expectedBody: `{
				"title": "Forbidden",
				"status": 403,
				"detail": "authentication error"
			}`,
		},
		{
			name:           "INSUFFICIENT_PERMISSIONS",
			path:           path,
			body:           fullBody,
			expectedStatus: http.StatusForbidden,
			expectedBody: `{
				"title": "Forbidden",
				"status": 403,
				"detail": "authentication error"
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
				})
				require.NoError(t, err)
			},
		},
		{
			name:           "APP_NOT_FOUND",
			path:           "/v1/steam/appinfo/999998/steaminputdbinfos",
			body:           fullBody,
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"title": "Not Found",
				"status": 404,
				"detail": "item not found"
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
			},
		},
		{
			name:           "SUCCESS_ADMIN",
			path:           path,
			body:           fullBody,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 5,
					"controller_support_notes": "notes",
					"mixed_input": {
						"type": 3,
						"glyph_flicker": true,
						"notes": "mixed notes",
						"mixed_input_mod_urls": ["https://example.com/a", "https://example.com/b"]
					},
					"glyphs": {
						"autodetect": true,
						"manual_select": true,
						"notes": "glyph notes",
						"controllers": [
							{"controller_type": "controller_ps5", "notes": "ps5"},
							{"controller_type": "controller_xboxone"}
						]
					},
					"steaminputapi_support": {
						"glyphs": [1, 2],
						"camera_support": 3,
						"pixels_per_360": "1234",
						"support_tags": [2, 3],
						"notes": "siapi notes"
					},
					"hw_features": [
						{"feature": 2, "controller_family": 0, "notes": "rumble"},
						{"feature": 3, "controller_family": 2}
					]
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
			},
		},
		{
			name:           "SUCCESS_TRUSTED_MOD",
			path:           path,
			body:           `{"controller_support_rating": 4}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 4
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					TrustedMod:  true,
				})
				require.NoError(t, err)
			},
		},
		{
			name:           "OMITTED_FIELDS_KEEP_EXISTING_INFOS",
			path:           path,
			body:           `{"controller_support_rating": 4}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 4,
					"controller_support_notes": "notes",
					"mixed_input": {
						"type": 3,
						"glyph_flicker": true,
						"notes": "mixed notes",
						"mixed_input_mod_urls": ["https://example.com/a", "https://example.com/b"]
					},
					"glyphs": {
						"autodetect": true,
						"manual_select": true,
						"notes": "glyph notes",
						"controllers": [
							{"controller_type": "controller_ps5", "notes": "ps5"},
							{"controller_type": "controller_xboxone"}
						]
					},
					"steaminputapi_support": {
						"glyphs": [1, 2],
						"camera_support": 3,
						"pixels_per_360": "1234",
						"support_tags": [2, 3],
						"notes": "siapi notes"
					},
					"hw_features": [
						{"feature": 2, "controller_family": 0, "notes": "rumble"},
						{"feature": 3, "controller_family": 2}
					]
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
				resp := api.Patch(path, cookieHeader, strings.NewReader(fullBody))
				require.Equal(t, http.StatusOK, resp.Code, "body: %s", resp.Body.String())
			},
		},
		{
			name: "EMPTY_SECTIONS_DO_NOT_CREATE_SECTIONS",
			path: path,
			body: `{
				"controller_support_rating": 4,
				"mixed_input": {"type": 0, "glyph_flicker": null, "notes": "", "mixed_input_mod_urls": []},
				"glyphs": {"autodetect": null, "manual_select": null, "notes": "", "controllers": []},
				"steaminputapi_support": {"glyphs": [], "camera_support": 0, "pixels_per_360": null, "notes": "", "support_tags": []},
				"hw_features": []
			}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 4
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
			},
		},
		{
			name: "GLYPH_TAGS_WITHOUT_GLYPHS_KEEP_GLYPH_INFO",
			path: path,
			body: `{
				"controller_support_rating": 5,
				"steaminputapi_support": {"glyphs": [3], "camera_support": 0, "pixels_per_360": ""}
			}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 5,
					"controller_support_notes": "notes",
					"mixed_input": {
						"type": 3,
						"glyph_flicker": true,
						"notes": "mixed notes",
						"mixed_input_mod_urls": ["https://example.com/a", "https://example.com/b"]
					},
					"glyphs": {
						"autodetect": true,
						"manual_select": true,
						"notes": "glyph notes",
						"controllers": [
							{"controller_type": "controller_ps5", "notes": "ps5"},
							{"controller_type": "controller_xboxone"}
						]
					},
					"steaminputapi_support": {
						"glyphs": [3],
						"camera_support": 0,
						"support_tags": [2, 3]
					},
					"hw_features": [
						{"feature": 2, "controller_family": 0, "notes": "rumble"},
						{"feature": 3, "controller_family": 2}
					]
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
				resp := api.Patch(path, cookieHeader, strings.NewReader(fullBody))
				require.Equal(t, http.StatusOK, resp.Code, "body: %s", resp.Body.String())
			},
		},
		{
			name: "ZERO_VALUES_AND_EMPTY_ARRAYS_CLEAR_EXISTING_INFOS",
			path: path,
			body: `{
				"controller_support_rating": 0,
				"controller_support_notes": "",
				"mixed_input": {"type": 0, "glyph_flicker": false, "notes": "", "mixed_input_mod_urls": [""]},
				"glyphs": {"autodetect": false, "manual_select": false, "notes": "", "controllers": []},
				"steaminputapi_support": {"glyphs": [], "camera_support": 0, "pixels_per_360": "", "notes": "", "support_tags": []},
				"hw_features": []
			}`,
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"app_id": 999999,
				"name": "Test App",
				"store_url_path": "",
				"type": "",
				"steaminputdb_info": {
					"controller_support_rating": 0,
					"controller_support_notes": "",
					"mixed_input": {
						"type": 0,
						"glyph_flicker": false
					},
					"glyphs": {
						"autodetect": false,
						"manual_select": false
					},
					"steaminputapi_support": {
						"camera_support": 0
					}
				}
			}`,
			setup: func(t *testing.T, api humatest.TestAPI, dal db.DAL) {
				err := dal.SteamUser().Insert(context.Background(), &models.SteamUser{
					SteamID:     76561197997352479,
					PersonaName: "TestUser",
					IsAdmin:     true,
				})
				require.NoError(t, err)
				resp := api.Patch(path, cookieHeader, strings.NewReader(fullBody))
				require.Equal(t, http.StatusOK, resp.Code, "body: %s", resp.Body.String())
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			api, dal, err := sidbtest.MockAPI(t)
			require.NoError(t, err)
			err = dal.AppInfo().Insert(context.Background(), &models.AppInfo{
				AppID: 999999,
				Name:  "Test App",
			})
			require.NoError(t, err)
			appinfo.RegisterRoute(api, dal)

			if tc.setup != nil {
				tc.setup(t, api, dal)
			}

			resp := api.Patch(tc.path, cookieHeader, strings.NewReader(tc.body))

			assert.Equal(t, tc.expectedStatus, resp.Code, "body: %s", resp.Body.String())
			assert.JSONEq(t, tc.expectedBody, resp.Body.String())

			if tc.expectedStatus != http.StatusOK {
				return
			}

			getResp := api.Get("/v1/steam/appinfo/999999?steaminputdb_info=true")
			assert.Equal(t, http.StatusOK, getResp.Code, "body: %s", getResp.Body.String())
			assert.JSONEq(t, tc.expectedBody, getResp.Body.String())
		})
	}
}
