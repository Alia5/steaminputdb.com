package models

import (
	"time"
)

type MixedInputSupportType int

const (
	MixedInputSupportUnknown MixedInputSupportType = 0
	MixedInputUnsupported    MixedInputSupportType = 1
	MixedInputSupported      MixedInputSupportType = 2
	MixedInputWithMod        MixedInputSupportType = 3
	MixedInputOther          MixedInputSupportType = 4
)

type AppMixedInputInfo struct {
	AppID                  uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	MixedInputSupport      MixedInputSupportType `gorm:"index"`
	MixedInputGlyphFlicker *bool
	MixedInputNotes        string
	MixedInputModLinks     []*MixedInputModLinks `gorm:"foreignKey:AppID"`
}

type MixedInputModLinks struct {
	AppID     uint32 `gorm:"primaryKey;autoIncrement:false"`
	URI       string `gorm:"primaryKey;autoIncrement:false;index"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
