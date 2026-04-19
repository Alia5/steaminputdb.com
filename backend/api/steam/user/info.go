package user

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/Alia5/steaminputdb.com/api/ctx"
	"github.com/Alia5/steaminputdb.com/api/memcache"
	"github.com/Alia5/steaminputdb.com/api/steam/auth"
	"github.com/Alia5/steaminputdb.com/db"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/danielgtaylor/huma/v2"
)

type PlayerInfo struct {
	CommunityVisibilityState int                    `json:"communityvisibilitystate"` // TODO:enum
	PersonaName              string                 `json:"personaname"`
	ProfileURL               string                 `json:"profileurl"`
	Avatar                   string                 `json:"avatar"`
	AvatarMedium             string                 `json:"avatarmedium"`
	AvatarFull               string                 `json:"avatarfull"`
	AvatarHash               string                 `json:"avatarhash"`
	LastLogOff               time.Time              `json:"lastlogoff"`
	PrimaryClanID            string                 `json:"primaryclanid"`
	TimeCreated              time.Time              `json:"timecreated"`
	LocCountryCode           string                 `json:"loccountrycode"`
	AvatarFrame              *AvatarFrame           `json:"avatarframe,omitempty"`
	ProfileBackground        *string                `json:"profilebackground,omitempty"`
	MiniProfileBackground    *MiniProfileBackground `json:"mini_profile_background,omitempty"`
}

type SteamInputDBInfo struct {
	IsAdmin      *bool     `json:"is_admin,omitempty,omitzero"`
	RegisteredAt time.Time `json:"registered_at"`
}

type UserInfoResponse struct {
	PlayerInfo
	SteamInputDBInfo *SteamInputDBInfo `json:"steam_input_db_info,omitempty,omitzero"`
}

type AvatarFrame struct {
	Small *string `json:"small"`
	Large *string `json:"large"`
}

type MiniProfileBackground struct {
	Image     *string `json:"image"`
	MovieWebm *string `json:"movie_webm"`
	MovieMp4  *string `json:"movie_mp4"`
}

type Response struct {
	Body UserInfoResponse
}

type UserInfoRequest struct {
	UserID                string `query:"user_id,omitempty,omitzero"` // this is the user_id query parameter, but it's named something else to avoid confusion with the steamID that can be extracted from the context. if this is 0, we'll attempt to use the steamID from the context instead.
	AvatarFrame           bool   `query:"include_avatar_frame,omitempty" default:"false"`
	ProfileBackground     bool   `query:"include_profile_background,omitempty" default:"false"`
	MiniProfileBackground bool   `query:"include_mini_profile_background,omitempty" default:"false"`
}

func RegisterRoutes(a huma.API, dal db.DAL, opts ...bool) {
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
			Method:  http.MethodGet,
			Path:    "/v1/steam/userinfo",
			Tags:    []string{"steam", "user"},
			Summary: "Get Steam user info",
			Description: `Retrieve user information from Steam for the provided userId,  
or for attempt authenticated user if no userId is provided  
Returns 401 if no id provided and token is invalid and 400 if everything is missing`,
			Errors: []int{
				http.StatusBadGateway, http.StatusUnauthorized, http.StatusNotFound, http.StatusBadRequest,
			},
			Middlewares: huma.Middlewares{
				auth.ExtractSteamIDMiddleware,
			},
		},
		func(c context.Context, req *UserInfoRequest) (*Response, error) {
			var steamID string
			queryJSON, _ := json.Marshal(req)

			if req.UserID == "" {
				var ok bool
				steamID, ok = c.Value(ctx.KeySteamID).(string)
				if !ok || steamID == "" {
					return nil, huma.Error401Unauthorized("authentication error")
				}
			} else {
				steamID = req.UserID
				if useMemCache {
					cached, ok := memcache.Get[*Response](
						cache,
						string(queryJSON),
					)
					if ok {
						return cached, nil
					}
				}
			}
			if steamID == "" {
				return nil, huma.Error400BadRequest("no user ID provided")
			}

			info, err := steamapi.DefaultClient.GetPlayerSummaries(c, steamID)
			if err != nil {
				if errors.Is(err, steamapi.ErrRequest) {
					return nil, huma.Error502BadGateway("failed to get steam user info: %v", err)
				}
				return nil, err
			}

			if len(info.Response.Players) == 0 {
				return nil, huma.Error404NotFound("steam user not found")
			}
			player := info.Response.Players[0]

			infoResp := UserInfoResponse{
				PlayerInfo: PlayerInfo{
					CommunityVisibilityState: player.Communityvisibilitystate,
					PersonaName:              player.Personaname,
					ProfileURL:               player.Profileurl,
					Avatar:                   player.Avatar,
					AvatarMedium:             player.Avatarmedium,
					AvatarFull:               player.Avatarfull,
					AvatarHash:               player.Avatarhash,
					LastLogOff:               time.Unix(int64(player.Lastlogoff), 0).UTC(),
					PrimaryClanID:            player.Primaryclanid,
					TimeCreated:              time.Unix(int64(player.Timecreated), 0).UTC(),
					LocCountryCode:           player.Loccountrycode,
				},
			}

			res := &Response{
				Body: infoResp,
			}

			steamID64, err := strconv.ParseUint(steamID, 10, 64)
			if err != nil {
				slog.Error("failed to parse steamID to uint64", "error", err, "steamID", steamID)
				return res, nil
			}

			if req.UserID == "" {
				userInfo, err := dal.SteamUser().Get(c, steamID64)
				if err == nil {
					infoResp.SteamInputDBInfo = &SteamInputDBInfo{
						RegisteredAt: userInfo.Timestamps.CreatedAt,
					}
					if userInfo.IsAdmin {
						infoResp.SteamInputDBInfo.IsAdmin = new(true)
					}
				}
			}

			wg := sync.WaitGroup{}
			if req.AvatarFrame {
				wg.Go(func() {
					avatarFrame, err := steamapi.DefaultClient.GetAvatarFrame(c, &steamapi.CPlayer_GetAvatarFrame_Request{
						Steamid: &steamID64,
					})
					if err != nil {
						slog.Error("failed to get avatar frame", "err", err, "steam_id", steamID)
						return
					}
					if avatarFrame != nil {
						res.Body.AvatarFrame = &AvatarFrame{
							Small: avatarFrame.AvatarFrame.ImageSmall,
							Large: avatarFrame.AvatarFrame.ImageLarge,
						}
					}
				})
			}
			if req.ProfileBackground {
				wg.Go(func() {
					profileBackground, err := steamapi.DefaultClient.GetProfileBackground(c, &steamapi.CPlayer_GetProfileBackground_Request{
						Steamid: &steamID64,
					})
					if err != nil {
						slog.Error("failed to get profile background", "err", err, "steam_id", steamID)
						return
					}
					if profileBackground != nil {
						res.Body.ProfileBackground = profileBackground.ProfileBackground.ImageLarge
					}
				})
			}
			if req.MiniProfileBackground {
				wg.Go(func() {
					miniProfileBackground, err := steamapi.DefaultClient.GetMiniProfileBackground(c, &steamapi.CPlayer_GetMiniProfileBackground_Request{
						Steamid: &steamID64,
					})
					if err != nil {
						slog.Error("failed to get mini profile background", "err", err, "steam_id", steamID)
						return
					}
					if miniProfileBackground != nil {
						res.Body.MiniProfileBackground = &MiniProfileBackground{
							Image:     miniProfileBackground.ProfileBackground.ImageLarge,
							MovieWebm: miniProfileBackground.ProfileBackground.MovieWebm,
							MovieMp4:  miniProfileBackground.ProfileBackground.MovieMp4,
						}
					}
				})
			}
			wg.Wait()

			if useMemCache && req.UserID != "" && len(queryJSON) != 0 {
				cache.Store(string(queryJSON), res)
			}
			return res, nil
		},
	)
}
