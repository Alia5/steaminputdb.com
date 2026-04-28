package steamtypes

type ControllerType string

const (
	ControllerTypeXbox360                 ControllerType = "controller_xbox360"
	ControllerTypeXboxOne                 ControllerType = "controller_xboxone"
	ControllerTypeXboxElite               ControllerType = "controller_xboxelite"
	ControllerTypePS3                     ControllerType = "controller_ps3"
	ControllerTypePS4                     ControllerType = "controller_ps4"
	ControllerTypePS5                     ControllerType = "controller_ps5"
	ControllerTypePS5Edge                 ControllerType = "controller_ps5_edge"
	ControllerTypeSteamController2015     ControllerType = "controller_steamcontroller_gordon"
	ControllerTypeSteamController         ControllerType = "controller_triton"
	ControllerTypeSteamControllerHeadcrab ControllerType = "controller_steamcontroller_headcrab"
	ControllerTypeSwitchPro               ControllerType = "controller_switch_pro"
	ControllerTypeSwitch2Pro              ControllerType = "controller_switch2_pro"
	ControllerTypeSwitchJoyConLeft        ControllerType = "controller_switch_joycon_left"
	ControllerTypeSwitchJoyConRight       ControllerType = "controller_switch_joycon_right"
	ControllerTypeSwitchJoyConPair        ControllerType = "controller_switch_joycon_pair"
	ControllerTypeSteamDeck               ControllerType = "controller_neptune"
	ControllerType8BitDo                  ControllerType = "controller_8bitdo"
	ControllerTypeLegionGoS               ControllerType = "controller_legion_go_s"
	ControllerHoriSteamDeck               ControllerType = "controller_hori_steam"
	ControllerRogAlly                     ControllerType = "controller_rog_ally"
	//
	ControllerTypeGeneric ControllerType = "controller_generic"
	// ControllerTypeNative       ControllerType = "controller_native" --- IGNORE ---
	ControllerTypeMobileTouch ControllerType = "controller_mobile_touch"
	ControllerTypeAndroid     ControllerType = "controller_android"
)

func (c *ControllerType) NiceName() string {
	if name, ok := controllerNiceNames[*c]; ok {
		return name
	}
	return ""
}

var controllerNiceNames map[ControllerType]string = map[ControllerType]string{
	ControllerTypeXbox360:                 "Xbox 360",
	ControllerTypeXboxOne:                 "Xbox One",
	ControllerTypeXboxElite:               "Xbox Elite",
	ControllerTypePS3:                     "DualShock 3",
	ControllerTypePS4:                     "DualShock 4",
	ControllerTypePS5:                     "DualSense",
	ControllerTypePS5Edge:                 "DualSense Edge",
	ControllerTypeSteamController2015:     "Steam Controller (2015)",
	ControllerTypeSteamController:         "Steam Controller",
	ControllerTypeSteamControllerHeadcrab: "Steam Controller (Headcrab)",
	ControllerTypeSwitchPro:               "Nintendo Switch Pro",
	ControllerTypeSwitch2Pro:              "Nintendo Switch 2 Pro",
	ControllerTypeSwitchJoyConLeft:        "Nintendo Switch Joy-Con (Left)",
	ControllerTypeSwitchJoyConRight:       "Nintendo Switch Joy-Con (Right)",
	ControllerTypeSwitchJoyConPair:        "Nintendo Switch Joy-Con (Pair)",
	ControllerTypeSteamDeck:               "Steam Deck",
	ControllerType8BitDo:                  "8BitDo",
	ControllerTypeLegionGoS:               "Lenovo Legion Go S",
	ControllerHoriSteamDeck:               "Horipad Steam",
	ControllerRogAlly:                     "ASUS ROG Ally",
	//
	ControllerTypeGeneric: "Generic",
	// ControllerTypeNative:       "Native", --- IGNORE ---
	ControllerTypeMobileTouch: "Mobile Touch",
	ControllerTypeAndroid:     "Android",
}

type ControllerTypeNumber int

const (
	ControllerTypeNumberNone                    ControllerTypeNumber = -1
	ControllerTypeNumberUnknown                 ControllerTypeNumber = 0
	ControllerTypeNumberSteamControllerHeadcrab ControllerTypeNumber = 1
	ControllerTypeNumberSteamController2015     ControllerTypeNumber = 2
	ControllerTypeNumberSteamController         ControllerTypeNumber = 3
	ControllerTypeNumberSteamDeck               ControllerTypeNumber = 4
	ControllerTypeNumberSteamController2        ControllerTypeNumber = 10
	ControllerTypeNumberFrontPanelBoard         ControllerTypeNumber = 20
	ControllerTypeNumberUnknownNonSteam         ControllerTypeNumber = 30
	ControllerTypeNumberXbox360                 ControllerTypeNumber = 31
	ControllerTypeNumberXboxOne                 ControllerTypeNumber = 32
	ControllerTypeNumberPS3                     ControllerTypeNumber = 33
	ControllerTypeNumberPS4                     ControllerTypeNumber = 34
	ControllerTypeNumberWii                     ControllerTypeNumber = 35
	ControllerTypeNumberApple                   ControllerTypeNumber = 36
	ControllerTypeNumberAndroid                 ControllerTypeNumber = 37
	ControllerTypeNumberSwitchPro               ControllerTypeNumber = 38
	ControllerTypeNumberSwitchJoyConLeft        ControllerTypeNumber = 39
	ControllerTypeNumberSwitchJoyConRight       ControllerTypeNumber = 40
	ControllerTypeNumberSwitchJoyConPair        ControllerTypeNumber = 41
	ControllerTypeNumberSwitchInputOnly         ControllerTypeNumber = 42
	ControllerTypeNumberMobileTouch             ControllerTypeNumber = 43
	ControllerTypeNumberXInputSwitch            ControllerTypeNumber = 44
	ControllerTypeNumberPS5                     ControllerTypeNumber = 45
	ControllerTypeNumberXboxElite               ControllerTypeNumber = 46
	ControllerTypeNumberLast                    ControllerTypeNumber = 47
	ControllerTypeNumberGenericKeyboard         ControllerTypeNumber = 400
	ControllerTypeNumberGenericMouse            ControllerTypeNumber = 800
)

var eControllerTypeMap = map[ControllerTypeNumber]ControllerType{
	ControllerTypeNumberSteamControllerHeadcrab: ControllerTypeSteamControllerHeadcrab,
	ControllerTypeNumberSteamController2015:     ControllerTypeSteamController2015,
	ControllerTypeNumberSteamController:         ControllerTypeSteamController,
	ControllerTypeNumberSteamController2:        ControllerTypeSteamController,
	ControllerTypeNumberSteamDeck:               ControllerTypeSteamDeck,
	ControllerTypeNumberXbox360:                 ControllerTypeXbox360,
	ControllerTypeNumberXboxOne:                 ControllerTypeXboxOne,
	ControllerTypeNumberPS3:                     ControllerTypePS3,
	ControllerTypeNumberPS4:                     ControllerTypePS4,
	ControllerTypeNumberAndroid:                 ControllerTypeAndroid,
	ControllerTypeNumberSwitchPro:               ControllerTypeSwitchPro,
	ControllerTypeNumberSwitchJoyConLeft:        ControllerTypeSwitchJoyConLeft,
	ControllerTypeNumberSwitchJoyConRight:       ControllerTypeSwitchJoyConRight,
	ControllerTypeNumberSwitchJoyConPair:        ControllerTypeSwitchJoyConPair,
	ControllerTypeNumberMobileTouch:             ControllerTypeMobileTouch,
	ControllerTypeNumberPS5:                     ControllerTypePS5,
	ControllerTypeNumberXboxElite:               ControllerTypeXboxElite,
}

func EControllerTypeFromInt(t ControllerTypeNumber) ControllerType {
	if ct, ok := eControllerTypeMap[t]; ok {
		return ct
	}
	return ControllerTypeGeneric
}

var controllerTypeNumberNiceNames = map[ControllerTypeNumber]string{
	ControllerTypeNumberNone:            "None",
	ControllerTypeNumberUnknown:         "Unknown",
	ControllerTypeNumberFrontPanelBoard: "Front Panel Board",
	ControllerTypeNumberUnknownNonSteam: "Unknown Non-Steam Controller",
	ControllerTypeNumberWii:             "Wii Controller",
	ControllerTypeNumberApple:           "Apple Controller",
	ControllerTypeNumberSwitchInputOnly: "Switch Input Only",
	ControllerTypeNumberXInputSwitch:    "XInput Switch Controller",
	ControllerTypeNumberLast:            "Last Controller",
	ControllerTypeNumberGenericKeyboard: "Generic Keyboard",
	ControllerTypeNumberGenericMouse:    "Generic Mouse",
}

func (n ControllerTypeNumber) NiceName() string {
	if ct, ok := eControllerTypeMap[n]; ok {
		return ct.NiceName()
	}
	if name, ok := controllerTypeNumberNiceNames[n]; ok {
		return name
	}
	return ""
}
