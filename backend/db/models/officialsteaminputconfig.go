package models

import (
	"time"

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
