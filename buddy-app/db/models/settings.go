package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Settings struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time

	AddDesktopUIEntries    *bool `gorm:"not null;default:true"`
	AddBigPictureUIEntries *bool `gorm:"not null;default:true"`

	DesktopUseSteamBrowser *bool `gorm:"not null;default:false"`

	SteamWaitTimeout *time.Duration `gorm:"not null;default:60000000000"`
}

func (s *Settings) BeforeCreate(_ *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
