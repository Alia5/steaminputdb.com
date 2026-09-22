package appinfo

import (
	"context"

	"github.com/Alia5/steaminputdb.com/db/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *dal) UpdateSteamInputDBInfo(ctx context.Context, appInfo *models.AppInfo) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := d.UpdateControllerSupportRating(ctx, appInfo.AppID, appInfo.ControllerSupportRating, tx)
		if err != nil {
			return err
		}
		if appInfo.ControllerSupportNotes != nil {
			err = d.UpdateControllerSupportNotes(ctx, appInfo.AppID, appInfo.ControllerSupportNotes, tx)
			if err != nil {
				return err
			}
		}
		if appInfo.MixedInputInfo != nil {
			err = d.UpdateMixedInputInfo(ctx, appInfo.AppID, appInfo.MixedInputInfo, tx)
			if err != nil {
				return err
			}
		}
		if appInfo.Glyphs != nil {
			err = d.UpdateGlyphs(ctx, appInfo.AppID, appInfo.Glyphs, tx)
			if err != nil {
				return err
			}
		}
		if appInfo.SteamInputAPISupport != nil {
			err = d.UpdateSteamInputAPISupport(ctx, appInfo.AppID, appInfo.SteamInputAPISupport, tx)
			if err != nil {
				return err
			}
		}
		if appInfo.HWFeatures != nil {
			err = d.UpdateHWFeatures(ctx, appInfo.AppID, appInfo.HWFeatures, tx)
			if err != nil {
				return err
			}
		}
		if appInfo.HWFeatureNotes != nil {
			err = d.UpdateHWFeatureNotes(ctx, appInfo.AppID, appInfo.HWFeatureNotes, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *dal) UpdateControllerSupportRating(
	ctx context.Context,
	appID uint32,
	rating models.ControllerSupportRating,
	tx *gorm.DB,
) error {
	if tx == nil {
		tx = d.db
	}
	tx = tx.WithContext(ctx)

	res := tx.Model(&models.AppInfo{}).
		Where("app_id = ?", appID).
		UpdateColumn("controller_support_rating", rating)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *dal) UpdateControllerSupportNotes(
	ctx context.Context,
	appID uint32,
	notes *string,
	tx *gorm.DB,
) error {
	if tx == nil {
		tx = d.db
	}
	tx = tx.WithContext(ctx)

	res := tx.Model(&models.AppInfo{}).
		Where("app_id = ?", appID).
		UpdateColumn("controller_support_notes", notes)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *dal) UpdateHWFeatureNotes(
	ctx context.Context,
	appID uint32,
	notes *string,
	tx *gorm.DB,
) error {
	if tx == nil {
		tx = d.db
	}
	tx = tx.WithContext(ctx)

	res := tx.Model(&models.AppInfo{}).
		Where("app_id = ?", appID).
		UpdateColumn("hw_feature_notes", notes)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (d *dal) UpdateMixedInputInfo(
	ctx context.Context,
	appID uint32,
	info *models.AppMixedInputInfo,
	tx *gorm.DB,
) error {
	if tx == nil {
		return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return d.UpdateMixedInputInfo(ctx, appID, info, tx)
		})
	}
	tx = tx.WithContext(ctx)

	info.AppID = appID
	err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "app_id",
			},
		},
		UpdateAll: true,
	}).Create(info).Error
	if err != nil {
		return err
	}

	if info.MixedInputModLinks == nil {
		return nil
	}
	mods := make([]string, len(info.MixedInputModLinks))
	for i, link := range info.MixedInputModLinks {
		link.AppID = appID
		mods[i] = link.Mod
	}
	return syncRelations(
		tx, appID, info.MixedInputModLinks,
		clause.OnConflict{
			DoNothing: true,
		},
		"mod", mods,
	)
}

func (d *dal) UpdateGlyphs(
	ctx context.Context,
	appID uint32,
	glyphs *models.AppGlyphs,
	tx *gorm.DB,
) error {
	if tx == nil {
		return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return d.UpdateGlyphs(ctx, appID, glyphs, tx)
		})
	}
	tx = tx.WithContext(ctx)

	glyphs.AppID = appID
	err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "app_id",
			},
		},
		UpdateAll: true,
	}).Create(glyphs).Error
	if err != nil {
		return err
	}

	if glyphs.GlyphCtrlSupport != nil {
		ctrlTypes := make([]string, len(glyphs.GlyphCtrlSupport))
		for i, ctrl := range glyphs.GlyphCtrlSupport {
			ctrl.AppID = appID
			ctrlTypes[i] = string(ctrl.ControllerType)
		}
		err = syncRelations(
			tx, appID, glyphs.GlyphCtrlSupport,
			clause.OnConflict{
				Columns: []clause.Column{
					{
						Name: "app_id",
					},
					{
						Name: "controller_type",
					},
				},
				UpdateAll: true,
			},
			"controller_type", ctrlTypes,
		)
		if err != nil {
			return err
		}
	}

	if glyphs.GlyphTags != nil {
		tags := make([]models.AppGlyphTagType, len(glyphs.GlyphTags))
		for i, tag := range glyphs.GlyphTags {
			tag.AppID = appID
			tags[i] = tag.Tag
		}
		err = syncRelations(
			tx, appID, glyphs.GlyphTags,
			clause.OnConflict{
				DoNothing: true,
			},
			"tag", tags,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *dal) UpdateSteamInputAPISupport(
	ctx context.Context,
	appID uint32,
	support *models.AppSteamInputAPISupport,
	tx *gorm.DB,
) error {
	if tx == nil {
		return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return d.UpdateSteamInputAPISupport(ctx, appID, support, tx)
		})
	}
	tx = tx.WithContext(ctx)

	support.AppID = appID
	err := tx.Omit(clause.Associations).Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{
				Name: "app_id",
			},
		},
		UpdateAll: true,
	}).Create(support).Error
	if err != nil {
		return err
	}

	if support.SIAPITypes == nil {
		return nil
	}
	types := make([]models.SteamInputAPISupportType, len(support.SIAPITypes))
	for i, t := range support.SIAPITypes {
		t.AppID = appID
		types[i] = t.Type
	}
	return syncRelations(
		tx, appID, support.SIAPITypes,
		clause.OnConflict{
			DoNothing: true,
		},
		"type", types,
	)
}

func (d *dal) UpdateHWFeatures(
	ctx context.Context,
	appID uint32,
	features []*models.AppHWFeatures,
	tx *gorm.DB,
) error {
	if tx == nil {
		return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return d.UpdateHWFeatures(ctx, appID, features, tx)
		})
	}
	tx = tx.WithContext(ctx)

	keys := make([][]any, len(features))
	for i, feature := range features {
		feature.AppID = appID
		keys[i] = []any{
			feature.HWFeature,
			feature.HWFeatureControllerType,
		}
	}
	return syncRelations(
		tx, appID, features,
		clause.OnConflict{
			Columns: []clause.Column{
				{
					Name: "app_id",
				},
				{
					Name: "hw_feature",
				},
				{
					Name: "hw_feature_controller_type",
				},
			},
			UpdateAll: true,
		},
		"(hw_feature, hw_feature_controller_type)", keys,
	)
}

func syncRelations[T any](
	tx *gorm.DB,
	appID uint32,
	rows []*T,
	onConflict clause.OnConflict,
	keyExpr string,
	keys any,
) error {
	del := tx.Where("app_id = ?", appID)
	if len(rows) > 0 {
		err := tx.Clauses(onConflict).Create(&rows).Error
		if err != nil {
			return err
		}
		del = del.Where(keyExpr+" NOT IN ?", keys)
	}
	return del.Delete(new(T)).Error
}
