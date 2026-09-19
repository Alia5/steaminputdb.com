package steamuser

import (
	"context"

	"github.com/Alia5/steaminputdb.com/db/models"
	"gorm.io/gorm"
)

type DAL interface {
	Get(ctx context.Context, steamID uint64) (*models.SteamUser, error)
	Insert(ctx context.Context, user *models.SteamUser) error
	Update(ctx context.Context, user *models.SteamUser) error
}

type dal struct {
	db *gorm.DB
}

func New(db *gorm.DB) DAL {
	return &dal{db: db}
}

func (d *dal) Get(ctx context.Context, steamID uint64) (*models.SteamUser, error) {
	user := &models.SteamUser{}
	res := d.db.WithContext(ctx).Limit(1).Find(user, "steam_id = ?", steamID)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}

func (d *dal) Insert(ctx context.Context, user *models.SteamUser) error {
	return d.db.WithContext(ctx).Create(user).Error
}

func (d *dal) Update(ctx context.Context, user *models.SteamUser) error {
	return d.db.WithContext(ctx).Model(user).Select("*").Omit("created_at").Updates(user).Error
}
