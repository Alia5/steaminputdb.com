package settings

import (
	"context"
	"time"

	"github.com/Alia5/steaminputdb.com/buddy-app/db/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DAL interface {
	Get(ctx context.Context) (*models.Settings, error)
	ResetToDefault(ctx context.Context) (*models.Settings, error)
	Update(ctx context.Context, settings *models.Settings) (*models.Settings, error)
}
type dal struct {
	db *gorm.DB
}

func New(db *gorm.DB) DAL {
	return &dal{db: db}
}

func (d *dal) Get(ctx context.Context) (*models.Settings, error) {
	settings := &models.Settings{}
	res := d.db.WithContext(ctx).Limit(1).Find(settings)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return settings, nil
}

func (d *dal) ResetToDefault(ctx context.Context) (*models.Settings, error) {
	settings, err := d.Get(ctx)
	if err != nil {
		return nil, err
	}
	settings.AddDesktopUIEntries = new(true)
	settings.AddBigPictureUIEntries = new(true)
	settings.DesktopUseSteamBrowser = new(false)
	settings.SteamWaitTimeout = new(60 * time.Second)

	err = d.db.WithContext(ctx).Model(settings).Updates(settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}

func (d *dal) Update(ctx context.Context, settings *models.Settings) (*models.Settings, error) {
	if settings == nil {
		return nil, nil
	}
	if settings.ID == uuid.Nil {
		dbSettings, err := d.Get(ctx)
		if err != nil {
			return nil, err
		}
		settings.ID = dbSettings.ID
	}
	err := d.db.WithContext(ctx).Model(settings).Updates(settings).Error
	if err != nil {
		return nil, err
	}
	return settings, nil
}
