package db

import (
	"log/slog"
	"strings"

	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/logging"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Init(cfg config.DB) (DAL, error) {
	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	isPostgres := strings.HasPrefix(cfg.DatabaseURL, "postgres://") || strings.HasPrefix(cfg.DatabaseURL, "postgresql://")

	var dialector gorm.Dialector
	if isPostgres {
		dialector = postgres.Open(cfg.DatabaseURL)
	} else {
		dialector = sqlite.Open(cfg.DatabaseURL)
	}

	db, err := gorm.Open(dialector, &gorm.Config{
		Logger:                 logging.NewGormLogger(cfg.SlowThreshold),
		TranslateError:         true,
		SkipDefaultTransaction: cfg.SkipDefaultTransaction,
		PrepareStmt:            cfg.PrepareStatement,
	})
	if err != nil {
		slog.Error("Could not open database", "error", err)
		return nil, err
	}

	sqldb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	if !isPostgres {
		sqldb.SetMaxOpenConns(1)
	}
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.MaxConnLifetime)
	sqldb.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	err = migrate(db)
	if err != nil {
		slog.Error("Could not migrate database", "error", err)
		return nil, err
	}

	return newDAL(db), nil
}
