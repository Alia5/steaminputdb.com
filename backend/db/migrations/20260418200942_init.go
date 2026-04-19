package migrations

import (
	"context"
	"log/slog"

	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		slog.Info("running init migration")

		_, err := db.NewCreateTable().Model((*models.AppCreatorRole)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		creatorRoles := []models.AppCreatorRole{
			{RoleID: models.AppCreatorRoleIDPublisher, Name: "publisher"},
			{RoleID: models.AppCreatorRoleIDDeveloper, Name: "developer"},
			{RoleID: models.AppCreatorRoleIDFranchise, Name: "franchise"},
		}
		_, err = db.NewInsert().Model(&creatorRoles).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppInfo)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppControllerSupport)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppAsset)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppLink)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppCreator)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.AppCreatorToApp)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		_, err = db.NewCreateTable().Model((*models.OfficialSteamInputConfig)(nil)).Exec(ctx)
		if err != nil {
			return err
		}

		return nil
	}, nil)
}
