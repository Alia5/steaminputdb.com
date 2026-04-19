package appinfo

import "fmt"

type ControllerCategory uint32

const (
	CategoryPartialControllerSupport ControllerCategory = 18
	CategoryFullControllerSupport    ControllerCategory = 28
	CategoryPS4Wired                 ControllerCategory = 55
	CategoryPS4Bluetooth             ControllerCategory = 56
	CategoryPS5Wired                 ControllerCategory = 57
	CategoryPS5Bluetooth             ControllerCategory = 58
	CategorySteamInputAPI            ControllerCategory = 59
)

var ControllerCategories = map[ControllerCategory]string{
	CategoryPartialControllerSupport: "Partial Controller Support",
	CategoryFullControllerSupport:    "Full Controller Support",
	CategoryPS4Wired:                 "PS4 Controller Support (Wired)",
	CategoryPS4Bluetooth:             "PS4 Controller Support (Bluetooth)",
	CategoryPS5Wired:                 "PS5 Controller Support (Wired)",
	CategoryPS5Bluetooth:             "PS5 Controller Support (Bluetooth)",
	CategorySteamInputAPI:            "Steam Input API Support",
}

func (c ControllerCategory) String() string {
	if s, ok := ControllerCategories[c]; ok {
		return s
	}
	return fmt.Sprintf("Unknown (%d)", c)
}

