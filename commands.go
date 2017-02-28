package main

import (
	"github.com/jroimartin/gocui"
)

type CommandFunc func(*gocui.Gui, *gocui.View) error

var COMMANDS map[string]func(string, *App) CommandFunc = map[string]func(string, *App) CommandFunc{
	"quit": func(_ string, _ *App) CommandFunc {
		return quit
	},
	"find": func(_ string, a *App) CommandFunc {
		return a.Find
	},
	"nextView": func(_ string, a *App) CommandFunc {
		return a.NextView
	},
	"prevView": func(_ string, a *App) CommandFunc {
		return a.PrevView
	},
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
