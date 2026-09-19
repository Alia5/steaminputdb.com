package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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
