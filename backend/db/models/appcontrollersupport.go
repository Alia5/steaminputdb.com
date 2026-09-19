package models

import (
	"time"

	"github.com/Alia5/steaminputdb.com/types"
)

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
