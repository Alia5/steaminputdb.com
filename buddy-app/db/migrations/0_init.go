package migrations

import (
	"context"
	"log/slog"

	"github.com/Alia5/steaminputdb.com/buddy-app/db/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		slog.Info("running init migration")

		_, err := db.NewCreateTable().Model((*models.Settings)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewInsert().Model(&models.Settings{}).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
