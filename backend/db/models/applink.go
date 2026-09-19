package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
