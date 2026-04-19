package db

import (
	"context"
	"database/sql"
	"log/slog"
	"strings"

	"github.com/Alia5/steaminputdb.com/config"
	"github.com/Alia5/steaminputdb.com/db/migrations"
	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/logging"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/driver/sqliteshim"
)

func Init(cfg config.DB) (DAL, error) {

	ctx := context.Background()

	if cfg.DatabaseURL == "" {
		cfg.DatabaseURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	var db *bun.DB
	if strings.HasPrefix(cfg.DatabaseURL, "postgres://") || strings.HasPrefix(cfg.DatabaseURL, "postgresql://") {
		sqldb := sql.OpenDB(pgdriver.NewConnector(
			pgdriver.WithDSN(cfg.DatabaseURL),
		))

		sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
		sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
		sqldb.SetConnMaxLifetime(cfg.MaxConnLifetime)
		sqldb.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

		db = bun.NewDB(sqldb, pgdialect.New())
	} else {
		sqldb, err := sql.Open(sqliteshim.ShimName, cfg.DatabaseURL)
		if err != nil {
			slog.Error("Could not open SQLite database", "error", err)
			return nil, err
		}
		// SEE: https://bun.uptrace.dev/guide/drivers.html#important-in-memory-database-configuration
		sqldb.SetMaxOpenConns(1)
		sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
		sqldb.SetConnMaxLifetime(cfg.MaxConnLifetime)
		sqldb.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

		db = bun.NewDB(sqldb, sqlitedialect.New())
	}

	db = db.WithQueryHook(logging.NewQueryHook())

	db.RegisterModel((*models.AppCreatorToApp)(nil))

	err := migrations.Migrate(ctx, db)
	if err != nil {
		return nil, err
	}

	return newDAL(db), nil
}
