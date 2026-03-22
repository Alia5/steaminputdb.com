package filedetails

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/search/configs"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
)

var ControllerConfigFileTypeReturn = uint32(12) // steam api returns 12, query uses 15? o.O

type Response struct {
	Body responseBody
}

type raw steamapi.CPublishedFile_GetDetails_Response

type responseBody interface {
	fileDetailsQueryResponse()
}

type ConfigDetailResponse struct {
	configs.ConfigItem
	PlaytimeSeconds  *uint64 `json:"playtime_seconds" example:"4474585092" doc:"Total playtime in seconds with the specified time period"`
	PlaytimeSessions *uint64 `json:"playtime_sessions" example:"1234" doc:"Total number of playtime sessions within the specified time period"`
}

func (r *raw) fileDetailsQueryResponse()                  {}
func (r *ConfigDetailResponse) fileDetailsQueryResponse() {}

type Request struct {
	FileID            uint32 `query:"file_id" required:"true"`
	PlaytimeStatsDays uint32 `query:"playtime_stats,omitempty,omitzero" default:"30" doc:"Number of days for playtime statistics"` // days
	Raw               bool   `query:"raw,omitempty" default:"false"`                                                               // if true, returns the raw steamapi.PublishedFileDetails instead of the processed ConfigResponseItem
}

func RegisterRoute(a huma.API, opts ...bool) {
	var useMemCache bool
	if len(opts) > 0 {
		useMemCache = opts[0]
	} else {
		useMemCache = true
	}
	cache := memcache.New(30*time.Minute, 1000)

	registry := a.OpenAPI().Components.Schemas

	huma.Register(
		a,
		huma.Operation{
			Method:  http.MethodGet,
			Path:    "/v1/steam/filedetails",
			Tags:    []string{"steam", "filedetails"},
			Summary: "Get Controller config details",
			Description: `
Retrieve details for a given controller config file ID.  
If a non-controller config file ID is provided, this will respond with a 404`,
			Errors: []int{
				http.StatusBadGateway, http.StatusNotFound,
			},
			Responses: map[string]*huma.Response{
				"200": {
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: &huma.Schema{
								OneOf: []*huma.Schema{
									registry.Schema(reflect.TypeFor[ConfigDetailResponse](), true, ""),
									registry.Schema(reflect.TypeFor[steamapi.CPublishedFile_GetDetails_Response](), true, ""),
								},
							},
						},
					},
				},
			},
		},
		func(c context.Context, req *Request) (*Response, error) {

			if useMemCache && !req.Raw {
				cached, ok := memcache.Get[*Response](
					cache,
					fmt.Sprintf("%v_%d", req.FileID, req.PlaytimeStatsDays),
				)
				if ok {
					return cached, nil
				}
			}

			info, err := steamapi.DefaultClient.GetFileDetails(c,
				&steamapi.CPublishedFile_GetDetails_Request{
					Publishedfileids: []uint64{
						uint64(req.FileID),
					},
					Includekvtags:       new(true),
					Includetags:         new(true),
					Includevotes:        new(true),
					Includemetadata:     new(true),
					ReturnPlaytimeStats: &req.PlaytimeStatsDays,
					Includereactions:    new(true),
				})
			if err != nil {
				if errors.Is(err, steamapi.ErrInfoNotFound) {
					return nil, huma.Error404NotFound("file not found")
				}
				if errors.Is(err, steamapi.ErrRequest) {
					return nil, huma.Error502BadGateway("failed to get steam file details: %v", err)
				}
				return nil, err
			}

			if len(info.Publishedfiledetails) == 0 {
				slog.Debug("no file details found", "file_id", req.FileID, "resp", info)
				return nil, huma.Error404NotFound("file not found")
			}

			if req.Raw {
				return &Response{
					Body: (*raw)(info),
				}, nil
			}

			var item *steamapi.PublishedFileDetails
			for _, itm := range info.Publishedfiledetails {
				if itm == nil {
					continue
				}
				if itm.FileType == nil || *itm.FileType != ControllerConfigFileTypeReturn {
					continue
				}
				item = itm
				break
			}
			if item == nil {
				slog.Debug("no file with correct type found", "file_id", req.FileID, "resp", info)
				return nil, huma.Error404NotFound("file not found")
			}

			resultInfo := &configs.ConfigItem{
				Title:       item.Title,
				Description: item.FileDescription,
				// AppID:       &appID,
				// AppIDString: &appIDStr,
				FileID:    item.Publishedfileid,
				FileName:  item.Filename,
				FileURL:   item.FileUrl,
				FileSize:  item.FileSize,
				CreatorID: fmt.Sprintf("%v", *item.Creator),
				// ControllerType: tags
				// ControllerTypeNice: tags
				// ControllerNative: tags
				// TimeCreated: item.TimeCreated,
				// TimeUpdated: item.TimeUpdated,
				// Playtime: item.LifetimePlaytime,
				LifetimePlaytimeSessions: item.LifetimePlaytimeSessions,
				Subscriptions:            item.LifetimeSubscriptions,
			}
			if item.LifetimePlaytime != nil {
				resultInfo.LifetimePlaytimeSeconds = item.LifetimePlaytime
			}
			if item.TimeCreated != nil {
				resultInfo.TimeCreated = time.Unix(int64(*item.TimeCreated), 0).UTC()
			}
			if item.TimeUpdated != nil {
				resultInfo.TimeUpdated = time.Unix(int64(*item.TimeUpdated), 0).UTC()
			}
			if item.Tags != nil {
				tags := make([]string, 0, len(item.Tags))
				for _, tag := range item.Tags {
					if tag == nil || tag.Tag == nil {
						continue
					}
					tags = append(tags, *tag.Tag)

					if *tag.Tag == "controller_native" {
						resultInfo.ControllerNative = true
					} else if strings.HasPrefix(*tag.Tag, "controller_") {
						resultInfo.ControllerType = (*steamtypes.ControllerType)(tag.Tag)
					}
				}
				resultInfo.Tags = &tags

			}
			if resultInfo.ControllerType != nil && *resultInfo.ControllerType != "" {
				resultInfo.ControllerTypeNice = resultInfo.ControllerType.NiceName()
			}
			if item.Kvtags != nil {
				var appIDStr string
				for _, tag := range item.Kvtags {
					if tag == nil || tag.Key == nil || tag.Value == nil {
						continue
					}
					kv := tag
					if *kv.Key == "app" {
						appIDStr = *kv.Value
						break
					}
				}
				var appID uint32
				if appIDStr != "" {
					var parsed uint64
					parsed, err = strconv.ParseUint(appIDStr, 10, 32)
					if err == nil {
						appID = uint32(parsed)
						resultInfo.AppID = &appID
					}
				}
				resultInfo.AppIDString = &appIDStr
			}

			if item.VoteData != nil {
				resultInfo.Votes.Score = item.VoteData.Score
				resultInfo.Votes.Up = item.VoteData.VotesUp
				resultInfo.Votes.Down = item.VoteData.VotesDown
			}
			res := &ConfigDetailResponse{
				ConfigItem: *resultInfo,
			}

			if item.PlaytimeStats != nil {
				if item.PlaytimeStats.PlaytimeSeconds != nil {
					res.PlaytimeSeconds = item.PlaytimeStats.PlaytimeSeconds
				}
				if item.PlaytimeStats.NumSessions != nil {
					res.PlaytimeSessions = item.PlaytimeStats.NumSessions
				}
			}

			response := &Response{
				Body: res,
			}
			if useMemCache {
				cache.Store(
					fmt.Sprintf("%v_%d", req.FileID, req.PlaytimeStatsDays),
					response,
				)
			}

			return response, nil
		},
	)
}
