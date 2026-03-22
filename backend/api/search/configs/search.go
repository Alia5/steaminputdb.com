package configs

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

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(a huma.API) {
	registry := a.OpenAPI().Components.Schemas

	huma.Register(
		a,
		huma.Operation{
			Method:      http.MethodPost,
			Path:        "/v1/search/configs",
			Tags:        []string{"search"},
			Summary:     "Search SteamInput configurations",
			Description: "Search for SteamInput configurations for Steam and non-Steam games",
			Errors: []int{
				422,
			},
			Responses: map[string]*huma.Response{
				"200": {
					Content: map[string]*huma.MediaType{
						"application/json": {
							Schema: &huma.Schema{
								OneOf: []*huma.Schema{
									registry.Schema(reflect.TypeFor[ConfigsResponse](), true, ""),
									registry.Schema(reflect.TypeFor[steamapi.CPublishedFile_QueryFiles_Response](), true, ""),
								},
							},
						},
					},
				},
			},
		},
		Handler,
	)
}

var SteamConfigsAppID = new(uint32(241100))
var ControllerConfigFileType = new(uint32(15)) // steam api returns 12, search uses 15? o.O

func Handler(ctx context.Context, req *Request) (*SearchResponse, error) {

	query := &steamapi.CPublishedFile_QueryFiles_Request{
		Appid:         SteamConfigsAppID,
		Filetype:      ControllerConfigFileType,
		ReturnKvTags:  new(true),
		ReturnDetails: new(true),
		RequiredKvTags: []*steamapi.CPublishedFile_QueryFiles_Request_KVTag{
			{
				Key:   new("visibility"),
				Value: new("public"),
			},
			{
				Key:   new("deleted"),
				Value: new("0"),
			},
		},
		ReturnTags:     new(true),
		ReturnMetadata: new(true),
		//
		ReturnVoteData: &req.Body.Include.Votes,
		//
		SearchText: &req.Body.QueryText,
		//
		Numperpage: &req.Body.Limit,
		Page:       &req.Body.Page,
		//
		QueryType: new((uint32)(req.Body.Rank.By.PublishedFileQueryType().Number())),
		Days:      &req.Body.Rank.TrendingPeriod,
	}
	slog.Debug("Querying Steam PublishedFile API for Configs with",
		"query", query,
		"queryType", *query.QueryType,
		"rank", req.Body.Rank,
		"include", req.Body.Include,
		"pagination", req.Body.ConfigPagination,
		"filter", req.Body.Filter,
	)

	if req.Body.Filter.AppID != "" {
		query.RequiredKvTags = append(query.RequiredKvTags, &steamapi.CPublishedFile_QueryFiles_Request_KVTag{
			Key:   new("app"),
			Value: &req.Body.Filter.AppID,
		})
	}

	if len(req.Body.Filter.Tags) > 0 {
		query.Requiredtags = append(query.Requiredtags, req.Body.Filter.Tags...)
	}
	if len(req.Body.Filter.ExcludedTags) > 0 {
		query.Excludedtags = append(query.Excludedtags, req.Body.Filter.ExcludedTags...)
	}

	if req.Body.Filter.Creator != nil && *req.Body.Filter.Creator != "" {
		creator64, err := strconv.ParseUint(*req.Body.Filter.Creator, 10, 64)
		if err != nil {
			return nil, huma.Error422UnprocessableEntity("invalid creator id", err)
		}
		accountID := uint32(creator64 & 0xFFFFFFFF)
		query.RequiredKvTags = append(query.RequiredKvTags, &steamapi.CPublishedFile_QueryFiles_Request_KVTag{
			Key:   new("owner"),
			Value: new(fmt.Sprintf("%d", accountID)),
		})

	}

	queryResp, err := steamapi.DefaultClient.QueryFiles(ctx, query)
	if err != nil {
		if errors.Is(err, steamapi.ErrRequest) {
			return nil, huma.Error502BadGateway("failed to get steam publishedFile info", err)
		}
		return nil, err
	}

	if req.Body.Raw {
		return &SearchResponse{
			Body: (*raw)(queryResp),
		}, nil
	}

	resultItems := make([]ConfigItem, len(queryResp.Publishedfiledetails))
	for i, item := range queryResp.Publishedfiledetails {
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
			}
		}
		resultItems[i] = ConfigItem{
			Title:       item.Title,
			Description: item.FileDescription,
			AppID:       &appID,
			AppIDString: &appIDStr,
			FileID:      item.Publishedfileid,
			FileName:    item.Filename,
			FileURL:     item.FileUrl,
			FileSize:    item.FileSize,
			CreatorID:   fmt.Sprintf("%v", item.Creator),
			// ControllerType: tags
			// ControllerTypeNice: tags
			// ControllerNative: tags
			// TimeCreated: item.TimeCreated,
			// TimeUpdated: item.TimeUpdated,
			// Playtime: item.LifetimePlaytime,
			LifetimePlaytimeSessions: item.LifetimePlaytimeSessions,
			Subscriptions:            item.LifetimeSubscriptions,
			// Votes: if requested
		}
		if item.LifetimePlaytime != nil {
			resultItems[i].LifetimePlaytimeSeconds = item.LifetimePlaytime
		}
		if item.TimeCreated != nil {
			resultItems[i].TimeCreated = time.Unix(int64(*item.TimeCreated), 0)
		}
		if item.TimeUpdated != nil {
			resultItems[i].TimeUpdated = time.Unix(int64(*item.TimeUpdated), 0)
		}
		if item.Tags != nil {
			for _, tag := range item.Tags {
				if tag == nil || tag.Tag == nil {
					continue
				}
				if *tag.Tag == "controller_native" {
					resultItems[i].ControllerNative = true
				} else if strings.HasPrefix(*tag.Tag, "controller_") {
					resultItems[i].ControllerType = (*steamtypes.ControllerType)(tag.Tag)
				}
			}
		}
		if resultItems[i].ControllerType != nil && *resultItems[i].ControllerType != "" {
			resultItems[i].ControllerTypeNice = resultItems[i].ControllerType.NiceName()
		}

		if req.Body.Include.Votes && item.VoteData != nil {
			resultItems[i].Votes.Score = item.VoteData.Score
			resultItems[i].Votes.Up = item.VoteData.VotesUp
			resultItems[i].Votes.Down = item.VoteData.VotesDown
		}

		if req.Body.Include.Tags && item.Tags != nil {
			tags := make([]string, 0, len(item.Tags))
			for _, tag := range item.Tags {
				if tag == nil || tag.Tag == nil {
					continue
				}
				tags = append(tags, *tag.Tag)
			}
			resultItems[i].Tags = &tags
		}

	}

	return &SearchResponse{
		Body: &ConfigsResponse{
			Total: int(*queryResp.Total),
			Items: resultItems,
		},
	}, nil
}
