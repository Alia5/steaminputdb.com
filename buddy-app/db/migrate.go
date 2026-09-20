package db

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/Alia5/steaminputdb.com/buddy-app/db/models"
	"gorm.io/gorm"
)

var tables = []any{
	&models.Settings{},
}

func migrate(db *gorm.DB) (bool, error) {
	if db.Migrator().HasTable("bun_migrations") {
		err := db.Transaction(rebuildFromBun)
		return false, err
	}

	isNew := !db.Migrator().HasTable(&models.Settings{})
	err := db.AutoMigrate(tables...)
	if err != nil {
		return false, err
	}
	err = ensureDefaults(db)
	if err != nil {
		return false, err
	}
	if isNew {
		slog.Info("created database")
	}
	return isNew, nil
}

func ensureDefaults(db *gorm.DB) error {
	var count int64
	err := db.Model(&models.Settings{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Create(&models.Settings{}).Error
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

func columns(tx *gorm.DB, table string) ([]string, error) {
	var names []string
	err := tx.Raw("SELECT name FROM pragma_table_info(?)", table).Scan(&names).Error
	return names, err
}

func rebuildFromBun(tx *gorm.DB) error {
	const suffix = "_bun"
	slog.Info("rebuilding bun tables with gorm")

	names := make([]string, 0, len(tables))
	for _, model := range tables {
		name, err := tableName(tx, model)
		if err != nil {
			return err
		}
		names = append(names, name)
	}

	for _, name := range names {
		err := tx.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", name, name+suffix)).Error
		if err != nil {
			return err
		}
	}

	err := tx.AutoMigrate(tables...)
	if err != nil {
		return err
	}

	for _, name := range names {
		legacyColumns, err := columns(tx, name+suffix)
		if err != nil {
			return err
		}
		newColumns, err := columns(tx, name)
		if err != nil {
			return err
		}
		shared := make([]string, 0, len(legacyColumns))
		for _, column := range legacyColumns {
			if slices.Contains(newColumns, column) {
				shared = append(shared, column)
			}
		}
		cols := strings.Join(shared, ", ")

		res := tx.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", name, cols, cols, name+suffix))
		if res.Error != nil {
			return fmt.Errorf("copying %s: %w", name, res.Error)
		}

		var legacyRows int64
		err = tx.Raw(fmt.Sprintf("SELECT count(*) FROM %s", name+suffix)).Scan(&legacyRows).Error
		if err != nil {
			return err
		}
		if res.RowsAffected != legacyRows {
			return fmt.Errorf("copying %s: copied %d of %d rows", name, res.RowsAffected, legacyRows)
		}
		slog.Info("copied table", "table", name, "rows", legacyRows)

		err = tx.Exec(fmt.Sprintf("DROP TABLE %s", name+suffix)).Error
		if err != nil {
			return err
		}
	}

	for _, name := range []string{
		"bun_migrations",
		"bun_migration_locks",
	} {
		err = tx.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name)).Error
		if err != nil {
			return err
		}
	}

	return ensureDefaults(tx)
}
