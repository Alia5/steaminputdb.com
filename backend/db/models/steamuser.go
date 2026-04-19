package models

type SteamUser struct {
	SteamID     uint64     `bun:",pk"`
	Timestamps  Timestamps `bun:",embed"`
	PersonaName string     `bun:"persona_name,notnull"`
	IsAdmin     bool       `bun:"is_admin,notnull,default:false"`
}
