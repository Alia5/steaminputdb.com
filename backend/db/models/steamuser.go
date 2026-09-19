package models

import "time"

type SteamUser struct {
	SteamID     uint64 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PersonaName string `gorm:"not null"`
	IsAdmin     bool   `gorm:"not null;default:false"`
	TrustedMod  bool   `gorm:"not null;default:false"`
}
