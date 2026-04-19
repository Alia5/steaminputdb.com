package types

type ControllerSupportLevel int

const (
	ControllerSupportLevelNone ControllerSupportLevel = iota
	ControllerSupportLevelPartial
	ControllerSupportLevelFull
)
