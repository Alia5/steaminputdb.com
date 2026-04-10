package db

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Alia5/steaminputdb.com/buddy-app/db/migrations"
	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	"github.com/Alia5/steaminputdb.com/logging"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func Init() (*DAL, bool, error) {
	dir := filepath.Dir(install.DefaultInstallPath())
	if err := os.MkdirAll(dir, 0755); err != nil {
		slog.Error("failed to create database directory", "path", dir, "error", err)
		return nil, false, err
	}
	databaseURL := fmt.Sprintf(
		"file:%s/buddy-app.db?cache=shared&mode=rwc&_journal_mode=WAL&_busy_timeout=5000",
		dir,
	)
	slog.Debug("database URL", "dsn", databaseURL)

	sqldb, err := sql.Open(sqliteshim.ShimName, databaseURL)
	if err != nil {
		slog.Error("failed to open SQLite database", "dsn", databaseURL, "error", err)
		return nil, false, err
	}
	sqldb.SetMaxOpenConns(1)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)
	db := bun.NewDB(sqldb, sqlitedialect.New())

	db = db.WithQueryHook(logging.NewQueryHook())

	hasMigrated, err := migrations.Migrate(context.Background(), db)
	if err != nil {
		return nil, hasMigrated, err
	}

	return newDal(db), hasMigrated, nil
}
