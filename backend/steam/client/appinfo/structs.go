package appinfo

type Info struct {
	AppID    uint32   `json:"appid"`
	Common   Common   `json:"common"`
	Extended Extended `json:"extended"`
	Config   Config   `json:"config"`
}

type Config struct {
	SteamControllerConfigDetails      map[string]ControllerConfigDetail `json:"steamcontrollerconfigdetails"`
	SteamControllerTouchConfigDetails map[string]ControllerConfigDetail `json:"steamcontrollertouchconfigdetails"`
}

type ControllerConfigDetail struct {
	ControllerType  string `json:"controller_type"`
	EnabledBranches string `json:"enabled_branches"`
	UseActionBlock  string `json:"use_action_block"`
}

type Common struct {
	Name                    string                     `json:"name"`
	Type                    string                     `json:"type"`
	ReleaseState            string                     `json:"releasestate"`
	OSList                  string                     `json:"oslist"`
	OSArch                  string                     `json:"osarch"`
	Icon                    string                     `json:"icon"`
	Logo                    string                     `json:"logo"`
	ClientIcon              string                     `json:"clienticon"`
	ClientTGA               string                     `json:"clienttga"`
	LinuxClientIcon         string                     `json:"linuxclienticon"`
	LogoSmall               string                     `json:"logo_small"`
	ControllerSupport       string                     `json:"controller_support"`
	SteamReleaseDate        string                     `json:"steam_release_date"`
	OriginalReleaseDate     string                     `json:"original_release_date"`
	MetacriticScore         int                        `json:"metacritic_score"`
	MetacriticName          string                     `json:"metacritic_name"`
	MetacriticURL           string                     `json:"metacritic_url"`
	MetacriticFullURL       string                     `json:"metacritic_fullurl"`
	SortAs                  string                     `json:"sortas"`
	ReviewScore             int                        `json:"review_score"`
	ReviewPercentage        int                        `json:"review_percentage"`
	ReviewScoreBombs        int                        `json:"review_score_bombs"`
	ReviewPercentageBombs   int                        `json:"review_percentage_bombs"`
	PrimaryGenre            uint32                     `json:"primary_genre"`
	GameID                  uint64                     `json:"gameid"`
	ParentAppID             uint32                     `json:"parent"`
	StoreAssetMtime         uint64                     `json:"store_asset_mtime"`
	HasAdultContent         bool                       `json:"has_adult_content"`
	HasAdultContentViolence bool                       `json:"has_adult_content_violence"`
	HasAdultContentSex      bool                       `json:"has_adult_content_sex"`
	MarketPresence          bool                       `json:"market_presence"`
	WorkshopVisible         bool                       `json:"workshop_visible"`
	CommunityHubVisible     bool                       `json:"community_hub_visible"`
	CommunityVisibleStats   bool                       `json:"community_visible_stats"`
	FreeOnDemand            bool                       `json:"freeondemand"`
	ExFGLS                  bool                       `json:"exfgls"`
	Category                map[string]int64           `json:"category"`
	Genres                  map[string]int64           `json:"genres"`
	StoreTags               map[string]int64           `json:"store_tags"`
	Languages               map[string]bool            `json:"languages"`
	SupportedLanguages      map[string]LanguageSupport `json:"supported_languages"`
	Associations            map[string]Association     `json:"associations"`
	SmallCapsule            map[string]string          `json:"small_capsule"`
	HeaderImage             map[string]string          `json:"header_image"`
	LibraryAssets           LibraryAssets              `json:"library_assets"`
	ContentDescriptors      map[string]int64           `json:"content_descriptors"`
	SteamDeckCompatibility  SteamDeckCompat            `json:"steam_deck_compatibility"`
	NameLocalized           map[string]string          `json:"name_localized"`
}

type LanguageSupport struct {
	Supported bool `json:"supported"`
	FullAudio bool `json:"full_audio"`
	Subtitles bool `json:"subtitles"`
}

type Association struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type LibraryAssets struct {
	LibraryCapsule string       `json:"library_capsule"`
	LibraryHero    string       `json:"library_hero"`
	LibraryLogo    string       `json:"library_logo"`
	LogoPosition   LogoPosition `json:"logo_position"`
}

type LogoPosition struct {
	PinnedPosition string  `json:"pinned_position"`
	WidthPct       float64 `json:"width_pct"`
	HeightPct      float64 `json:"height_pct"`
}

type SteamDeckCompat struct {
	Category      int             `json:"category"`
	TestTimestamp uint64          `json:"test_timestamp"`
	TestedBuildID uint64          `json:"tested_build_id"`
	Configuration SteamDeckConfig `json:"configuration"`
}

type SteamDeckConfig struct {
	RecommendedRuntime                string `json:"recommended_runtime"`
	RequiresH264                      bool   `json:"requires_h264"`
	RequiresInternetForSetup          bool   `json:"requires_internet_for_setup"`
	RequiresInternetForSingleplayer   bool   `json:"requires_internet_for_singleplayer"`
	RequiresManualKeyboardInvoke      bool   `json:"requires_manual_keyboard_invoke"`
	RequiresNonControllerLauncherNav  bool   `json:"requires_non_controller_launcher_nav"`
	SmallText                         bool   `json:"small_text"`
	SupportedInput                    string `json:"supported_input"`
	GamescopeFrameLimiterNotSupported bool   `json:"gamescope_frame_limiter_not_supported"`
	NonDeckDisplayGlyphs              bool   `json:"non_deck_display_glyphs"`
	PrimaryPlayerIsControllerSlot0    bool   `json:"primary_player_is_controller_slot_0"`
}

type Extended struct {
	Developer    string `json:"developer"`
	Publisher    string `json:"publisher"`
	Homepage     string `json:"homepage"`
	DeveloperURL string `json:"developer_url"`
	GameDir      string `json:"gamedir"`
	IsFreeApp    bool   `json:"isfreeapp"`
	DemoOfAppID  uint32 `json:"demoofappid"`
	DLCForAppID  uint32 `json:"dlcforappid"`
	ListOfDLC    string `json:"listofdlc"`
	Icon         string `json:"icon"`
}
