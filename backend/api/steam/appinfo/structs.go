package appinfo

import (
	"github.com/Alia5/steaminputdb.com/api/search/games"
	"github.com/Alia5/steaminputdb.com/db/models"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/types"
)

type AppInfoItem struct {
	games.AppItem
	ControllerSupport *ControllerSupport     `json:"controller_support,omitempty"`
	OfficialConfigs   *officialConfigs       `json:"official_configs,omitempty"`
	SteamInputDBInfo  *SIDBControllerSupport `json:"steaminputdb_info,omitempty,omitzero"`
}

type ControllerSupport struct {
	SupportLevel         *types.ControllerSupportLevel `json:"support_level,omitempty" enum:"0,1,2"`
	DS4WiredSupport      *bool                         `json:"ds4_wired_support,omitempty"`
	DS4WirelessSupport   *bool                         `json:"ds4_wireless_support,omitempty"`
	DS5WiredSupport      *bool                         `json:"ds5_wired_support,omitempty"`
	DS5WirelessSupport   *bool                         `json:"ds5_wireless_support,omitempty"`
	SteamInputAPISupport *bool                         `json:"steaminputapi_support,omitempty"`
}

type SIDBControllerSupport struct {
	ControllerSupportRating models.ControllerSupportRating `json:"controller_support_rating"`
	ControllerSupportNotes  *string                        `json:"controller_support_notes,omitempty"`
	MixedInputInfo          *MixedInputInfo                `json:"mixed_input,omitempty"`
	GlyphInfo               *GlyphInfo                     `json:"glyphs,omitempty"`
	SteamInputAPISupport    *SteamInputAPISupport          `json:"steaminputapi_support,omitempty"`
	HWFeatures              []HWFeature                    `json:"hw_features,omitempty"`
	HWFeatureNotes          *string                        `json:"hw_feature_notes,omitempty"`
}

type MixedInputInfo struct {
	MixedInputType    models.MixedInputSupportType `json:"type"`
	GlyphFlicker      *bool                        `json:"glyph_flicker,omitempty,omitzero" nullable:"true"`
	Notes             string                       `json:"notes,omitzero,omitempty"`
	MixedInputModURLs []MixedInputMods             `json:"mixed_input_mods,omitempty"`
}

type MixedInputMods struct {
	URL  string `json:"url" format:"url" required:"true"`
	Name string `json:"name,omitempty,omitzero"`
}

type GlyphInfo struct {
	AutoDetect   *bool                    `json:"autodetect,omitempty,omitzero" nullable:"true"`
	ManualSelect *bool                    `json:"manual_select,omitempty,omitzero" nullable:"true"`
	Notes        string                   `json:"notes,omitempty,omitzero"`
	Controllers  []GlyphControllerSupport `json:"controllers,omitempty"`
}

type GlyphControllerSupport struct {
	ControllerType *steamtypes.ControllerType `json:"controller_type,omitempty,omitzero" example:"controller_steamcontroller_gordon" doc:"Type of controller this configuration is designed for"`
	Notes          string                     `json:"notes,omitempty,omitzero"`
}

type SteamInputAPISupport struct {
	Glyphs []models.AppGlyphTagType `json:"glyphs,omitempty,omitzero"`

	CameraSupport  models.SteamInputCameraSupport    `json:"camera_support"`
	PixelsPer360   *string                           `json:"pixels_per_360,omitempty,omitzero" nullable:"true"`
	SupportTags    []models.SteamInputAPISupportType `json:"support_tags,omitempty"`
	SteamInputType *models.SteamInputType            `json:"steam_input_type,omitempty,omitzero" nullable:"true"`

	Notes string `json:"notes,omitempty,omitzero"`
}

type HWFeature struct {
	Feature          models.HWFeature               `json:"feature"`
	ControllerFamily models.HWFeatureControllerType `json:"controller_family"`
	Notes            string                         `json:"notes,omitempty,omitzero"`
}

type configID uint64
type officialConfigs map[steamtypes.ControllerType]configID
