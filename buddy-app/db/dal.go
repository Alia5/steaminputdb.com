package db

import (
	"github.com/Alia5/steaminputdb.com/buddy-app/db/dal/settings"
	"github.com/uptrace/bun"
)

type DAL struct {
	DB       *bun.DB
	Settings settings.DAL
}

func newDal(db *bun.DB) *DAL {
	return &DAL{
		DB:       db,
		Settings: settings.New(db),
	}
}
