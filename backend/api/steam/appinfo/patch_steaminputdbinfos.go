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
	Body  UpdateSteamInputDBInfosBody
}

type UpdateSteamInputDBInfosBody struct {
	ControllerSupportRating *models.ControllerSupportRating `json:"controller_support_rating,omitempty,omitzero"`
	ControllerSupportNotes  *string                         `json:"controller_support_notes,omitempty,omitzero"`
	MixedInputInfo          *MixedInputInfo                 `json:"mixed_input,omitempty,omitzero"`
	GlyphInfo               *GlyphInfo                      `json:"glyphs,omitempty,omitzero"`
	SteamInputAPISupport    *SteamInputAPISupport           `json:"steaminputapi_support,omitempty,omitzero"`
	HWFeatures              []HWFeature                     `json:"hw_features,omitempty,omitzero"`
	HWFeatureNotes          *string                         `json:"hw_feature_notes,omitempty,omitzero"`
}

func registerPatchSteamInputDBInfos(a huma.API, dal db.DAL, _ huma.Registry, sc client.Client, useMemCache bool, cache *memcache.Cache) {
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

func mapSteamInputDBInfoToAppInfo(appID uint32, info *UpdateSteamInputDBInfosBody, dbInfo *models.AppInfo) *models.AppInfo {
	appInfo := &models.AppInfo{
		AppID:                   appID,
		ControllerSupportRating: dbInfo.ControllerSupportRating,
		ControllerSupportNotes:  info.ControllerSupportNotes,
		HWFeatureNotes:          info.HWFeatureNotes,
	}
	if info.ControllerSupportRating != nil {
		appInfo.ControllerSupportRating = *info.ControllerSupportRating
	}

	mixedInput := info.MixedInputInfo
	if mixedInput != nil {
		createsEmptyMixedInput := dbInfo.MixedInputInfo == nil
		createsEmptyMixedInput = createsEmptyMixedInput && mixedInput.MixedInputType == models.MixedInputSupportUnknown
		createsEmptyMixedInput = createsEmptyMixedInput && mixedInput.GlyphFlicker == nil
		createsEmptyMixedInput = createsEmptyMixedInput && mixedInput.Notes == ""
		createsEmptyMixedInput = createsEmptyMixedInput && len(mixedInput.MixedInputModURLs) == 0
		if createsEmptyMixedInput {
			mixedInput = nil
		}
	}
	if mixedInput != nil {
		appInfo.MixedInputInfo = &models.AppMixedInputInfo{
			MixedInputSupport:      mixedInput.MixedInputType,
			MixedInputGlyphFlicker: mixedInput.GlyphFlicker,
			MixedInputNotes:        mixedInput.Notes,
		}
		if mixedInput.MixedInputModURLs != nil {
			appInfo.MixedInputInfo.MixedInputModLinks = make(
				[]*models.MixedInputModLinks,
				0,
				len(mixedInput.MixedInputModURLs),
			)
			for _, mod := range mixedInput.MixedInputModURLs {
				uri := strings.TrimSpace(mod.URL)
				if uri == "" {
					continue
				}
				appInfo.MixedInputInfo.MixedInputModLinks = append(
					appInfo.MixedInputInfo.MixedInputModLinks,
					&models.MixedInputModLinks{
						URI:  uri,
						Name: strings.TrimSpace(mod.Name),
					},
				)
			}
		}
	}

	glyphs := info.GlyphInfo
	if glyphs != nil {
		createsEmptyGlyphs := dbInfo.Glyphs == nil
		createsEmptyGlyphs = createsEmptyGlyphs && glyphs.AutoDetect == nil
		createsEmptyGlyphs = createsEmptyGlyphs && glyphs.ManualSelect == nil
		createsEmptyGlyphs = createsEmptyGlyphs && glyphs.Notes == ""
		createsEmptyGlyphs = createsEmptyGlyphs && len(glyphs.Controllers) == 0
		if createsEmptyGlyphs {
			glyphs = nil
		}
	}
	if glyphs != nil {
		appInfo.Glyphs = &models.AppGlyphs{
			AutoGlyphDetect:   glyphs.AutoDetect,
			ManualGlyphSelect: glyphs.ManualSelect,
			GlyphNotes:        glyphs.Notes,
		}
		if glyphs.Controllers != nil {
			appInfo.Glyphs.GlyphCtrlSupport = make(
				[]*models.AppGlyphCtrlSupport,
				0,
				len(glyphs.Controllers),
			)
			addedControllerTypes := make(map[steamtypes.ControllerType]bool, len(glyphs.Controllers))
			for _, ctrl := range glyphs.Controllers {
				if ctrl.ControllerType == nil {
					continue
				}
				if addedControllerTypes[*ctrl.ControllerType] {
					continue
				}
				addedControllerTypes[*ctrl.ControllerType] = true
				appInfo.Glyphs.GlyphCtrlSupport = append(
					appInfo.Glyphs.GlyphCtrlSupport,
					&models.AppGlyphCtrlSupport{
						ControllerType: *ctrl.ControllerType,
						Notes:          ctrl.Notes,
					},
				)
			}
		}
	}

	siapi := info.SteamInputAPISupport
	if siapi != nil {
		if siapi.PixelsPer360 != nil {
			if *siapi.PixelsPer360 == "" {
				siapi.PixelsPer360 = nil
			}
		}
		createsEmptySIAPISupport := dbInfo.SteamInputAPISupport == nil
		createsEmptySIAPISupport = createsEmptySIAPISupport && dbInfo.Glyphs == nil
		createsEmptySIAPISupport = createsEmptySIAPISupport && siapi.CameraSupport == models.SteamInputCameraSupportUnknown
		createsEmptySIAPISupport = createsEmptySIAPISupport && siapi.PixelsPer360 == nil
		createsEmptySIAPISupport = createsEmptySIAPISupport && siapi.SteamInputType == nil
		createsEmptySIAPISupport = createsEmptySIAPISupport && siapi.Notes == ""
		createsEmptySIAPISupport = createsEmptySIAPISupport && len(siapi.SupportTags) == 0
		createsEmptySIAPISupport = createsEmptySIAPISupport && len(siapi.Glyphs) == 0
		if createsEmptySIAPISupport {
			siapi = nil
		}
	}
	if siapi != nil {
		appInfo.SteamInputAPISupport = &models.AppSteamInputAPISupport{
			SIAPICameraSupport: siapi.CameraSupport,
			SIAPIPixelsPer360:  siapi.PixelsPer360,
			SIAPIType:          siapi.SteamInputType,
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
