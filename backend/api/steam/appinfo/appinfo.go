package appinfo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/api/ctx"
	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/steam/auth"
	"github.com/Alia5/steaminputdb.com/db"
	appinfodal "github.com/Alia5/steaminputdb.com/db/dal/appinfo"
	"github.com/Alia5/steaminputdb.com/steam/client"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

const (
	dbMaxAge         = 24 * time.Hour
	reconnectTimeout = 10 * time.Second
)

func RegisterRoute(a huma.API, dal db.DAL, opts ...bool) {

	registry := a.OpenAPI().Components.Schemas

	var useMemCache bool
	if len(opts) > 0 {
		useMemCache = opts[0]
	} else {
		useMemCache = true
	}
	cache := memcache.New(30*time.Minute, 1000)

	sc := client.New()
	bgCtx := context.Background()
	scDetails := client.LoginDetails{Anonymous: true, Language: "english"}
	if err := sc.Connect(bgCtx); err != nil {
		slog.Error("steam client connect failed", "error", err)
	} else if err := sc.Login(bgCtx, scDetails); err != nil {
		slog.Error("steam client login failed", "error", err)
	}
	sc.EnableAutoReconnect(scDetails, reconnectTimeout)

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodGet,
			Path:        "/v1/steam/appinfo",
			Tags:        []string{"steam", "app"},
			Summary:     "Get Steam app info",
			Description: "Retrieve app information from Steam Store for a given app ID",
			Errors: []int{
				http.StatusBadGateway, http.StatusNotFound, http.StatusForbidden,
			},
			Responses: map[string]*huma.Response{
				"200": {
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: registry.Schema(reflect.TypeFor[AppInfoItem](), true, ""),
						},
					},
				},
			},
			Middlewares: huma.Middlewares{
				auth.ExtractSteamIDMiddleware,
			},
		},
		func(c context.Context, req *AppInfoRequest) (*AppInfoResponse, error) {

			rawAllowed, err := isAllowedRawInfo(c, dal, req)
			if err != nil {
				return nil, err
			}

			if req.Raw && !rawAllowed {
				return nil, huma.Error403Forbidden("")
			}

			cacheKey := fmt.Sprint(req.AppID)

			if useMemCache && !req.Raw && !req.ForceRefresh {
				cached, ok := memcache.Get[*AppInfoItem](cache, cacheKey)
				if ok {
					slog.Debug("returning cached app info", "app_id", req.AppID)
					data, err := json.Marshal(cached)
					if err != nil {
						return nil, err
					}
					var res AppInfoItem
					if err := json.Unmarshal(data, &res); err != nil {
						return nil, err
					}
					if !req.ControllerSupport {
						res.ControllerSupport = nil
					}
					if !req.OfficialConfigs {
						res.OfficialConfigs = nil
					}
					return &AppInfoResponse{Body: &res}, nil
				}
			}

			dbInfo, dbErr := dal.AppInfo().Get(c, req.AppID, appinfodal.AppInfoInclude{
				ControllerSupport: true,
				Assets:            true,
				Links:             true,
				Creators:          true,
				OfficialConfigs:   true,
				SteamInputDBInfo:  true,
			})
			if dbErr != nil && !errors.Is(dbErr, gorm.ErrRecordNotFound) {
				return nil, huma.Error502BadGateway("database error", dbErr)
			}

			if dbErr == nil && !req.ForceRefresh && time.Since(dbInfo.UpdatedAt) < dbMaxAge {
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
				if !req.SteamInputDBInfo {
					res.SteamInputDBInfo = nil
				}
				return &AppInfoResponse{Body: &res}, nil
			}

			storeItem, err := fetchFromSteamAPI(c, req)
			if err != nil {
				return nil, err
			}

			if storeItem.raw == nil {
				return nil, huma.Error404NotFound("item not found")
			}

			if req.Raw {
				return &AppInfoResponse{
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
				if err := dal.AppInfo().Insert(c, appInfo); err != nil {
					slog.Error("db insert failed", "error", err)
					return nil, err
				}
			} else {
				if err := dal.AppInfo().UpdateBaseInfo(c, appInfo); err != nil {
					slog.Error("db update failed", "error", err)
					return nil, err
				}
			}

			infoItem := mapModelToResponse(appInfo)
			if useMemCache {
				cache.Store(cacheKey, infoItem)
			}
			res := *infoItem
			if !req.ControllerSupport {
				res.ControllerSupport = nil
			}
			if !req.OfficialConfigs {
				res.OfficialConfigs = nil
			}
			if !req.SteamInputDBInfo {
				res.SteamInputDBInfo = nil
			}
			return &AppInfoResponse{Body: &res}, nil
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
			return nil, huma.Error502BadGateway("failed to get steam app info", err)
		}
		return nil, err
	}
	result := &storeAPIResult{raw: resp}
	if len(resp.StoreItems) > 0 {
		result.item = resp.StoreItems[0]
	}
	return result, nil
}

func isAllowedRawInfo(c context.Context, dal db.DAL, req *AppInfoRequest) (bool, error) {
	rawAllowed := os.Getenv("DEV") == "1"
	if !rawAllowed && (req.Raw || req.ForceRefresh) {
		steamID, ok := c.Value(ctx.KeySteamID).(string)
		if !ok || steamID == "" {
			return false, huma.Error403Forbidden("authentication error")
		}

		steamID64, err := strconv.ParseUint(steamID, 10, 64)
		if err != nil {
			return false, huma.Error403Forbidden("invalid steam ID")
		}

		userInfo, err := dal.SteamUser().Get(c, steamID64)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false, huma.Error403Forbidden("user not found")
			}
			return false, huma.Error502BadGateway("database error", err)
		}
		if !userInfo.IsAdmin {
			return false, huma.Error403Forbidden("insufficient permissions")
		}
		rawAllowed = true
	}
	return rawAllowed, nil
}
