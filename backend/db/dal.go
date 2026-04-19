package db

import (
	"github.com/Alia5/steaminputdb.com/db/dal/appinfo"
	"github.com/Alia5/steaminputdb.com/db/dal/steamuser"
	"github.com/uptrace/bun"
)

type DAL interface {
	AppInfo() appinfo.DAL
	SteamUser() steamuser.DAL
	DB() *bun.DB
}

type dal struct {
	db        *bun.DB
	appInfo   appinfo.DAL
	steamUser steamuser.DAL
}

func newDAL(db *bun.DB) DAL {
	return &dal{
		db:        db,
		appInfo:   appinfo.New(db),
		steamUser: steamuser.New(db),
	}
}

func (d *dal) AppInfo() appinfo.DAL {
	return d.appInfo
}

func (d *dal) SteamUser() steamuser.DAL {
	return d.steamUser
}

func (d *dal) DB() *bun.DB {
	return d.db
}
