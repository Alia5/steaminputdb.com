package db

import (
	"github.com/Alia5/steaminputdb.com/buddy-app/db/dal/settings"
	"gorm.io/gorm"
)

type DAL struct {
	DB       *gorm.DB
	Settings settings.DAL
}

func newDal(db *gorm.DB) *DAL {
	return &DAL{
		DB:       db,
		Settings: settings.New(db),
	}
}
