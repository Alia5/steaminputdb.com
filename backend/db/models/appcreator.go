package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
