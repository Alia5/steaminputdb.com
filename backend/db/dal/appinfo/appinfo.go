package appinfo

import (
	"context"

	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/uptrace/bun"
)

type DAL interface {
	Get(ctx context.Context, appID uint32, include AppInfoInclude) (*models.AppInfo, error)
	Insert(ctx context.Context, appInfo *models.AppInfo) error
	Update(ctx context.Context, appInfo *models.AppInfo) error
}

type dal struct {
	db *bun.DB
}

func New(db *bun.DB) DAL {
	return &dal{db: db}
}

type AppInfoInclude struct {
	ControllerSupport bool
	Assets            bool
	Links             bool
	Creators          bool
	OfficialConfigs   bool
}

func (d *dal) Get(ctx context.Context, appID uint32, include AppInfoInclude) (*models.AppInfo, error) {
	appInfo := &models.AppInfo{AppID: appID}
	q := d.db.NewSelect().Model(appInfo).WherePK()

	if include.ControllerSupport {
		q = q.Relation("ControllerSupport")
	}
	if include.Assets {
		q = q.Relation("Assets")
	}
	if include.Links {
		q = q.Relation("Links")
	}
	if include.Creators {
		q = q.Relation("CreatorLinks").Relation("CreatorLinks.AppCreator")
	}
	if include.OfficialConfigs {
		q = q.Relation("OfficialConfigs")
	}

	err := q.Scan(ctx)
	if err != nil {
		return nil, err
	}
	return appInfo, nil
}

func (d *dal) Insert(ctx context.Context, appInfo *models.AppInfo) error {
	return d.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewInsert().Model(appInfo).
			On("CONFLICT (app_id) DO UPDATE").
			Exec(ctx)
		if err != nil {
			return err
		}

		if appInfo.ControllerSupport != nil {
			appInfo.ControllerSupport.AppID = appInfo.AppID
			_, err = tx.NewInsert().Model(appInfo.ControllerSupport).
				On("CONFLICT (app_id) DO UPDATE").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		if appInfo.Assets != nil {
			appInfo.Assets.AppID = appInfo.AppID
			_, err = tx.NewInsert().Model(appInfo.Assets).
				On("CONFLICT (app_id) DO UPDATE").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(appInfo.Links) > 0 {
			_, err = tx.NewDelete().Model((*models.AppLink)(nil)).
				Where("app_id = ?", appInfo.AppID).Exec(ctx)
			if err != nil {
				return err
			}
			for _, link := range appInfo.Links {
				link.AppID = appInfo.AppID
			}
			_, err = tx.NewInsert().Model(&appInfo.Links).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(appInfo.OfficialConfigs) > 0 {
			_, err = tx.NewDelete().Model((*models.OfficialSteamInputConfig)(nil)).
				Where("app_id = ?", appInfo.AppID).Exec(ctx)
			if err != nil {
				return err
			}
			for _, cfg := range appInfo.OfficialConfigs {
				cfg.AppID = appInfo.AppID
			}
			_, err = tx.NewInsert().Model(&appInfo.OfficialConfigs).Exec(ctx)
			if err != nil {
				return err
			}
		}

		if len(appInfo.CreatorLinks) > 0 {
			_, err = tx.NewDelete().Model((*models.AppCreatorToApp)(nil)).
				Where("app_id = ?", appInfo.AppID).Exec(ctx)
			if err != nil {
				return err
			}
			for _, link := range appInfo.CreatorLinks {
				if link.AppCreator != nil {
					creator, err := findOrCreateCreator(ctx, tx, link.AppCreator)
					if err != nil {
						return err
					}
					link.AppCreator = creator
					link.AppCreatorID = creator.ID
				}
				link.AppID = appInfo.AppID
			}
			_, err = tx.NewInsert().Model(&appInfo.CreatorLinks).Exec(ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *dal) Update(ctx context.Context, appInfo *models.AppInfo) error {
	return d.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		_, err := tx.NewUpdate().Model(appInfo).WherePK().Exec(ctx)
		if err != nil {
			return err
		}

		if appInfo.ControllerSupport != nil {
			appInfo.ControllerSupport.AppID = appInfo.AppID
			_, err = tx.NewInsert().Model(appInfo.ControllerSupport).
				On("CONFLICT (app_id) DO UPDATE").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		if appInfo.Assets != nil {
			appInfo.Assets.AppID = appInfo.AppID
			_, err = tx.NewInsert().Model(appInfo.Assets).
				On("CONFLICT (app_id) DO UPDATE").
				Exec(ctx)
			if err != nil {
				return err
			}
		}

		if appInfo.Links != nil {
			_, err = tx.NewDelete().Model((*models.AppLink)(nil)).
				Where("app_id = ?", appInfo.AppID).
				Exec(ctx)
			if err != nil {
				return err
			}
			if len(appInfo.Links) > 0 {
				for _, link := range appInfo.Links {
					link.AppID = appInfo.AppID
				}
				_, err = tx.NewInsert().Model(&appInfo.Links).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}

		if appInfo.OfficialConfigs != nil {
			_, err = tx.NewDelete().Model((*models.OfficialSteamInputConfig)(nil)).
				Where("app_id = ?", appInfo.AppID).
				Exec(ctx)
			if err != nil {
				return err
			}
			if len(appInfo.OfficialConfigs) > 0 {
				for _, cfg := range appInfo.OfficialConfigs {
					cfg.AppID = appInfo.AppID
				}
				_, err = tx.NewInsert().Model(&appInfo.OfficialConfigs).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}

		if appInfo.CreatorLinks != nil {
			_, err = tx.NewDelete().Model((*models.AppCreatorToApp)(nil)).
				Where("app_id = ?", appInfo.AppID).
				Exec(ctx)
			if err != nil {
				return err
			}
			if len(appInfo.CreatorLinks) > 0 {
				for _, link := range appInfo.CreatorLinks {
					if link.AppCreator != nil {
						creator, err := findOrCreateCreator(ctx, tx, link.AppCreator)
						if err != nil {
							return err
						}
						link.AppCreator = creator
						link.AppCreatorID = creator.ID
					}
					link.AppID = appInfo.AppID
				}
				_, err = tx.NewInsert().Model(&appInfo.CreatorLinks).Exec(ctx)
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func findOrCreateCreator(ctx context.Context, tx bun.Tx, creator *models.AppCreator) (*models.AppCreator, error) {
	existing := &models.AppCreator{}
	q := tx.NewSelect().Model(existing).Where("name = ?", creator.Name)
	if creator.CreatorClanAccountID != nil {
		q = q.Where("creator_clan_account_id = ?", *creator.CreatorClanAccountID)
	} else {
		q = q.Where("creator_clan_account_id IS NULL")
	}
	err := q.Limit(1).Scan(ctx)
	if err == nil {
		return existing, nil
	}

	_, err = tx.NewInsert().Model(creator).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return creator, nil
}
