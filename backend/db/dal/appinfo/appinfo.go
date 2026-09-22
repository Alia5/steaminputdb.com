package appinfo

import (
	"context"

	"github.com/Alia5/steaminputdb.com/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DAL interface {
	Get(ctx context.Context, appID uint32, include AppInfoInclude) (*models.AppInfo, error)
	Insert(ctx context.Context, appInfo *models.AppInfo) error
	UpdateBaseInfo(ctx context.Context, appInfo *models.AppInfo) error
	UpdateSteamInputDBInfo(ctx context.Context, appInfo *models.AppInfo) error
	UpdateControllerSupportRating(ctx context.Context, appID uint32, rating models.ControllerSupportRating, tx *gorm.DB) error
	UpdateControllerSupportNotes(ctx context.Context, appID uint32, notes *string, tx *gorm.DB) error
	UpdateMixedInputInfo(ctx context.Context, appID uint32, info *models.AppMixedInputInfo, tx *gorm.DB) error
	UpdateGlyphs(ctx context.Context, appID uint32, glyphs *models.AppGlyphs, tx *gorm.DB) error
	UpdateSteamInputAPISupport(ctx context.Context, appID uint32, support *models.AppSteamInputAPISupport, tx *gorm.DB) error
	UpdateHWFeatures(ctx context.Context, appID uint32, features []*models.AppHWFeatures, tx *gorm.DB) error
	UpdateHWFeatureNotes(ctx context.Context, appID uint32, notes *string, tx *gorm.DB) error
}

type dal struct {
	db *gorm.DB
}

func New(db *gorm.DB) DAL {
	return &dal{db: db}
}

type AppInfoInclude struct {
	ControllerSupport bool
	Assets            bool
	Links             bool
	Creators          bool
	OfficialConfigs   bool
	SteamInputDBInfo  bool
}

func (d *dal) Get(ctx context.Context, appID uint32, include AppInfoInclude) (*models.AppInfo, error) {
	q := d.db.WithContext(ctx)

	if include.ControllerSupport {
		q = q.Joins("ControllerSupport")
	}
	if include.Assets {
		q = q.Joins("Assets")
	}
	if include.Links {
		q = q.Preload("Links")
	}
	if include.Creators {
		q = q.Preload("CreatorLinks", func(links *gorm.DB) *gorm.DB {
			return links.Joins("AppCreator")
		})
	}
	if include.OfficialConfigs {
		q = q.Preload("OfficialConfigs")
	}
	if include.SteamInputDBInfo {
		q = q.Joins("MixedInputInfo").
			Preload("MixedInputInfo.MixedInputModLinks").
			Joins("SteamInputAPISupport").
			Preload("SteamInputAPISupport.SIAPITypes").
			Preload("HWFeatures").
			Joins("Glyphs").
			Preload("Glyphs.GlyphCtrlSupport").
			Preload("Glyphs.GlyphTags")
	}

	appInfo := &models.AppInfo{}
	res := q.Limit(1).Find(appInfo, "app_infos.app_id = ?", appID)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return appInfo, nil
}

func (d *dal) Insert(ctx context.Context, appInfo *models.AppInfo) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Omit(
			clause.Associations,
			"controller_support_rating", "controller_support_notes", "hw_feature_notes",
		).Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "app_id",
				},
			},
			UpdateAll: true,
		}).Create(appInfo).Error
		if err != nil {
			return err
		}
		return writebaseRelations(tx, appInfo)
	})
}

func (d *dal) UpdateBaseInfo(ctx context.Context, appInfo *models.AppInfo) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Model(appInfo).Select("*").Omit(
			clause.Associations,
			"created_at", "controller_support_rating", "controller_support_notes", "hw_feature_notes",
		).Updates(appInfo).Error
		if err != nil {
			return err
		}
		return writebaseRelations(tx, appInfo)
	})
}

func writebaseRelations(tx *gorm.DB, appInfo *models.AppInfo) error {
	if appInfo.ControllerSupport != nil {
		appInfo.ControllerSupport.AppID = appInfo.AppID
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "app_id",
				},
			},
			UpdateAll: true,
		}).Create(appInfo.ControllerSupport).Error
		if err != nil {
			return err
		}
	}

	if appInfo.Assets != nil {
		appInfo.Assets.AppID = appInfo.AppID
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "app_id",
				},
			},
			UpdateAll: true,
		}).Create(appInfo.Assets).Error
		if err != nil {
			return err
		}
	}

	if appInfo.Links != nil {
		err := tx.Where("app_id = ?", appInfo.AppID).Delete(&models.AppLink{}).Error
		if err != nil {
			return err
		}
		if len(appInfo.Links) > 0 {
			for _, link := range appInfo.Links {
				link.AppID = appInfo.AppID
			}
			err = tx.Create(&appInfo.Links).Error
			if err != nil {
				return err
			}
		}
	}

	if appInfo.OfficialConfigs != nil {
		err := tx.Where("app_id = ?", appInfo.AppID).Delete(&models.OfficialSteamInputConfig{}).Error
		if err != nil {
			return err
		}
		if len(appInfo.OfficialConfigs) > 0 {
			for _, cfg := range appInfo.OfficialConfigs {
				cfg.AppID = appInfo.AppID
			}
			err = tx.Create(&appInfo.OfficialConfigs).Error
			if err != nil {
				return err
			}
		}
	}

	if appInfo.CreatorLinks != nil {
		err := tx.Where("app_id = ?", appInfo.AppID).Delete(&models.AppCreatorToApp{}).Error
		if err != nil {
			return err
		}
		if len(appInfo.CreatorLinks) > 0 {
			for _, link := range appInfo.CreatorLinks {
				link.AppID = appInfo.AppID
				if link.AppCreator != nil {
					err = findOrCreateCreator(tx, link.AppCreator)
					if err != nil {
						return err
					}
					link.AppCreatorID = link.AppCreator.ID
				}
			}
			err = tx.Omit(clause.Associations).Create(&appInfo.CreatorLinks).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func findOrCreateCreator(tx *gorm.DB, creator *models.AppCreator) error {
	onConflict := clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "name",
			},
			{
				Name: "creator_clan_account_id",
			},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"name",
		}),
	}
	returning := clause.Returning{
		Columns: []clause.Column{
			{
				Name: "id",
			},
		},
	}
	return tx.Omit(clause.Associations).Clauses(onConflict, returning).Create(creator).Error
}
