package testing

import (
	"testing"

	"github.com/Alia5/steaminputdb.com/buddy-app/config"
	"github.com/Alia5/steaminputdb.com/buddy-app/db"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func NewMemDB(tb testing.TB) (*db.DAL, error) {
	tb.Helper()

	dal, _, err := db.Open("file::memory:")
	if err != nil {
		return nil, err
	}
	return dal, nil
}

func MockConfig(tb testing.TB) *config.Config {
	tb.Helper()
	return &config.Config{
		LogLevel:    "debug",
		TrayDisplay: false,
		API: config.API{
			ListenAddress: "localhost:5119",
			CORSOrigins:   "*",
		},
		Steam: config.Steam{
			CEFRemoteDebugPort: 8080,
		},
	}
}

func MockTestParams(t testing.TB) (humatest.TestAPI, *db.DAL, *config.Config, error) {
	t.Helper()
	_, api := humatest.New(TWithSuppressedHumaLogs(t))
	cfg := MockConfig(t)
	dal, err := NewMemDB(t)
	if err != nil {
		return nil, nil, nil, err
	}
	return api, dal, cfg, nil
}

// TWithSuppressedHumaLogs returns a testing.TB that suppresses logs from Huma during tests.
//
// Usage:
//
//	_, api := humatest.New(buddyTesting.TWithSuppressedHumaLogs(t))
func TWithSuppressedHumaLogs(tb testing.TB) testing.TB {
	return &suppressedHumaLogs{TB: tb}
}

type suppressedHumaLogs struct {
	testing.TB
}

func (s *suppressedHumaLogs) Log(_ ...any) {
	// no-op
}
