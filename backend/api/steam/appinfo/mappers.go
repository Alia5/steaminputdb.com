package appinfo

import (
	"time"

	"github.com/Alia5/steaminputdb.com/api/search/games"
	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/Alia5/steaminputdb.com/types"
)

func mapStoreItemToModel(item *steamapi.StoreItem) *models.AppInfo {
	appInfo := &models.AppInfo{
		AppID:        item.GetAppid(),
		Name:         item.GetName(),
		StoreURLPath: item.GetStoreUrlPath(),
	}
	if item.Type != nil {
		appInfo.Type = games.TypeToString(item.Type)
	}
	if item.BasicInfo != nil {
		desc := item.BasicInfo.GetShortDescription()
		if desc != "" {
			appInfo.ShortDescription = &desc
		}
	}
	if item.Platforms != nil {
		appInfo.Platforms = models.AppPlatforms{
			Windows:      item.Platforms.Windows,
			Mac:          item.Platforms.Mac,
			SteamOSLinux: item.Platforms.SteamosLinux,
		}
	}
	if item.Release != nil {
		if item.Release.SteamReleaseDate != nil {
			t := time.Unix(int64(*item.Release.SteamReleaseDate), 0)
			appInfo.Release.SteamReleaseDate = &t
		}
		if item.Release.OriginalReleaseDate != nil && *item.Release.OriginalReleaseDate != 0 {
			t := time.Unix(int64(*item.Release.OriginalReleaseDate), 0)
			appInfo.Release.OriginalReleaseDate = &t
		}
	}
	if item.Assets != nil {
		appInfo.Assets = &models.AppAsset{
			AppID:              item.GetAppid(),
			AssetURLFormat:     item.Assets.AssetUrlFormat,
			MainCapsule:        item.Assets.MainCapsule,
			SmallCapsule:       item.Assets.SmallCapsule,
			Header:             item.Assets.Header,
			PackageHeader:      item.Assets.PackageHeader,
			PageBackground:     item.Assets.PageBackground,
			HeroCapsule:        item.Assets.HeroCapsule,
			HeroCapsule2X:      item.Assets.HeroCapsule_2X,
			LibraryCapsule:     item.Assets.LibraryCapsule,
			LibraryCapsule2X:   item.Assets.LibraryCapsule_2X,
			LibraryHero:        item.Assets.LibraryHero,
			LibraryHero2X:      item.Assets.LibraryHero_2X,
			CommunityIcon:      item.Assets.CommunityIcon,
			ClanAvatar:         item.Assets.ClanAvatar,
			PageBackgroundPath: item.Assets.PageBackgroundPath,
			RawPageBackground:  item.Assets.RawPageBackground,
		}
	}
	appInfo.Links = make([]*models.AppLink, 0)
	if item.Links != nil {
		for _, link := range item.Links {
			if link == nil || link.Url == nil {
				continue
			}
			appInfo.Links = append(appInfo.Links, &models.AppLink{
				AppID: item.GetAppid(),
				URL:   *link.Url,
			})
		}
	}
	appInfo.CreatorLinks = make([]*models.AppCreatorToApp, 0)
	if item.BasicInfo != nil {
		for _, pub := range item.BasicInfo.Publishers {
			if pub == nil || pub.Name == nil {
				continue
			}
			appInfo.CreatorLinks = append(appInfo.CreatorLinks, &models.AppCreatorToApp{
				RoleID: models.AppCreatorRoleIDPublisher,
				AppCreator: &models.AppCreator{
					Name:                 *pub.Name,
					CreatorClanAccountID: pub.GetCreatorClanAccountId(),
				},
			})
		}
		for _, dev := range item.BasicInfo.Developers {
			if dev == nil || dev.Name == nil {
				continue
			}
			appInfo.CreatorLinks = append(appInfo.CreatorLinks, &models.AppCreatorToApp{
				RoleID: models.AppCreatorRoleIDDeveloper,
				AppCreator: &models.AppCreator{
					Name:                 *dev.Name,
					CreatorClanAccountID: dev.GetCreatorClanAccountId(),
				},
			})
		}
		for _, fr := range item.BasicInfo.Franchises {
			if fr == nil || fr.Name == nil {
				continue
			}
			appInfo.CreatorLinks = append(appInfo.CreatorLinks, &models.AppCreatorToApp{
				RoleID: models.AppCreatorRoleIDFranchise,
				AppCreator: &models.AppCreator{
					Name:                 *fr.Name,
					CreatorClanAccountID: fr.GetCreatorClanAccountId(),
				},
			})
		}
	}
	return appInfo
}

func mapSteamInputDBInfo(appInfo *models.AppInfo) *SIDBControllerSupport {
	siDBInfo := &SIDBControllerSupport{
		ControllerSupportRating: appInfo.ControllerSupportRating,
		ControllerSupportNotes:  appInfo.ControllerSupportNotes,
		HWFeatureNotes:          appInfo.HWFeatureNotes,
	}
	if appInfo.MixedInputInfo != nil {
		siDBInfo.MixedInputInfo = &MixedInputInfo{
			MixedInputType: appInfo.MixedInputInfo.MixedInputSupport,
			GlyphFlicker:   appInfo.MixedInputInfo.MixedInputGlyphFlicker,
			Notes:          appInfo.MixedInputInfo.MixedInputNotes,
		}
		mixedInputModLen := len(appInfo.MixedInputInfo.MixedInputModLinks)
		if mixedInputModLen > 0 {
			siDBInfo.MixedInputInfo.MixedInputModURLs = make([]MixedInputMods, mixedInputModLen)
			for i, modLink := range appInfo.MixedInputInfo.MixedInputModLinks {
				siDBInfo.MixedInputInfo.MixedInputModURLs[i] = MixedInputMods{
					URL:  modLink.URI,
					Name: modLink.Name,
				}
			}
		}
	}
	if appInfo.Glyphs != nil {
		siDBInfo.GlyphInfo = &GlyphInfo{
			AutoDetect:   appInfo.Glyphs.AutoGlyphDetect,
			ManualSelect: appInfo.Glyphs.ManualGlyphSelect,
			Notes:        appInfo.Glyphs.GlyphNotes,
		}
		glyphControllersLen := len(appInfo.Glyphs.GlyphCtrlSupport)
		if glyphControllersLen > 0 {
			siDBInfo.GlyphInfo.Controllers = make([]GlyphControllerSupport, glyphControllersLen)
			for i, ctrlSupport := range appInfo.Glyphs.GlyphCtrlSupport {
				siDBInfo.GlyphInfo.Controllers[i].ControllerType = &ctrlSupport.ControllerType
				siDBInfo.GlyphInfo.Controllers[i].Notes = ctrlSupport.Notes
			}
		}
	}
	if appInfo.SteamInputAPISupport != nil || appInfo.Glyphs != nil {
		var siapisupport *SteamInputAPISupport
		if appInfo.SteamInputAPISupport != nil {
			siapisupport = &SteamInputAPISupport{
				CameraSupport:  appInfo.SteamInputAPISupport.SIAPICameraSupport,
				PixelsPer360:   appInfo.SteamInputAPISupport.SIAPIPixelsPer360,
				SteamInputType: appInfo.SteamInputAPISupport.SIAPIType,
				Notes:          appInfo.SteamInputAPISupport.Notes,
			}
			lenSupportTags := len(appInfo.SteamInputAPISupport.SIAPITypes)
			if lenSupportTags > 0 {
				siapisupport.SupportTags = make([]models.SteamInputAPISupportType, lenSupportTags)
				for i, tag := range appInfo.SteamInputAPISupport.SIAPITypes {
					siapisupport.SupportTags[i] = tag.Type
				}
			}
		}
		if siapisupport == nil {
			siapisupport = &SteamInputAPISupport{}
		}
		if appInfo.Glyphs != nil {
			lenGlyphTags := len(appInfo.Glyphs.GlyphTags)
			if lenGlyphTags > 0 {
				siapisupport.Glyphs = make([]models.AppGlyphTagType, lenGlyphTags)
				for i, tag := range appInfo.Glyphs.GlyphTags {
					siapisupport.Glyphs[i] = tag.Tag
				}
			}
		}
		siDBInfo.SteamInputAPISupport = siapisupport
	}
	lenHWFeatures := len(appInfo.HWFeatures)
	if lenHWFeatures > 0 {
		siDBInfo.HWFeatures = make([]HWFeature, lenHWFeatures)
		for i, hwFeature := range appInfo.HWFeatures {
			siDBInfo.HWFeatures[i] = HWFeature{
				Feature:          hwFeature.HWFeature,
				ControllerFamily: hwFeature.HWFeatureControllerType,
				Notes:            hwFeature.Notes,
			}
		}
	}
	return siDBInfo
}

func mapModelToResponse(appInfo *models.AppInfo) *AppInfoItem {
	wrapper := &AppInfoItem{
		AppItem: games.AppItem{
			AppID:        &appInfo.AppID,
			Name:         &appInfo.Name,
			StoreURLPath: &appInfo.StoreURLPath,
			Type:         appInfo.Type,
			Platforms: games.AppsPlatforms{
				Windows:      appInfo.Platforms.Windows,
				SteamOSLinux: appInfo.Platforms.SteamOSLinux,
				Mac:          appInfo.Platforms.Mac,
			},
		},
	}
	wrapper.SteamInputDBInfo = mapSteamInputDBInfo(appInfo)

	if appInfo.Release.SteamReleaseDate != nil {
		wrapper.Release.SteamReleaseDate = *appInfo.Release.SteamReleaseDate
	}
	if appInfo.Release.OriginalReleaseDate != nil {
		wrapper.Release.OriginalReleaseDate = *appInfo.Release.OriginalReleaseDate
	}
	if appInfo.Assets != nil {
		wrapper.Assets = &steamapi.StoreItem_Assets{
			AssetUrlFormat:     appInfo.Assets.AssetURLFormat,
			MainCapsule:        appInfo.Assets.MainCapsule,
			SmallCapsule:       appInfo.Assets.SmallCapsule,
			Header:             appInfo.Assets.Header,
			PackageHeader:      appInfo.Assets.PackageHeader,
			PageBackground:     appInfo.Assets.PageBackground,
			HeroCapsule:        appInfo.Assets.HeroCapsule,
			HeroCapsule_2X:     appInfo.Assets.HeroCapsule2X,
			LibraryCapsule:     appInfo.Assets.LibraryCapsule,
			LibraryCapsule_2X:  appInfo.Assets.LibraryCapsule2X,
			LibraryHero:        appInfo.Assets.LibraryHero,
			LibraryHero_2X:     appInfo.Assets.LibraryHero2X,
			CommunityIcon:      appInfo.Assets.CommunityIcon,
			ClanAvatar:         appInfo.Assets.ClanAvatar,
			PageBackgroundPath: appInfo.Assets.PageBackgroundPath,
			RawPageBackground:  appInfo.Assets.RawPageBackground,
		}
	}
	if len(appInfo.Links) > 0 {
		links := make([]string, 0, len(appInfo.Links))
		for _, link := range appInfo.Links {
			links = append(links, link.URL)
		}
		wrapper.Links = &links
	}
	if appInfo.ShortDescription != nil || len(appInfo.CreatorLinks) > 0 {
		bi := &steamapi.StoreItem_BasicInfo{}
		if appInfo.ShortDescription != nil {
			bi.ShortDescription = appInfo.ShortDescription
		}
		for _, cl := range appInfo.CreatorLinks {
			creator := cl.AppCreator
			if creator == nil {
				continue
			}
			name := creator.Name
			link := &steamapi.StoreItem_BasicInfo_CreatorHomeLink{
				Name: &name,
			}
			if creator.CreatorClanAccountID != 0 {
				clanID := creator.CreatorClanAccountID
				link.CreatorClanAccountId = &clanID
			}
			switch cl.RoleID {
			case models.AppCreatorRoleIDPublisher:
				bi.Publishers = append(bi.Publishers, link)
			case models.AppCreatorRoleIDDeveloper:
				bi.Developers = append(bi.Developers, link)
			case models.AppCreatorRoleIDFranchise:
				bi.Franchises = append(bi.Franchises, link)
			}
		}
		wrapper.BasicInfo = bi
	}
	if appInfo.ControllerSupport != nil {
		wrapper.ControllerSupport = &ControllerSupport{
			SupportLevel:         appInfo.ControllerSupport.SupportLevel,
			DS4WiredSupport:      appInfo.ControllerSupport.DS4WiredSupport,
			DS4WirelessSupport:   appInfo.ControllerSupport.DS4WirelessSupport,
			DS5WiredSupport:      appInfo.ControllerSupport.DS5WiredSupport,
			DS5WirelessSupport:   appInfo.ControllerSupport.DS5WirelessSupport,
			SteamInputAPISupport: appInfo.ControllerSupport.SteamInputAPISupport,
		}
		if wrapper.SteamInputDBInfo != nil {
			if appInfo.ControllerSupport.SupportLevel != nil &&
				*appInfo.ControllerSupport.SupportLevel == types.ControllerSupportLevelNone {
				wrapper.SteamInputDBInfo.ControllerSupportRating = models.ControllerSupportRatingWood
			}
		}
	}
	if len(appInfo.OfficialConfigs) > 0 {
		oc := make(officialConfigs, len(appInfo.OfficialConfigs))
		for _, cfg := range appInfo.OfficialConfigs {
			oc[cfg.ControllerType] = configID(cfg.ConfigID)
		}
		wrapper.OfficialConfigs = &oc
	}

	return wrapper
}
