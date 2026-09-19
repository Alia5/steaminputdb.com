package models

import "time"

type ControllerSupportRating int

const (
	ControllerSupportRatingNone     ControllerSupportRating = 0
	ControllerSupportRatingWood     ControllerSupportRating = 1
	ControllerSupportRatingStone    ControllerSupportRating = 2
	ControllerSupportRatingBronze   ControllerSupportRating = 3
	ControllerSupportRatingSilver   ControllerSupportRating = 4
	ControllerSupportRatingGold     ControllerSupportRating = 5
	ControllerSupportRatingPlatinum ControllerSupportRating = 6
)

type AppInfo struct {
	AppID            uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Name             string
	StoreURLPath     string `gorm:"column:store_url"`
	Type             string
	ShortDescription *string

	Platforms AppPlatforms `gorm:"embedded;embeddedPrefix:platform_"`
	Release   AppRelease   `gorm:"embedded;embeddedPrefix:release_"`

	ControllerSupport *AppControllerSupport       `gorm:"foreignKey:AppID"`
	Assets            *AppAsset                   `gorm:"foreignKey:AppID"`
	Links             []*AppLink                  `gorm:"foreignKey:AppID"`
	CreatorLinks      []*AppCreatorToApp          `gorm:"foreignKey:AppID"`
	OfficialConfigs   []*OfficialSteamInputConfig `gorm:"foreignKey:AppID"`

	MixedInputInfo       *AppMixedInputInfo       `gorm:"foreignKey:AppID"`
	SteamInputAPISupport *AppSteamInputAPISupport `gorm:"foreignKey:AppID"`
	HWFeatures           []*AppHWFeatures         `gorm:"foreignKey:AppID"`
	Glyphs               *AppGlyphs               `gorm:"foreignKey:AppID"`

	ControllerSupportRating ControllerSupportRating
	ControllerSupportNotes  *string
}

type AppPlatforms struct {
	Windows      *bool
	Mac          *bool
	SteamOSLinux *bool `gorm:"column:steamos_linux"`
}

type AppRelease struct {
	SteamReleaseDate    *time.Time
	OriginalReleaseDate *time.Time
}
