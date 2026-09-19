package models

import (
	"time"

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/types"
	"github.com/google/uuid"
	"gorm.io/gorm"
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

type AppControllerSupport struct {
	AppID                uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	SupportLevel         *types.ControllerSupportLevel
	DS4WiredSupport      *bool
	DS4WirelessSupport   *bool
	DS5WiredSupport      *bool
	DS5WirelessSupport   *bool
	SteamInputAPISupport *bool
}

type AppAsset struct {
	AppID              uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
	AssetURLFormat     *string
	MainCapsule        *string
	SmallCapsule       *string
	Header             *string
	PackageHeader      *string
	PageBackground     *string
	HeroCapsule        *string
	HeroCapsule2X      *string `gorm:"column:hero_capsule_2x"`
	LibraryCapsule     *string
	LibraryCapsule2X   *string `gorm:"column:library_capsule_2x"`
	LibraryHero        *string
	LibraryHero2X      *string `gorm:"column:library_hero_2x"`
	CommunityIcon      *string
	ClanAvatar         *string
	PageBackgroundPath *string
	RawPageBackground  *string
}

type AppLink struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	AppID     uint32 `gorm:"not null;index"`
	URL       string `gorm:"not null"`
}

func (l *AppLink) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}

type OfficialSteamInputConfig struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	AppID          uint32                    `gorm:"not null;index"`
	ControllerType steamtypes.ControllerType `gorm:"not null"`
	ConfigID       uint64                    `gorm:"not null"`
}

func (c *OfficialSteamInputConfig) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type AppCreatorRoleID int

const (
	AppCreatorRoleIDPublisher AppCreatorRoleID = 1
	AppCreatorRoleIDDeveloper AppCreatorRoleID = 2
	AppCreatorRoleIDFranchise AppCreatorRoleID = 3
)

type AppCreatorRole struct {
	RoleID    AppCreatorRoleID `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string `gorm:"not null;unique"`
}

type AppCreator struct {
	ID                   uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
	Name                 string `gorm:"not null;uniqueIndex:creator_identity"`
	CreatorClanAccountID uint32 `gorm:"uniqueIndex:creator_identity"`
}

func (c *AppCreator) BeforeCreate(_ *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

type AppCreatorToApp struct {
	ID           uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	AppID        uint32           `gorm:"not null;index"`
	AppCreatorID uuid.UUID        `gorm:"not null;type:uuid"`
	RoleID       AppCreatorRoleID `gorm:"not null"`
	AppCreator   *AppCreator
	Role         *AppCreatorRole
}

func (l *AppCreatorToApp) BeforeCreate(_ *gorm.DB) error {
	if l.ID == uuid.Nil {
		l.ID = uuid.New()
	}
	return nil
}
