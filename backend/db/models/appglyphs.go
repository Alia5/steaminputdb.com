package models

import "github.com/Alia5/steaminputdb.com/steam/steamtypes"

type AppGlyphTagType int

const (
	AppGlyphTagUnknown              AppGlyphTagType = 0
	AppGlyphTagSteamInputAPIButtons AppGlyphTagType = 1
	AppGlyphTagInGameButtons        AppGlyphTagType = 2
	AppGlyphTagSteamInputAPITexts   AppGlyphTagType = 3
)

type AppGlyphs struct {
	AppID uint32 `gorm:"primaryKey;autoIncrement:false"`

	AutoGlyphDetect   *bool
	ManualGlyphSelect *bool
	GlyphCtrlSupport  []*AppGlyphCtrlSupport `gorm:"foreignKey:AppID"`
	GlyphTags         []*AppGlyphTag         `gorm:"foreignKey:AppID"`
	GlyphNotes        string
}

type AppGlyphCtrlSupport struct {
	AppID          uint32                    `gorm:"primaryKey;autoIncrement:false"`
	ControllerType steamtypes.ControllerType `gorm:"primaryKey;index"`
	Notes          string
}

type AppGlyphTag struct {
	AppID uint32          `gorm:"primaryKey;autoIncrement:false"`
	Tag   AppGlyphTagType `gorm:"primaryKey;autoIncrement:false;index"`
}
