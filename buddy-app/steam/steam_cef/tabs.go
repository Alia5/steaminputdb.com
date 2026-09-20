package steamcef

import "strings"

const DefaultInjectTab = "SharedJSContext"
const DesktopTab = "Steam"

var BigPictureWindowTitles = []string{
	"وضع الصورة الكبيرة لـSteam",
	"Steam — Modo Big Picture",
	"Steam режим „Голям екран“",
	"Steam – režim Big Picture",
	"Steam Big Picture-tilstand",
	"Steam: Big Picture-modus",
	"Steam Big Picture Mode",
	"Steamin televisiotila",
	"Steam : mode Big Picture",
	"Big-Picture-Modus",
	"Steam – Λειτουργία Big Picture",
	"Steam Nagy Kép mód",
	"Mode Big Picture Steam",
	"Modalità Big Picture di Steam",
	"Steam Big Pictureモード",
	"Steam Big Picture 모드",
	"Modo Big Picture de Steam",
	"Mod Gambar Besar Steam",
	"Steam – Big Picture-modus",
	"Tryb Big Picture Steam",
	"Steam: Big Picture",
	"Steam – modul Big Picture",
	"Режим Big Picture",
	"Steam 大屏幕模式",
	"蒸汽平台大屏幕模式",
	"Steams Big Picture-läge",
	"Steam Big Picture 模式",
	"โหมด Steam Big Picture",
	"Steam Geniş Ekran Modu",
	"Steam у режимі Big Picture",
	"Chế độ Big Picture trên Steam",
}

var UIModTabs = [][]string{
	{
		DesktopTab,
	},
	BigPictureWindowTitles,
	{
		DefaultInjectTab,
	},
}

func matchesTab(title string, tabs []string) bool {
	for _, tab := range tabs {
		if strings.EqualFold(title, tab) {
			return true
		}
	}
	return false
}

func IsUIModTab(title string) bool {
	for _, tabs := range UIModTabs {
		if matchesTab(title, tabs) {
			return true
		}
	}
	return false
}
