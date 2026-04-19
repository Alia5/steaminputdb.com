package appinfo_test

import (
	"os"
	"strings"
	"testing"

	"github.com/Alia5/steaminputdb.com/api/steam/appinfo"
	"github.com/Alia5/steaminputdb.com/steamapi"
	sidbtest "github.com/Alia5/steaminputdb.com/testing"
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

func BenchmarkAppInfo(b *testing.B) {

	type testCase struct {
		name        string
		useMemCache bool
		freshDB     bool
	}

	testCases := []testCase{
		{name: "With MemCache", useMemCache: true},
		{name: "Without MemCache", useMemCache: false},
		{name: "With Fresh DB every time", useMemCache: false, freshDB: true},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {

			api, dal, err := sidbtest.MockAPI(b)
			require.NoError(b, err)
			appinfo.RegisterRoute(api, dal, tc.useMemCache)
			//warmup
			_ = api.Get("/v1/steam/appinfo?app_id=250900")
			for b.Loop() {
				if tc.freshDB {
					b.StopTimer()
					api, dal, err = sidbtest.MockAPI(b)
					require.NoError(b, err)
					appinfo.RegisterRoute(api, dal, tc.useMemCache)
					b.StartTimer()
				}
				resp := api.Get("/v1/steam/appinfo?app_id=250900")
				b.StopTimer()
				if resp.Body == nil {
					b.Fatal("response body is nil")
				}
				if resp.Code != 200 {
					b.Fatalf("expected status 200, got %d", resp.Code)
				}
				b.StartTimer()
			}

		})
	}

}
