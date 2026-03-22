//go:build windows

package main

import (
	"github.com/Alia5/steaminputdb.com/buddy-app/util"
)

func init() {
	if util.IsRunFromGUI() {
		util.HideConsoleWindow()
	}
}
