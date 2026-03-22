package models

import "time"

type Settings struct {
	Base `bun:",embed"`

	AddDesktopUIEntries    *bool `bun:",notnull,default:true"`
	AddBigPictureUIEntries *bool `bun:",notnull,default:true"`

	DesktopUseSteamBrowser *bool `bun:",notnull,default:false"`

	SteamWaitTimeout *time.Duration `bun:",notnull,default:60000000000"`
}
