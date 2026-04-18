package appinfo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/search/games"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
)

type responseBody interface {
	searchSuggestionsResponse()
}

type raw steamapi.CStoreBrowse_GetItems_Response

func (r *raw) searchSuggestionsResponse()            {}
func (r *AppInfoWrapper) searchSuggestionsResponse() {}

type AppInfoWrapper struct {
	games.AppItem
}

type Response struct {
	Body responseBody
}

type Request struct {
	AppID uint32 `query:"app_id" required:"true"`
	Raw   bool   `query:"raw" default:"false"`
}

func RegisterRoute(a huma.API, opts ...bool) {
	var useMemCache bool
	if len(opts) > 0 {
		useMemCache = opts[0]
	} else {
		useMemCache = true
	}
	cache := memcache.New(30*time.Minute, 1000)

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
		func(c context.Context, req *Request) (*Response, error) {

			if useMemCache && !req.Raw {
				cached, ok := memcache.Get[*Response](cache, fmt.Sprint(req.AppID))
				if ok {
					slog.Debug("returning cached app info", "app_id", req.AppID)
					return cached, nil
				}
			}

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

			if req.Raw {
				return &Response{
					Body: (*raw)(resp),
				}, nil
			}

			if len(resp.StoreItems) == 0 {
				return nil, huma.Error404NotFound("item not found")
			}

			storeItem := resp.StoreItems[0]

			responseItem := &games.AppItem{
				AppID:        storeItem.Appid,
				Name:         storeItem.Name,
				StoreURLPath: storeItem.StoreUrlPath,
			}
			if storeItem.Type != nil {
				responseItem.Type = games.TypeToString(storeItem.Type)
			}
			if storeItem.Links != nil {
				responseItem.Links = &[]string{}
				for _, link := range storeItem.Links {
					if link == nil || link.Url == nil {
						continue
					}
					*responseItem.Links = append(*responseItem.Links, *link.Url)
				}
			}
			responseItem.BasicInfo = storeItem.BasicInfo
			responseItem.Assets = storeItem.Assets
			if storeItem.Platforms != nil {
				responseItem.Platforms = games.AppsPlatforms{
					Windows:      storeItem.Platforms.Windows,
					SteamOSLinux: storeItem.Platforms.SteamosLinux,
					Mac:          storeItem.Platforms.Mac,
				}
			}

			if storeItem.Release != nil {
				responseItem.Release = games.AppsRelease{}
				if storeItem.Release.SteamReleaseDate != nil {
					responseItem.Release.SteamReleaseDate = time.Unix(int64(*storeItem.Release.SteamReleaseDate), 0)
				}
				if storeItem.Release.OriginalReleaseDate != nil && *storeItem.Release.OriginalReleaseDate != 0 {
					responseItem.Release.OriginalReleaseDate = time.Unix(int64(*storeItem.Release.OriginalReleaseDate), 0)
				}
			}

			res := &Response{
				Body: &AppInfoWrapper{
					*responseItem,
				},
			}
			if useMemCache {
				cache.Store(fmt.Sprint(req.AppID), res)
			}

			return res, nil
		},
	)
}
