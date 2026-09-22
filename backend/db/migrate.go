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
	&models.AppMixedInputInfo{},
	&models.MixedInputModLinks{},
	&models.AppGlyphs{},
	&models.AppGlyphCtrlSupport{},
	&models.AppGlyphTag{},
	&models.AppSteamInputAPISupport{},
	&models.AppSIAPITypes{},
	&models.AppHWFeatures{},
}

func migrate(db *gorm.DB) error {
	err := renameMixedInputModToURI(db)
	if err != nil {
		return err
	}

	if db.Migrator().HasTable("bun_migrations") {
		return db.Transaction(rebuildFromBun)
	}

	return autoMigrate(db)
}

func renameMixedInputModToURI(db *gorm.DB) error {
	model := &models.MixedInputModLinks{}
	name, err := tableName(db, model)
	if err != nil {
		return err
	}

	m := db.Migrator()
	if !m.HasTable(name) || !m.HasColumn(model, "mod") {
		return nil
	}

	if !m.HasColumn(model, "uri") {
		return db.Transaction(func(tx *gorm.DB) error {
			return renameMixedInputModColumn(tx, model, name)
		})
	}

	return db.Transaction(func(tx *gorm.DB) error {
		return mergeMixedInputModColumn(tx, model, name)
	})
}

func renameMixedInputModColumn(tx *gorm.DB, model any, name string) error {
	m := tx.Migrator()
	err := m.RenameColumn(model, "mod", "uri")
	if err != nil {
		return err
	}
	slog.Info("renamed column", "table", name, "from", "mod", "to", "uri")

	oldIndex := tx.NamingStrategy.IndexName(name, "mod")
	newIndex := tx.NamingStrategy.IndexName(name, "uri")
	if m.HasIndex(model, oldIndex) && !m.HasIndex(model, newIndex) {
		return m.RenameIndex(model, oldIndex, newIndex)
	}
	return nil
}

func mergeMixedInputModColumn(tx *gorm.DB, model any, name string) error {
	res := tx.Exec(fmt.Sprintf("UPDATE %s SET uri = mod WHERE uri IS NULL AND mod IS NOT NULL", name))
	if res.Error != nil {
		return res.Error
	}
	slog.Info("backfilled column", "table", name, "from", "mod", "to", "uri", "rows", res.RowsAffected)

	err := tx.Exec(fmt.Sprintf("DELETE FROM %s WHERE uri IS NULL", name)).Error
	if err != nil {
		return err
	}

	dedupe := fmt.Sprintf(
		"DELETE FROM %s a USING %s b WHERE a.ctid < b.ctid AND a.app_id = b.app_id AND a.uri = b.uri",
		name, name,
	)
	if tx.Name() != "postgres" {
		dedupe = fmt.Sprintf(
			"DELETE FROM %s WHERE rowid NOT IN (SELECT MAX(rowid) FROM %s GROUP BY app_id, uri)",
			name, name,
		)
	}
	err = tx.Exec(dedupe).Error
	if err != nil {
		return err
	}

	err = tx.Migrator().DropColumn(model, "mod")
	if err != nil {
		return err
	}
	slog.Info("dropped column", "table", name, "column", "mod")

	if tx.Name() != "postgres" {
		return nil
	}
	err = tx.Exec(fmt.Sprintf("ALTER TABLE %s ALTER COLUMN uri SET NOT NULL", name)).Error
	if err != nil {
		return err
	}
	var hasPrimaryKey bool
	err = tx.Raw(
		"SELECT EXISTS (SELECT 1 FROM pg_index WHERE indrelid = ?::regclass AND indisprimary)",
		name,
	).Scan(&hasPrimaryKey).Error
	if err != nil {
		return err
	}
	if hasPrimaryKey {
		return nil
	}
	return tx.Exec(fmt.Sprintf("ALTER TABLE %s ADD PRIMARY KEY (app_id, uri)", name)).Error
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
