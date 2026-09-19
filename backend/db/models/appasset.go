package models

import "time"

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
