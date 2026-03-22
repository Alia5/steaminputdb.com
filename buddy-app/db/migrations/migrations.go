package migrations

import (
	"context"
	"log/slog"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

// add back in case we ever use sql based migrations
// //go:embed *.sql
// var sqlMigrations  embed.FS

var Migrations = migrate.NewMigrations()

// func init() {
// 	if err := Migrations.Discover(sqlMigrations ); err != nil {
// 		panic(err)
// 	}
// }

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)
	err := migrator.Init(ctx)
	if err != nil {
		return err
	}
	if err := migrator.Lock(ctx); err != nil {
		return err
	}
	defer migrator.Unlock(ctx) //nolint:errcheck

	group, err := migrator.Migrate(ctx)
	if err != nil {
		return err
	}
	if group.IsZero() {
		slog.Debug("there are no new migrations to run (database is up to date)")
		return nil
	}
	slog.Info("migrated", "to", group)
	return nil
}
