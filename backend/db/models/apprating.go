package models

import "time"

type AppRating struct {
	AppID     uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
