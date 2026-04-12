package steam_js

import (
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"
)

type GetAppsArgs struct {
	NonSteamOnly  bool
	InstalledOnly bool
}

type GetAppsExecutor interface {
	steamcef.Executor[*GetAppsArgs, []AppInfo]
}

//go:embed templates/dist/get_apps.js.tmpl
var getAppsJSTmpl string

var getAppsJS = template.Must(template.New("getApps").Delims("<<%", "%>>").Parse(getAppsJSTmpl))

func NewGetApps(cfg *appconfig.Steam) GetAppsExecutor {
	return steamcef.NewExecutor[*GetAppsArgs, []AppInfo](cfg, getAppsJS)
}

type PerClientData struct {
	ClientID                     string `json:"clientid"`
	ClientName                   string `json:"client_name"`
	DisplayStatus                int    `json:"display_status"`
	Installed                    *bool  `json:"installed,omitempty"`
	IsAvailableOnCurrentPlatform bool   `json:"is_available_on_current_platform"`
	StatusPercentage             *int   `json:"status_percentage,omitempty"`
}

type AppInfo struct {
	AppType                     int             `json:"app_type"`
	AppID                       uint32          `json:"appid"`
	BitfieldSupportedLanguages  string          `json:"bitfield_supported_languages"`
	CanonicalAppType            int             `json:"canonicalAppType"`
	DisplayName                 string          `json:"display_name"`
	DisplayNameELanguage        int             `json:"display_name_elanguage"`
	HeaderFilename              string          `json:"header_filename"`
	IconHash                    string          `json:"icon_hash"`
	LibraryCapsuleFilename      string          `json:"library_capsule_filename"`
	LocalCacheVersion           int64           `json:"local_cache_version"`
	MetacriticScore             int             `json:"metacritic_score"`
	MinutesPlaytimeForever      int             `json:"minutes_playtime_forever"`
	MinutesPlaytimeLastTwoWeeks int             `json:"minutes_playtime_last_two_weeks"`
	MostAvailableClientID       string          `json:"most_available_clientid"`
	NumberOfCopies              int             `json:"number_of_copies"`
	PerClientData               []PerClientData `json:"per_client_data"`
	ReviewPercentageWithBombs   int             `json:"review_percentage_with_bombs"`
	ReviewPercentageWithout     int             `json:"review_percentage_without_bombs"`
	ReviewScoreWithBombs        int             `json:"review_score_with_bombs"`
	ReviewScoreWithout          int             `json:"review_score_without_bombs"`
	RtLastTimePlayed            int64           `json:"rt_last_time_played"`
	RtLastTimePlayedOrInstalled int64           `json:"rt_last_time_played_or_installed"`
	RtOriginalReleaseDate       int64           `json:"rt_original_release_date"`
	RtPurchasedTime             int64           `json:"rt_purchased_time"`
	RtRecentActivityTime        int64           `json:"rt_recent_activity_time"`
	RtSteamReleaseDate          int64           `json:"rt_steam_release_date"`
	RtStoreAssetMtime           int64           `json:"rt_store_asset_mtime"`
	SelectedClientID            string          `json:"selected_clientid"`
	SortAs                      string          `json:"sort_as"`
	SteamHwCompatCategoryPacked int             `json:"steam_hw_compat_category_packed"`
	SubscribedTo                bool            `json:"subscribed_to"`
	VisibleInGameList           bool            `json:"visible_in_game_list"`
}
