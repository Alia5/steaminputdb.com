package appinfo

import (
	"github.com/Alia5/steaminputdb.com/api/search/games"
	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
	"github.com/Alia5/steaminputdb.com/steamapi"
	"github.com/Alia5/steaminputdb.com/types"
)

type Response struct {
	Body responseBody
}

type AppInfoRequest struct {
	AppID             uint32 `query:"app_id" required:"true"`
	Raw               bool   `query:"raw" default:"false"`
	ControllerSupport bool   `query:"controller_support" default:"false"`
	OfficialConfigs   bool   `query:"official_configs" default:"false"`
}

type responseBody interface {
	searchSuggestionsResponse()
}

type raw steamapi.CStoreBrowse_GetItems_Response

func (r *raw) searchSuggestionsResponse()            {}
func (r *AppInfoWrapper) searchSuggestionsResponse() {}

type AppInfoWrapper struct {
	games.AppItem
	ControllerSupport *ControllerSupport `json:"controller_support,omitempty"`
	OfficialConfigs   *officialConfigs   `json:"official_configs,omitempty"`
}

type ControllerSupport struct {
	SupportLevel         *types.ControllerSupportLevel `json:"support_level,omitempty" enum:"0,1,2"`
	DS4WiredSupport      *bool                         `json:"ds4_wired_support,omitempty"`
	DS4WirelessSupport   *bool                         `json:"ds4_wireless_support,omitempty"`
	DS5WiredSupport      *bool                         `json:"ds5_wired_support,omitempty"`
	DS5WirelessSupport   *bool                         `json:"ds5_wireless_support,omitempty"`
	SteamInputAPISupport *bool                         `json:"steam_input_api_support,omitempty"`
}

type configId uint64
type officialConfigs map[steamtypes.ControllerType]configId
