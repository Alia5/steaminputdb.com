package db

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/Alia5/steaminputdb.com/buddy-app/install"
	"github.com/Alia5/steaminputdb.com/logging"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init() (*DAL, bool, error) {
	dir := filepath.Dir(install.DefaultInstallPath())
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		slog.Error("failed to create database directory", "path", dir, "error", err)
		return nil, false, err
	}
	databaseURL := fmt.Sprintf(
		"file:%s/buddy-app.db?cache=shared&mode=rwc&_journal_mode=WAL&_busy_timeout=5000",
		dir,
	)
	return Open(databaseURL)
}

func Open(databaseURL string) (*DAL, bool, error) {
	slog.Debug("database URL", "dsn", databaseURL)

	db, err := gorm.Open(sqlite.Open(databaseURL), &gorm.Config{
		Logger:                 logging.NewGormLogger(200 * time.Millisecond),
		TranslateError:         true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		slog.Error("failed to open SQLite database", "dsn", databaseURL, "error", err)
		return nil, false, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, false, err
	}
	sqldb.SetMaxOpenConns(1)
	sqldb.SetMaxIdleConns(10)
	sqldb.SetConnMaxLifetime(5 * time.Minute)
	sqldb.SetConnMaxIdleTime(5 * time.Minute)

	hasMigrated, err := migrate(db)
	if err != nil {
		return nil, hasMigrated, err
	}

	return newDal(db), hasMigrated, nil
}
