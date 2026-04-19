package models

import (
	"time"

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/types"
	"github.com/google/uuid"
)

type AppInfo struct {
	AppID            uint32     `bun:"app_id,pk"`
	Timestamps       Timestamps `bun:",embed"`
	Name             string     `bun:"name"`
	StoreURLPath     string     `bun:"store_url"`
	Type             string     `bun:"type"`
	ShortDescription *string    `bun:"short_description"`

	Platforms AppPlatforms `bun:",embed:platform_"`
	Release   AppRelease   `bun:",embed:release_"`

	ControllerSupport *AppControllerSupport       `bun:"rel:has-one,join:app_id=app_id"`
	Assets            *AppAsset                   `bun:"rel:has-one,join:app_id=app_id"`
	Links             []*AppLink                  `bun:"rel:has-many,join:app_id=app_id"`
	Creators          []*AppCreator               `bun:"m2m:app_creator_to_apps,join:AppInfo=AppCreator"`
	CreatorLinks      []*AppCreatorToApp          `bun:"rel:has-many,join:app_id=app_id"`
	OfficialConfigs   []*OfficialSteamInputConfig `bun:"rel:has-many,join:app_id=app_id"`
}

type AppPlatforms struct {
	Windows      *bool `bun:"windows"`
	Mac          *bool `bun:"mac"`
	SteamOSLinux *bool `bun:"steamos_linux"`
}

type AppRelease struct {
	SteamReleaseDate    *time.Time `bun:"steam_release_date"`
	OriginalReleaseDate *time.Time `bun:"original_release_date"`
}

type AppControllerSupport struct {
	AppID                uint32                        `bun:"app_id,pk"`
	Timestamps           Timestamps                    `bun:",embed"`
	SupportLevel         *types.ControllerSupportLevel `bun:"support_level"`
	DS4WiredSupport      *bool                         `bun:"ds4_wired_support"`
	DS4WirelessSupport   *bool                         `bun:"ds4_wireless_support"`
	DS5WiredSupport      *bool                         `bun:"ds5_wired_support"`
	DS5WirelessSupport   *bool                         `bun:"ds5_wireless_support"`
	SteamInputAPISupport *bool                         `bun:"steam_input_api_support"`
}

type AppAsset struct {
	AppID              uint32     `bun:"app_id,pk"`
	Timestamps         Timestamps `bun:",embed"`
	AssetURLFormat     *string    `bun:"asset_url_format"`
	MainCapsule        *string    `bun:"main_capsule"`
	SmallCapsule       *string    `bun:"small_capsule"`
	Header             *string    `bun:"header"`
	PackageHeader      *string    `bun:"package_header"`
	PageBackground     *string    `bun:"page_background"`
	HeroCapsule        *string    `bun:"hero_capsule"`
	HeroCapsule2X      *string    `bun:"hero_capsule_2x"`
	LibraryCapsule     *string    `bun:"library_capsule"`
	LibraryCapsule2X   *string    `bun:"library_capsule_2x"`
	LibraryHero        *string    `bun:"library_hero"`
	LibraryHero2X      *string    `bun:"library_hero_2x"`
	CommunityIcon      *string    `bun:"community_icon"`
	ClanAvatar         *string    `bun:"clan_avatar"`
	PageBackgroundPath *string    `bun:"page_background_path"`
	RawPageBackground  *string    `bun:"raw_page_background"`
}

type AppLink struct {
	Base  `bun:",embed"`
	AppID uint32 `bun:"app_id,notnull"`
	URL   string `bun:"url,notnull"`
}

type OfficialSteamInputConfig struct {
	Base           `bun:",embed"`
	AppID          uint32                    `bun:"app_id,notnull"`
	ControllerType steamtypes.ControllerType `bun:"controller_type,notnull"`
	ConfigID       uint64                    `bun:"config_id,notnull"`
}

type AppCreatorRoleID int

const (
	AppCreatorRoleIDPublisher AppCreatorRoleID = 1
	AppCreatorRoleIDDeveloper AppCreatorRoleID = 2
	AppCreatorRoleIDFranchise AppCreatorRoleID = 3
)

type AppCreatorRole struct {
	Base   `bun:",embed"`
	RoleID AppCreatorRoleID `bun:"role_id,notnull,unique"`
	Name   string           `bun:"name,notnull,unique"`
}

type AppCreator struct {
	Base                 `bun:",embed"`
	Name                 string `bun:"name,notnull,unique:creator_identity"`
	CreatorClanAccountID uint32 `bun:"creator_clan_account_id,unique:creator_identity"`

	Apps []*AppInfo `bun:"m2m:app_creator_to_apps,join:AppCreator=AppInfo"`
}

type AppCreatorToApp struct {
	Base         `bun:",embed"`
	AppID        uint32           `bun:"app_id,notnull"`
	AppCreatorID uuid.UUID        `bun:"app_creator_id,notnull,type:uuid"`
	RoleID       AppCreatorRoleID `bun:"role_id,notnull"`
	AppInfo      *AppInfo         `bun:"rel:belongs-to,join:app_id=app_id"`
	AppCreator   *AppCreator      `bun:"rel:belongs-to,join:app_creator_id=id"`
	Role         *AppCreatorRole  `bun:"rel:belongs-to,join:role_id=id"`
}
