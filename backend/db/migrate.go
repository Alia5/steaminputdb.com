package db

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/Alia5/steaminputdb.com/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var tables = []any{
	&models.AppCreatorRole{},
	&models.AppInfo{},
	&models.AppControllerSupport{},
	&models.AppAsset{},
	&models.AppLink{},
	&models.AppCreator{},
	&models.AppCreatorToApp{},
	&models.OfficialSteamInputConfig{},
	&models.SteamUser{},
}

func migrate(db *gorm.DB) error {
	if db.Migrator().HasTable("bun_migrations") {
		return db.Transaction(rebuildFromBun)
	}

	return autoMigrate(db)
}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	roles := []models.AppCreatorRole{
		{
			RoleID: models.AppCreatorRoleIDPublisher,
			Name:   "publisher",
		},
		{
			RoleID: models.AppCreatorRoleIDDeveloper,
			Name:   "developer",
		},
		{
			RoleID: models.AppCreatorRoleIDFranchise,
			Name:   "franchise",
		},
	}
	onConflict := clause.OnConflict{
		DoNothing: true,
	}
	return db.Clauses(onConflict).Create(&roles).Error
}

func tableName(db *gorm.DB, model any) (string, error) {
	stmt := &gorm.Statement{
		DB: db,
	}
	err := stmt.Parse(model)
	if err != nil {
		return "", err
	}
	return stmt.Schema.Table, nil
}

func rebuildFromBun(tx *gorm.DB) error {
	const legacy = "bun_legacy"
	slog.Info("rebuilding bun tables with gorm")

	names := make([]string, 0, len(tables))
	for _, model := range tables {
		name, err := tableName(tx, model)
		if err != nil {
			return err
		}
		names = append(names, name)
	}

	err := tx.Exec("CREATE SCHEMA " + legacy).Error
	if err != nil {
		return err
	}
	for _, name := range append([]string{"bun_migrations", "bun_migration_locks"}, names...) {
		err = tx.Exec(fmt.Sprintf("ALTER TABLE IF EXISTS %s SET SCHEMA %s", name, legacy)).Error
		if err != nil {
			return err
		}
	}

	err = autoMigrate(tx)
	if err != nil {
		return err
	}
	err = tx.Exec("DELETE FROM app_creator_roles").Error
	if err != nil {
		return err
	}

	for _, name := range names {
		var columns []string
		err = tx.Raw(
			"SELECT quote_ident(l.column_name) FROM information_schema.columns l"+
				" JOIN information_schema.columns n ON n.table_schema = current_schema() AND n.table_name = l.table_name AND n.column_name = l.column_name"+
				" WHERE l.table_schema = ? AND l.table_name = ? ORDER BY l.ordinal_position",
			legacy, name,
		).Scan(&columns).Error
		if err != nil {
			return err
		}
		cols := strings.Join(columns, ", ")

		res := tx.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s.%s", name, cols, cols, legacy, name))
		if res.Error != nil {
			return fmt.Errorf("copying %s: %w", name, res.Error)
		}

		var legacyRows int64
		err = tx.Raw(fmt.Sprintf("SELECT count(*) FROM %s.%s", legacy, name)).Scan(&legacyRows).Error
		if err != nil {
			return err
		}
		if res.RowsAffected != legacyRows {
			return fmt.Errorf("copying %s: copied %d of %d rows", name, res.RowsAffected, legacyRows)
		}
		slog.Info("copied table", "table", name, "rows", legacyRows)
	}

	return tx.Exec("DROP SCHEMA " + legacy + " CASCADE").Error
}
