package appinfo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/search/games"
	"github.com/Alia5/steaminputdb.com/db"
	appinfodal "github.com/Alia5/steaminputdb.com/db/dal/appinfo"
	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/steam/client"
	clientappinfo "github.com/Alia5/steaminputdb.com/steam/client/appinfo"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/Alia5/steaminputdb.com/types"
	"github.com/danielgtaylor/huma/v2"
)

const dbMaxAge = 24 * time.Hour

func RegisterRoute(a huma.API, dal db.DAL, opts ...bool) {
	var useMemCache bool
	if len(opts) > 0 {
		useMemCache = opts[0]
	} else {
		useMemCache = true
	}
	cache := memcache.New(30*time.Minute, 1000)

	sc := client.New()
	ctx := context.Background()
	if err := sc.Connect(ctx); err != nil {
		slog.Error("steam client connect failed", "error", err)
	} else if err := sc.Login(ctx, client.LoginDetails{Anonymous: true, Language: "english"}); err != nil {
		slog.Error("steam client login failed", "error", err)
	}

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/appinfo",
			Tags:        []string{"steam", "app"},
			Summary:     "Get Steam app info",
			Description: "Retrieve app information from Steam Store for a given app ID",
			Errors: []int{
				http.StatusBadGateway, http.StatusNotFound,
			},
		},
		func(c context.Context, req *AppInfoRequest) (*Response, error) {

			if req.Raw && os.Getenv("DEV") != "1" {
				return nil, huma.Error403Forbidden("")
			}

			cacheKey := fmt.Sprint(req.AppID)

			if useMemCache && !req.Raw && !req.ForceRefresh {
				cached, ok := memcache.Get[*AppInfoWrapper](cache, cacheKey)
				if ok {
					slog.Debug("returning cached app info", "app_id", req.AppID)
					res := *cached
					if !req.ControllerSupport {
						res.ControllerSupport = nil
					}
					if !req.OfficialConfigs {
						res.OfficialConfigs = nil
					}
					return &Response{Body: &res}, nil
				}
			}

			dbInfo, dbErr := dal.AppInfo().GetAppInfo(c, req.AppID, appinfodal.AppInfoInclude{
				ControllerSupport: true,
				Assets:            true,
				Links:             true,
				Creators:          true,
				OfficialConfigs:   true,
			})
			if dbErr != nil && !errors.Is(dbErr, sql.ErrNoRows) {
				return nil, huma.Error502BadGateway("database error", dbErr)
			}

			if dbErr == nil && !req.ForceRefresh && time.Since(dbInfo.Timestamps.UpdatedAt) < dbMaxAge {
				wrapper := mapModelToResponse(dbInfo)
				if useMemCache {
					cache.Store(cacheKey, wrapper)
				}
				res := *wrapper
				if !req.ControllerSupport {
					res.ControllerSupport = nil
				}
				if !req.OfficialConfigs {
					res.OfficialConfigs = nil
				}
				return &Response{Body: &res}, nil
			}

			storeItem, err := fetchFromSteamAPI(c, req)
			if err != nil {
				return nil, err
			}

			if req.Raw {
				return &Response{
					Body: (*raw)(storeItem.raw),
				}, nil
			}

			if len(storeItem.raw.StoreItems) == 0 {
				return nil, huma.Error404NotFound("item not found")
			}

			appInfo := mapStoreItemToModel(storeItem.item)

			picsInfo, err := fetchFromSteamClient(c, sc, req.AppID)
			if err != nil {
				if strings.Contains(err.Error(), "no PICS info") {
					return nil, huma.Error404NotFound("item not found")
				}
				slog.Error("steam client fetch failed", "error", err)
				return nil, huma.Error502BadGateway("steam client error", err)
			}
			enrichModelFromPICS(appInfo, picsInfo)

			if dbErr != nil {
				if err := dal.AppInfo().InsertAppInfo(c, appInfo); err != nil {
					slog.Error("db insert failed", "error", err)
					return nil, err
				}
			} else {
				if err := dal.AppInfo().UpdateAppInfo(c, appInfo); err != nil {
					slog.Error("db update failed", "error", err)
					return nil, err
				}
			}

			wrapper := mapModelToResponse(appInfo)
			if useMemCache && wrapper.ControllerSupport != nil && wrapper.OfficialConfigs != nil {
				cache.Store(cacheKey, wrapper)
			}
			res := *wrapper
			if !req.ControllerSupport {
				res.ControllerSupport = nil
			}
			if !req.OfficialConfigs {
				res.OfficialConfigs = nil
			}
			return &Response{Body: &res}, nil
		},
	)
}

type storeAPIResult struct {
	raw  *steamapi.CStoreBrowse_GetItems_Response
	item *steamapi.StoreItem
}

func fetchFromSteamAPI(c context.Context, req *AppInfoRequest) (*storeAPIResult, error) {
	resp, err := steamapi.DefaultClient.GetItems(c, &steamapi.CStoreBrowse_GetItems_Request{
		Ids: []*steamapi.StoreItemID{
			{
				Appid: &req.AppID,
			},
		},
		Context: &steamapi.StoreBrowseContext{
			Language:    new("english"),
			CountryCode: new("US"),
		},
		DataRequest: &steamapi.StoreBrowseItemDataRequest{
			IncludeAssets:    new(true),
			IncludeBasicInfo: new(true),
			IncludeLinks:     new(true),
			IncludeRatings:   new(true),
			IncludePlatforms: new(true),
			IncludeRelease:   new(true),
		},
	})
	if err != nil {
		if errors.Is(err, steamapi.ErrRequest) {
			if strings.Contains(err.Error(), "HTTP error 404") {
				return nil, huma.Error404NotFound("app not found", err)
			}
			return nil, huma.Error502BadGateway("failed to get steam app info: %v", err)
		}
		return nil, err
	}
	result := &storeAPIResult{raw: resp}
	if len(resp.StoreItems) > 0 {
		result.item = resp.StoreItems[0]
	}
	return result, nil
}

func fetchFromSteamClient(ctx context.Context, sc client.Client, appID uint32) (*clientappinfo.Info, error) {
	infos, err := clientappinfo.Get(ctx, sc, appID)
	if err != nil {
		if errors.Is(err, client.ErrDisconnected) {
			if err := sc.Connect(ctx); err != nil {
				return nil, fmt.Errorf("reconnect: %w", err)
			}
			if err := sc.Login(ctx, client.LoginDetails{Anonymous: true, Language: "english"}); err != nil {
				return nil, fmt.Errorf("relogin: %w", err)
			}
			infos, err = clientappinfo.Get(ctx, sc, appID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	if len(infos) == 0 {
		return nil, fmt.Errorf("no PICS info for app %d", appID)
	}
	return &infos[0], nil
}

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
	if item.BasicInfo != nil {
		for _, pub := range item.BasicInfo.Publishers {
			if pub == nil || pub.Name == nil {
				continue
			}
			appInfo.CreatorLinks = append(appInfo.CreatorLinks, &models.AppCreatorToApp{
				RoleID: models.AppCreatorRoleIDPublisher,
				AppCreator: &models.AppCreator{
					Name:                 *pub.Name,
					CreatorClanAccountID: pub.CreatorClanAccountId,
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
					CreatorClanAccountID: dev.CreatorClanAccountId,
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
					CreatorClanAccountID: fr.CreatorClanAccountId,
				},
			})
		}
	}
	return appInfo
}

func enrichModelFromPICS(appInfo *models.AppInfo, info *clientappinfo.Info) {
	cs := &models.AppControllerSupport{AppID: appInfo.AppID}
	switch info.Common.ControllerSupport {
	case "full":
		cs.SupportLevel = new(types.ControllerSupportLevelFull)
	case "partial":
		cs.SupportLevel = new(types.ControllerSupportLevelPartial)
	default:
		cs.SupportLevel = new(types.ControllerSupportLevelNone)
	}

	if appInfo.Name == "" && info.Common.Name != "" {
		appInfo.Name = info.Common.Name
	}

	cat := info.Common.Category
	check := func(id int64) *bool {
		_, ok := cat[fmt.Sprintf("category_%d", id)]
		return &ok
	}
	cs.DS4WiredSupport = check(int64(clientappinfo.CategoryPS4Wired))
	cs.DS4WirelessSupport = check(int64(clientappinfo.CategoryPS4Bluetooth))
	cs.DS5WiredSupport = check(int64(clientappinfo.CategoryPS5Wired))
	cs.DS5WirelessSupport = check(int64(clientappinfo.CategoryPS5Bluetooth))
	cs.SteamInputAPISupport = check(int64(clientappinfo.CategorySteamInputAPI))
	appInfo.ControllerSupport = cs

	appInfo.OfficialConfigs = nil
	for idStr, detail := range info.Config.SteamControllerConfigDetails {
		appInfo.OfficialConfigs = appendConfigIfDefault(appInfo.OfficialConfigs, appInfo.AppID, idStr, detail)
	}
	for idStr, detail := range info.Config.SteamControllerTouchConfigDetails {
		appInfo.OfficialConfigs = appendConfigIfDefault(appInfo.OfficialConfigs, appInfo.AppID, idStr, detail)
	}
}

func appendConfigIfDefault(configs []*models.OfficialSteamInputConfig, appID uint32, idStr string, detail clientappinfo.ControllerConfigDetail) []*models.OfficialSteamInputConfig {
	branches := strings.Split(detail.EnabledBranches, ",")
	hasDefault := false
	for _, b := range branches {
		if strings.TrimSpace(b) == "default" {
			hasDefault = true
			break
		}
	}
	if !hasDefault {
		return configs
	}
	configID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return configs
	}
	return append(configs, &models.OfficialSteamInputConfig{
		AppID:          appID,
		ControllerType: steamtypes.ControllerType(detail.ControllerType),
		ConfigID:       configID,
	})
}

func mapModelToResponse(appInfo *models.AppInfo) *AppInfoWrapper {
	wrapper := &AppInfoWrapper{
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
	if appInfo.Release.SteamReleaseDate != nil {
		wrapper.AppItem.Release.SteamReleaseDate = *appInfo.Release.SteamReleaseDate
	}
	if appInfo.Release.OriginalReleaseDate != nil {
		wrapper.AppItem.Release.OriginalReleaseDate = *appInfo.Release.OriginalReleaseDate
	}
	if appInfo.Assets != nil {
		wrapper.AppItem.Assets = &steamapi.StoreItem_Assets{
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
		wrapper.AppItem.Links = &links
	}
	if appInfo.ShortDescription != nil || len(appInfo.CreatorLinks) > 0 {
		bi := &steamapi.StoreItem_BasicInfo{}
		if appInfo.ShortDescription != nil {
			bi.ShortDescription = appInfo.ShortDescription
		}
		for _, cl := range appInfo.CreatorLinks {
			if cl.AppCreator == nil {
				continue
			}
			link := &steamapi.StoreItem_BasicInfo_CreatorHomeLink{
				Name:                 &cl.AppCreator.Name,
				CreatorClanAccountId: cl.AppCreator.CreatorClanAccountID,
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
		wrapper.AppItem.BasicInfo = bi
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
	}
	if len(appInfo.OfficialConfigs) > 0 {
		oc := make(officialConfigs, len(appInfo.OfficialConfigs))
		for _, cfg := range appInfo.OfficialConfigs {
			oc[cfg.ControllerType] = configId(cfg.ConfigID)
		}
		wrapper.OfficialConfigs = &oc
	}
	return wrapper
}
