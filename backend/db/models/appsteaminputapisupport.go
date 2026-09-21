package models

import (
	"time"
)

type SteamInputAPISupportType int

const (
	SteamInputAPISupportUnknown SteamInputAPISupportType = 0
	SteamInputAPISupportNone    SteamInputAPISupportType = 1
	SteamInputAPIInGameActions  SteamInputAPISupportType = 2
	SteamInputAPIActionSets     SteamInputAPISupportType = 3
	SteamInputAPIXInputActions  SteamInputAPISupportType = 4
	SteamInputAPIDiscontinued   SteamInputAPISupportType = 5
)

type SteamInputCameraSupport int

const (
	SteamInputCameraSupportUnknown SteamInputCameraSupport = 0
	SteamInputCameraSupportNone    SteamInputCameraSupport = 1
	SteamInputCameraSupportPartial SteamInputCameraSupport = 2
	SteamInputCameraSupportFull    SteamInputCameraSupport = 3
)

type AppSteamInputAPISupport struct {
	AppID     uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time

	SIAPITypes         []*AppSIAPITypes `gorm:"foreignKey:AppID"`
	SIAPICameraSupport SteamInputCameraSupport
	SIAPIPixelsPer360  *string
	SIAPICameraNotes   string

	Notes string
}

type AppSIAPITypes struct {
	AppID uint32                   `gorm:"primaryKey;autoIncrement:false"`
	Type  SteamInputAPISupportType `gorm:"primaryKey;autoIncrement:false;index"`
}
