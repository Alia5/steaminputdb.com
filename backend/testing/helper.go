package testing

import (
	"testing"
	"time"

	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func NewMemDB(tb testing.TB, shared bool) (db.DAL, error) {

	uri := "file::memory:"
	if shared {
		uri = uri + "?cache=shared"
	}
	return db.Init(config.DB{
		DatabaseURL: uri,
		// DatabaseURL:     "file:./testdb.sql?cache=shared",
		MaxOpenConns:    1,
		MaxIdleConns:    100,
		MaxConnLifetime: 5 * time.Minute,
		MaxConnIdleTime: 5 * time.Minute,
	})
}

func MockAPI(t testing.TB) (humatest.TestAPI, db.DAL, error) {
	t.Helper()
	_, api := humatest.New(TWithSuppressedHumaLogs(t))
	dal, err := NewMemDB(t, false)
	if err != nil {
		return nil, nil, err
	}
	return api, dal, nil
}

// TWithSuppressedHumaLogs returns a testing.TB that suppresses logs from Huma during tests.
//
// Usage:
//
// _, api := humatest.New(studioTesting.TWithSuppressedHumaLogs(t))
func TWithSuppressedHumaLogs(tb testing.TB) testing.TB {
	return &suppressedHumaLogs{TB: tb}
}

type suppressedHumaLogs struct {
	testing.TB
}

func (s *suppressedHumaLogs) Log(_ ...any) {
	// no-op
}
