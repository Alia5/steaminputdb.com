package settings_test

import (
	"strings"
	"testing"

	"github.com/Alia5/steaminputdb.com/buddy-app/api/settings"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	buddyTesting "github.com/Alia5/steaminputdb.com/buddy-app/testing"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/stretchr/testify/assert"
)

func TestGetSettings(t *testing.T) {

	type testCase struct {
		name             string
		expectedStatus   int
		expectedResponse string
		setup            func(d *db.DAL) error
	}

	testCases := []testCase{
		{
			name:           "Get default settings",
			expectedStatus: 200,
			expectedResponse: `{
				"addDesktopUIEntries": true,
				"addBigPictureUIEntries": true,
				"desktopUseSteamBrowser": false,
				"steamWaitTimeout": "1m0s",
				"autoStart": false
			}`,
		},
		{
			name: "Get settings after update",
			setup: func(d *db.DAL) error {
				s, err := d.Settings.Get(t.Context())
				if err != nil {
					return err
				}
				s.AddDesktopUIEntries = new(false)
				s.AddBigPictureUIEntries = new(false)
				s.DesktopUseSteamBrowser = new(true)
				_, err = d.Settings.Update(t.Context(), s)
				return err
			},
			expectedStatus: 200,
			expectedResponse: `{
				"addDesktopUIEntries": false,
				"addBigPictureUIEntries": false,
				"desktopUseSteamBrowser": true,
				"steamWaitTimeout": "1m0s",
				"autoStart": false
			}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, api := humatest.New(buddyTesting.TWithSuppressedHumaLogs(t))
			dal, err := buddyTesting.NewMemDB(t)
			if err != nil {
				t.Fatalf("failed to create mock DB: %v", err)
			}

			settings.RegisterRoutes(api, dal, buddyTesting.MockConfig(t))

			if tc.setup != nil {
				err := tc.setup(dal)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}
			resp := api.Get("/v1/settings")
			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedResponse, resp.Body.String())
		})
	}
}

func TestUpdateSettings(t *testing.T) {

	type testCase struct {
		name             string
		body             string
		expectedStatus   int
		expectedResponse string
		setup            func(d *db.DAL) error
		postAsserts      func(d *db.DAL)
	}

	testCases := []testCase{
		{
			name:           "Updates all settings",
			expectedStatus: 200,
			body: `{
				"addDesktopUIEntries": false,
				"addBigPictureUIEntries": false,
				"desktopUseSteamBrowser": true,
				"steamWaitTimeout": "2m30s"
			}`,
			expectedResponse: `{
				"addDesktopUIEntries": false,
				"addBigPictureUIEntries": false,
				"desktopUseSteamBrowser": true,
				"steamWaitTimeout": "2m30s"
			}`,
			postAsserts: func(d *db.DAL) {
				s, err := d.Settings.Get(t.Context())
				assert.NoError(t, err)
				assert.Equal(t, new(false), s.AddDesktopUIEntries)
				assert.Equal(t, new(false), s.AddBigPictureUIEntries)
				assert.Equal(t, new(true), s.DesktopUseSteamBrowser)
			},
		},
		{
			name:           "Updates partial settings",
			expectedStatus: 200,
			setup: func(d *db.DAL) error {
				s, err := d.Settings.Get(t.Context())
				if err != nil {
					return err
				}
				s.AddDesktopUIEntries = new(false)
				s.DesktopUseSteamBrowser = new(true)
				_, err = d.Settings.Update(t.Context(), s)
				return err
			},
			body: `{
				"addDesktopUIEntries": true,
				"addBigPictureUIEntries": true,
				"desktopUseSteamBrowser": false,
				"steamWaitTimeout": "1m0s"
			}`,
			expectedResponse: `{
				"addDesktopUIEntries": true,
				"addBigPictureUIEntries": true,
				"desktopUseSteamBrowser": false,
				"steamWaitTimeout": "1m0s"
			}`,
			postAsserts: func(d *db.DAL) {
				s, err := d.Settings.Get(t.Context())
				assert.NoError(t, err)
				assert.Equal(t, new(true), s.AddDesktopUIEntries)
				assert.Equal(t, new(true), s.AddBigPictureUIEntries)
				assert.Equal(t, new(false), s.DesktopUseSteamBrowser)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, api := humatest.New(buddyTesting.TWithSuppressedHumaLogs(t))
			dal, err := buddyTesting.NewMemDB(t)
			if err != nil {
				t.Fatalf("failed to create mock DB: %v", err)
			}

			settings.RegisterRoutes(api, dal, buddyTesting.MockConfig(t))

			if tc.setup != nil {
				err := tc.setup(dal)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}
			resp := api.Put("/v1/settings", strings.NewReader(tc.body))
			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedResponse, resp.Body.String())
			if tc.postAsserts != nil {
				tc.postAsserts(dal)
			}
		})
	}
}

func TestResetSettings(t *testing.T) {

	type testCase struct {
		name             string
		expectedStatus   int
		expectedResponse string
		setup            func(d *db.DAL) error
		postAsserts      func(d *db.DAL)
	}

	testCases := []testCase{
		{
			name: "Resets to defaults after modification",
			setup: func(d *db.DAL) error {
				s, err := d.Settings.Get(t.Context())
				if err != nil {
					return err
				}
				s.AddDesktopUIEntries = new(false)
				s.AddBigPictureUIEntries = new(false)
				s.DesktopUseSteamBrowser = new(true)
				_, err = d.Settings.Update(t.Context(), s)
				return err
			},
			expectedStatus: 200,
			expectedResponse: `{
				"addDesktopUIEntries": true,
				"addBigPictureUIEntries": true,
				"desktopUseSteamBrowser": false,
				"steamWaitTimeout": "1m0s",
				"autoStart": false
			}`,
			postAsserts: func(d *db.DAL) {
				s, err := d.Settings.Get(t.Context())
				assert.NoError(t, err)
				assert.Equal(t, new(true), s.AddDesktopUIEntries)
				assert.Equal(t, new(true), s.AddBigPictureUIEntries)
				assert.Equal(t, new(false), s.DesktopUseSteamBrowser)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, api := humatest.New(buddyTesting.TWithSuppressedHumaLogs(t))
			dal, err := buddyTesting.NewMemDB(t)
			if err != nil {
				t.Fatalf("failed to create mock DB: %v", err)
			}

			settings.RegisterRoutes(api, dal, buddyTesting.MockConfig(t))

			if tc.setup != nil {
				err := tc.setup(dal)
				if err != nil {
					t.Fatalf("setup failed: %v", err)
				}
			}
			resp := api.Post("/v1/settings/reset")
			assert.Equal(t, tc.expectedStatus, resp.Code)
			assert.JSONEq(t, tc.expectedResponse, resp.Body.String())
			if tc.postAsserts != nil {
				tc.postAsserts(dal)
			}
		})
	}
}
