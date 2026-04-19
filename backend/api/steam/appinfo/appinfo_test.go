package appinfo_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Alia5/steaminputdb.com/api/steam/appinfo"
	"github.com/Alia5/steaminputdb.com/steamapi"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	key := os.Getenv("STEAM_API_KEY")
	if key == "" {
		if b, err := os.ReadFile(".env"); err == nil {
			for _, line := range strings.Split(string(b), "\n") {
				if k, v, ok := strings.Cut(line, "="); ok && k == "STEAM_API_KEY" {
					key = v
				}
			}
		}
	}
	if key != "" {
		steamapi.DefaultClient = steamapi.NewClient(key)
	}
}

func TestSteamAppInfo(t *testing.T) {
	type testCase struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
		contains       string
		setup          func(t *testing.T)
	}

	testCases := []testCase{
		{
			name:           "SUCCESS",
			path:           "/v1/steam/appinfo?app_id=250900",
			expectedStatus: http.StatusOK,
			contains:       "The Binding of Isaac: Rebirth",
		},
		{
			name:           "MISSING_APP_ID",
			path:           "/v1/steam/appinfo",
			expectedStatus: http.StatusUnprocessableEntity, contains: "app_id",
		},
		{
			name:           "APP_NOT_FOUND",
			path:           "/v1/steam/appinfo?app_id=999999",
			expectedStatus: http.StatusNotFound,
			expectedBody: `{
				"title": "Not Found",
				"status": 404,
				"detail": "item not found"
			}`,
		},
		{
			name:           "STEAM_API_ERROR",
			path:           "/v1/steam/appinfo?app_id=250900",
			expectedStatus: http.StatusBadGateway,
			contains:       "failed to get steam app info",
			setup: func(t *testing.T) {
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "", http.StatusServiceUnavailable)
				}))
				t.Cleanup(srv.Close)
				orig := steamapi.DefaultClient
				steamapi.DefaultClient = steamapi.NewClientWithBaseURL("test-key", srv.URL)
				t.Cleanup(func() { steamapi.DefaultClient = orig })
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setup != nil {
				tc.setup(t)
			}

			api, dal, err := sidbtest.MockAPI(t)
			require.NoError(t, err)
			appinfo.RegisterRoute(api, dal)
			resp := api.Get(tc.path)

			assert.Equal(t, tc.expectedStatus, resp.Code, "body: %s", resp.Body.String())
			if tc.expectedBody != "" {
				assert.JSONEq(t, tc.expectedBody, resp.Body.String())
			}
			if tc.contains != "" {
				assert.Contains(t, resp.Body.String(), tc.contains)
			}
		})
	}
}
