package models

import "time"

type HWFeatureControllerType int

const (
	HWFeatureControllerFamilyCommon      HWFeatureControllerType = 0
	HWFeatureControllerFamilySteam       HWFeatureControllerType = 1
	HWFeatureControllerFamilyPlayStation HWFeatureControllerType = 2
	HWFeatureControllerFamilyXbox        HWFeatureControllerType = 3
	HWFeatureControllerFamilyNintendo    HWFeatureControllerType = 4
)

type HWFeature int

const (
	HWFeatureMotionInputs          HWFeature = 0
	HWFeatureNativeGyroCamera      HWFeature = 1
	HWFeatureRumble                HWFeature = 2
	HWFeatureHDHaptics             HWFeature = 3
	HWFeatureLightbar              HWFeature = 4
	HWFeatureAdaptiveTriggersWired HWFeature = 5
	HWFeatureImpulseTriggers       HWFeature = 6
	HWFeatureTouchpads             HWFeature = 7
	HWFeatureAudioHaptics          HWFeature = 8
)

type AppHWFeatures struct {
	AppID                   uint32 `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
	HWFeature               HWFeature               `gorm:"primaryKey;autoIncrement:false"`
	HWFeatureControllerType HWFeatureControllerType `gorm:"primaryKey;autoIncrement:false"`
	Notes                   string
}
