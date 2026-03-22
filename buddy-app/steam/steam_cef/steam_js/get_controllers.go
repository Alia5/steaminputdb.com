package steam_js

import (
	_ "embed"
	"text/template"

	appconfig "github.com/Alia5/steaminputdb.com/buddy-app/config"
	steamcef "github.com/Alia5/steaminputdb.com/buddy-app/steam/steam_cef"

	"github.com/Alia5/steaminputdb.com/steam/steamtypes"
)

type GetControllersExecutor interface {
	steamcef.Executor[*struct{}, []ControllerInfo]
}

//go:embed templates/get_controllers.js.tmpl
var getControllersJSTmpl string

var getControllersJS = template.Must(template.New("getControllers").Parse(getControllersJSTmpl))

func NewGetControllers(cfg *appconfig.Steam) GetControllersExecutor {
	return steamcef.NewExecutor[*struct{}, []ControllerInfo](cfg, getControllersJS)
}

type ControllerInfo struct {
	Name                                 string                          `json:"strName"`
	ControllerType                       steamtypes.ControllerTypeNumber `json:"eControllerType"`
	ControllerStyle                      int                             `json:"eControllerStyle"`
	XInputIndex                          int64                           `json:"nXInputIndex"`
	ControllerIndex                      int                             `json:"nControllerIndex"`
	RumblePreference                     int                             `json:"eRumblePreference"`
	Wireless                             bool                            `json:"bWireless"`
	VendorID                             uint16                          `json:"unVendorID"`
	ProductID                            uint16                          `json:"unProductID"`
	Capabilities                         string                          `json:"unCapabilities"`
	FirmwareBuildTime                    string                          `json:"strFirmwareBuildTime"`
	SerialNumber                         string                          `json:"strSerialNumber"`
	ChipID                               string                          `json:"strChipID"`
	LEDBrightness                        int                             `json:"flLEDBrightness"`
	LEDSaturation                        int                             `json:"flLEDSaturation"`
	TurnOnSound                          int                             `json:"nTurnOnSound"`
	TurnOffSound                         int                             `json:"nTurnOffSound"`
	LEDColorR                            uint8                           `json:"nLEDColorR"`
	LEDColorG                            uint8                           `json:"nLEDColorG"`
	LEDColorB                            uint8                           `json:"nLEDColorB"`
	LStickDeadzone                       int                             `json:"nLStickDeadzone"`
	RStickDeadzone                       int                             `json:"nRStickDeadzone"`
	Haptics                              bool                            `json:"bHaptics"`
	SWAntiDrift                          bool                            `json:"bSWAntiDrift"`
	IMUOneEuroFilter                     bool                            `json:"bIMUOneEuroFilter"`
	LHapticStrength                      int                             `json:"nLHapticStrength"`
	RHapticStrength                      int                             `json:"nRHapticStrength"`
	LPadPressureCurve                    int                             `json:"flLPadPressureCurve"`
	RPadPressureCurve                    int                             `json:"flRPadPressureCurve"`
	LeftStickTouchDisablesLeftTrackPad   bool                            `json:"bLeftStickTouchDisablesLeftTrackPad"`
	RightStickTouchDisablesRightTrackPad bool                            `json:"bRightStickTouchDisablesRightTrackPad"`
	PlayerSlotLEDSetting                 int                             `json:"ePlayerSlotLEDSetting"`
	NintendoLayout                       bool                            `json:"bNintendoLayout"`
	UseReversedLayout                    bool                            `json:"bUseReversedLayout"`
	UseUniversalFaceButtonGlyphs         bool                            `json:"bUseUniversalFaceButtonGlyphs"`
	GyroStationaryTolerance              float64                         `json:"flGyroStationaryTolerance"`
	AccelerometerStationaryTolerance     float64                         `json:"flAccelerometerStationaryTolerance"`
	AuxCapSenseThreshold                 int                             `json:"nAuxCapSenseThreshold"`
	AuxCapSenseHysteresis                int                             `json:"nAuxCapSenseHysteresis"`
	RemoteDevice                         bool                            `json:"bRemoteDevice"`
	Bluetooth                            bool                            `json:"bBluetooth"`
	VecMacAddrs                          []string                        `json:"vecMacAddrs"`
	UcBatteryLevel                       uint8                           `json:"ucBatteryLevel"`
	Charging                             bool                            `json:"bCharging"`
	HasTouchscreen                       bool                            `json:"bHasTouchscreen"`
}
