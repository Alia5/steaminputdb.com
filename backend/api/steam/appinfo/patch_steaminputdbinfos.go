package appinfo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Alia5/steaminputdb.com/api/ctx"
	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/steam/auth"
	"github.com/Alia5/steaminputdb.com/db"
	appinfodal "github.com/Alia5/steaminputdb.com/db/dal/appinfo"
	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/steam/client"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type PatchSteamInputDBInfosRequest struct {
	AppID uint32 `path:"app_id"`
	Body  SIDBControllerSupport
}

type UpdateSteamInputDBInfosBody struct {
	// Define the fields that can be updated here
}

func registerPatchSteamInputDBInfos(a huma.API, dal db.DAL, registry huma.Registry, sc client.Client, useMemCache bool, cache *memcache.Cache) {
	huma.Register(a, huma.Operation{
		Method:  "PATCH",
		Path:    "/v1/steam/appinfo/{app_id}/steaminputdbinfos",
		Summary: "Update Steam Input DB Controller support info",
		Middlewares: huma.Middlewares{
			auth.ExtractSteamIDMiddleware,
		},
	}, func(c context.Context, req *PatchSteamInputDBInfosRequest) (*AppInfoResponse, error) {
		steamID, ok := c.Value(ctx.KeySteamID).(string)
		if !ok || steamID == "" {
			return nil, huma.Error403Forbidden("authentication error")
		}
		steamID64, err := strconv.ParseUint(steamID, 10, 64)
		if err != nil {
			return nil, huma.Error403Forbidden("invalid steam ID")
		}

		steamUser, err := dal.SteamUser().Get(c, steamID64)
		if steamUser == nil || err != nil {
			return nil, huma.Error403Forbidden("authentication error")
		}
		isAllowedInfoUpdate := steamUser.IsAdmin || steamUser.TrustedMod
		if !isAllowedInfoUpdate {
			return nil, huma.Error403Forbidden("authentication error")
		}

		// ---

		dbInfo, err := dal.AppInfo().Get(c, req.AppID, appinfodal.AppInfoInclude{
			SteamInputDBInfo: true,
		})
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, huma.Error404NotFound("item not found")
			}
			return nil, err
		}

		update := mapSteamInputDBInfoToAppInfo(req.AppID, &req.Body, dbInfo)
		err = dal.AppInfo().UpdateSteamInputDBInfo(c, update)
		if err != nil {
			slog.Error("db update failed", "error", err)
			return nil, huma.Error502BadGateway("database error", err)
		}

		updated, err := dal.AppInfo().Get(c, req.AppID, appinfodal.AppInfoInclude{
			ControllerSupport: true,
			Assets:            true,
			Links:             true,
			Creators:          true,
			OfficialConfigs:   true,
			SteamInputDBInfo:  true,
		})
		if err != nil {
			return nil, huma.Error502BadGateway("database error", err)
		}
		picsInfo, err := fetchFromSteamClient(c, sc, req.AppID)
		if err != nil {
			slog.Error("steam client fetch failed", "error", err)
		} else {
			enrichModelFromPICS(updated, picsInfo)
		}
		appInfo := mapModelToResponse(updated)
		if useMemCache {
			cache.Store(fmt.Sprint(req.AppID), appInfo)
		}
		return &AppInfoResponse{
			Body: appInfo,
		}, nil
	})
}

func mapSteamInputDBInfoToAppInfo(appID uint32, info *SIDBControllerSupport, dbInfo *models.AppInfo) *models.AppInfo {
	appInfo := &models.AppInfo{
		AppID:                   appID,
		ControllerSupportRating: info.ControllerSupportRating,
		ControllerSupportNotes:  info.ControllerSupportNotes,
	}

	mixedInput := info.MixedInputInfo
	if mixedInput != nil {
		appInfo.MixedInputInfo = &models.AppMixedInputInfo{
			MixedInputSupport:      mixedInput.MixedInputType,
			MixedInputGlyphFlicker: mixedInput.GlyphFlicker,
			MixedInputNotes:        mixedInput.Notes,
		}
		if mixedInput.MixedInputModURLS != nil {
			appInfo.MixedInputInfo.MixedInputModLinks = make([]*models.MixedInputModLinks, 0, len(mixedInput.MixedInputModURLS))
			for _, url := range mixedInput.MixedInputModURLS {
				url = strings.TrimSpace(url)
				if url == "" {
					continue
				}
				appInfo.MixedInputInfo.MixedInputModLinks = append(appInfo.MixedInputInfo.MixedInputModLinks, &models.MixedInputModLinks{
					Mod: url,
				})
			}
		}
	}

	glyphs := info.GlyphInfo
	if glyphs != nil {
		appInfo.Glyphs = &models.AppGlyphs{
			AutoGlyphDetect:   glyphs.AutoDetect,
			ManualGlyphSelect: glyphs.ManualSelect,
			GlyphNotes:        glyphs.Notes,
		}
		if glyphs.Controllers != nil {
			appInfo.Glyphs.GlyphCtrlSupport = make([]*models.AppGlyphCtrlSupport, 0, len(glyphs.Controllers))
			addedControllerTypes := make(map[steamtypes.ControllerType]bool, len(glyphs.Controllers))
			for _, ctrl := range glyphs.Controllers {
				if ctrl.ControllerType == nil {
					continue
				}
				if addedControllerTypes[*ctrl.ControllerType] {
					continue
				}
				addedControllerTypes[*ctrl.ControllerType] = true
				appInfo.Glyphs.GlyphCtrlSupport = append(appInfo.Glyphs.GlyphCtrlSupport, &models.AppGlyphCtrlSupport{
					ControllerType: *ctrl.ControllerType,
					Notes:          ctrl.Notes,
				})
			}
		}
	}

	siapi := info.SteamInputAPISupport
	if siapi != nil {
		appInfo.SteamInputAPISupport = &models.AppSteamInputAPISupport{
			SIAPICameraSupport: siapi.CameraSupport,
			SIAPIPixelsPer360:  siapi.PixelsPer360,
			Notes:              siapi.Notes,
		}
		if dbInfo.SteamInputAPISupport != nil {
			appInfo.SteamInputAPISupport.SIAPICameraNotes = dbInfo.SteamInputAPISupport.SIAPICameraNotes
		}
		if siapi.SupportTags != nil {
			appInfo.SteamInputAPISupport.SIAPITypes = make([]*models.AppSIAPITypes, len(siapi.SupportTags))
			for i, tag := range siapi.SupportTags {
				appInfo.SteamInputAPISupport.SIAPITypes[i] = &models.AppSIAPITypes{
					Type: tag,
				}
			}
		}
		if siapi.Glyphs != nil {
			if appInfo.Glyphs == nil {
				appInfo.Glyphs = &models.AppGlyphs{}
				if dbInfo.Glyphs != nil {
					appInfo.Glyphs.AutoGlyphDetect = dbInfo.Glyphs.AutoGlyphDetect
					appInfo.Glyphs.ManualGlyphSelect = dbInfo.Glyphs.ManualGlyphSelect
					appInfo.Glyphs.GlyphNotes = dbInfo.Glyphs.GlyphNotes
				}
			}
			appInfo.Glyphs.GlyphTags = make([]*models.AppGlyphTag, len(siapi.Glyphs))
			for i, tag := range siapi.Glyphs {
				appInfo.Glyphs.GlyphTags[i] = &models.AppGlyphTag{
					Tag: tag,
				}
			}
		}
	}

	if info.HWFeatures != nil {
		appInfo.HWFeatures = make([]*models.AppHWFeatures, 0, len(info.HWFeatures))
		addedHWFeatures := make(map[[2]int]bool, len(info.HWFeatures))
		for _, feature := range info.HWFeatures {
			key := [2]int{
				int(feature.Feature),
				int(feature.ControllerFamily),
			}
			if addedHWFeatures[key] {
				continue
			}
			addedHWFeatures[key] = true
			appInfo.HWFeatures = append(appInfo.HWFeatures, &models.AppHWFeatures{
				HWFeature:               feature.Feature,
				HWFeatureControllerType: feature.ControllerFamily,
				Notes:                   feature.Notes,
			})
		}
	}

	return appInfo
}
