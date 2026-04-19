package steamuser

import (
	"context"

	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/uptrace/bun"
)

type DAL interface {
	Get(ctx context.Context, steamID uint64) (*models.SteamUser, error)
	Insert(ctx context.Context, user *models.SteamUser) error
	Update(ctx context.Context, user *models.SteamUser) error
}

type dal struct {
	db *bun.DB
}

func New(db *bun.DB) DAL {
	return &dal{db: db}
}

func (d *dal) Get(ctx context.Context, steamID uint64) (*models.SteamUser, error) {
	user := &models.SteamUser{SteamID: steamID}
	err := d.db.NewSelect().Model(user).WherePK().Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d *dal) Insert(ctx context.Context, user *models.SteamUser) error {
	_, err := d.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (d *dal) Update(ctx context.Context, user *models.SteamUser) error {
	_, err := d.db.NewUpdate().Model(user).WherePK().Exec(ctx)
	return err
}
