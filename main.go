package main

import (
	"readHan/config"

	"errors"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/jroimartin/gocui"
	runewidth "github.com/mattn/go-runewidth"
)

const (
	FileView         = "file"
	FileFilterView   = "filefilter"
	StringFilterView = "stringfilter"
	ResultView       = "result"
	ErrorView        = "error"
	PopUpView        = "popup"
)

func getViewValue(g *gocui.Gui, name string) string {
	v, err := g.View(name)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(v.Buffer())
}

func setViewDefaults(v *gocui.View) {
	v.Frame = true
	v.Wrap = false
}

func setViewTextAndCursor(v *gocui.View, s string) {
	v.Clear()
	fmt.Fprint(v, s)
	v.SetCursor(len(s), 0)
}

func parseKey(k string) (interface{}, gocui.Modifier, error) {
	mod := gocui.ModNone
	if strings.Index(k, "Alt") == 0 {
		mod = gocui.ModAlt
		k = k[3:]
	}
	switch len(k) {
	case 0:
		return 0, 0, errors.New("Empty key string")
	case 1:
		if mod != gocui.ModNone {
			k = strings.ToLower(k)
		}
		return rune(k[0]), mod, nil
	}

	key, found := KEYS[k]
	if !found {
		return 0, 0, fmt.Errorf("Unknown key: %v", k)
	}
	return key, mod, nil
}

var KEYS = map[string]gocui.Key{
	"F1":             gocui.KeyF1,
	"F2":             gocui.KeyF2,
	"F3":             gocui.KeyF3,
	"F4":             gocui.KeyF4,
	"F5":             gocui.KeyF5,
	"F6":             gocui.KeyF6,
	"F7":             gocui.KeyF7,
	"F8":             gocui.KeyF8,
	"F9":             gocui.KeyF9,
	"F10":            gocui.KeyF10,
	"F11":            gocui.KeyF11,
	"F12":            gocui.KeyF12,
	"Insert":         gocui.KeyInsert,
	"Delete":         gocui.KeyDelete,
	"Home":           gocui.KeyHome,
	"End":            gocui.KeyEnd,
	"PageUp":         gocui.KeyPgup,
	"PageDown":       gocui.KeyPgdn,
	"ArrowUp":        gocui.KeyArrowUp,
	"ArrowDown":      gocui.KeyArrowDown,
	"ArrowLeft":      gocui.KeyArrowLeft,
	"ArrowRight":     gocui.KeyArrowRight,
	"CtrlTilde":      gocui.KeyCtrlTilde,
	"Ctrl2":          gocui.KeyCtrl2,
	"CtrlSpace":      gocui.KeyCtrlSpace,
	"CtrlA":          gocui.KeyCtrlA,
	"CtrlB":          gocui.KeyCtrlB,
	"CtrlC":          gocui.KeyCtrlC,
	"CtrlD":          gocui.KeyCtrlD,
	"CtrlE":          gocui.KeyCtrlE,
	"CtrlF":          gocui.KeyCtrlF,
	"CtrlG":          gocui.KeyCtrlG,
	"Backspace":      gocui.KeyBackspace,
	"CtrlH":          gocui.KeyCtrlH,
	"Tab":            gocui.KeyTab,
	"CtrlI":          gocui.KeyCtrlI,
	"CtrlJ":          gocui.KeyCtrlJ,
	"CtrlK":          gocui.KeyCtrlK,
	"CtrlL":          gocui.KeyCtrlL,
	"Enter":          gocui.KeyEnter,
	"CtrlM":          gocui.KeyCtrlM,
	"CtrlN":          gocui.KeyCtrlN,
	"CtrlO":          gocui.KeyCtrlO,
	"CtrlP":          gocui.KeyCtrlP,
	"CtrlQ":          gocui.KeyCtrlQ,
	"CtrlR":          gocui.KeyCtrlR,
	"CtrlS":          gocui.KeyCtrlS,
	"CtrlT":          gocui.KeyCtrlT,
	"CtrlU":          gocui.KeyCtrlU,
	"CtrlV":          gocui.KeyCtrlV,
	"CtrlW":          gocui.KeyCtrlW,
	"CtrlX":          gocui.KeyCtrlX,
	"CtrlY":          gocui.KeyCtrlY,
	"CtrlZ":          gocui.KeyCtrlZ,
	"Esc":            gocui.KeyEsc,
	"CtrlLsqBracket": gocui.KeyCtrlLsqBracket,
	"Ctrl3":          gocui.KeyCtrl3,
	"Ctrl4":          gocui.KeyCtrl4,
	"CtrlBackslash":  gocui.KeyCtrlBackslash,
	"Ctrl5":          gocui.KeyCtrl5,
	"CtrlRsqBracket": gocui.KeyCtrlRsqBracket,
	"Ctrl6":          gocui.KeyCtrl6,
	"Ctrl7":          gocui.KeyCtrl7,
	"CtrlSlash":      gocui.KeyCtrlSlash,
	"CtrlUnderscore": gocui.KeyCtrlUnderscore,
	"Space":          gocui.KeySpace,
	"Backspace2":     gocui.KeyBackspace2,
	"Ctrl8":          gocui.KeyCtrl8,
}

const WINDOWS_OS = "windows"

var VIEWS []string = []string{
	FileView,
	FileFilterView,
	StringFilterView,
	ResultView,
}

const MinWidth = 60
const MinHeight = 20

type App struct {
	viewIndex    int
	currentPopup string
}

func (a *App) NextView(g *gocui.Gui, v *gocui.View) error {
	a.viewIndex = (a.viewIndex + 1) % len(VIEWS)
	return a.setView(g)
}

func (a *App) PrevView(g *gocui.Gui, v *gocui.View) error {
	a.viewIndex = (a.viewIndex - 1 + len(VIEWS)) % len(VIEWS)
	return a.setView(g)
}

func (a *App) closePopup(g *gocui.Gui, viewname string) {
	_, err := g.View(viewname)
	if err == nil {
		a.currentPopup = ""
		g.DeleteView(viewname)
		g.SetCurrentView(VIEWS[a.viewIndex%len(VIEWS)])
		g.Cursor = true
	}
}

func (a *App) setView(g *gocui.Gui) error {
	a.closePopup(g, a.currentPopup)
	_, err := g.SetCurrentView(VIEWS[a.viewIndex])
	return err
}

func (a *App) setKey(g *gocui.Gui, keyStr, commandStr, viewName string) error {
	if commandStr == "" {
		return nil
	}
	key, mod, err := parseKey(keyStr)
	if err != nil {
		return err
	}
	commandParts := strings.SplitN(commandStr, " ", 2)
	command := commandParts[0]
	var commandArgs string
	if len(commandParts) == 2 {
		commandArgs = commandParts[1]
	}
	keyFnGen, found := COMMANDS[command]
	if !found {
		return fmt.Errorf("Unknown command: %v", command)
	}
	keyFn := keyFnGen(commandArgs, a)
	if err := g.SetKeybinding(viewName, key, mod, keyFn); err != nil {
		return fmt.Errorf("Failed to set key '%v': %v", keyStr, err)
	}
	return nil
}

func (a *App) Find(g *gocui.Gui, _ *gocui.View) error {
	vres, _ := g.View("result")
	vres.Clear()
	file := getViewValue(g, "file")
	readHan(file, vres)

	return nil
}

func (a *App) SetKeys(g *gocui.Gui) error {
	for viewName, keys := range config.DefaultKeys {
		if viewName == "global" {
			viewName = ""
		}
		for keyStr, commandStr := range keys {
			if err := a.setKey(g, keyStr, commandStr, viewName); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *App) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if maxX < MinWidth || maxY < MinHeight {
		if v, err := g.SetView(ErrorView, 0, 0, maxX-1, maxY-1); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			setViewDefaults(v)
			v.Title = "Error"
			g.Cursor = false
			fmt.Fprintln(v, "Terminal is too small")
		}
	}
	if _, err := g.View(ErrorView); err == nil {
		g.DeleteView("error")
		g.Cursor = true
		a.setView(g)
	}

	//splitY := int(0.25 * float32(maxY-3))
	if v, err := g.SetView(FileView, 0, 0, maxX-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setViewDefaults(v)
		v.Title = "File - press ctrl+r to find"
		v.Editable = true
		v.Editor = &defaultEditor
		tmp, _ := os.Getwd()
		setViewTextAndCursor(v, tmp)
	}

	if v, err := g.SetView(FileFilterView, 0, 3, maxX-1, 6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setViewDefaults(v)
		v.Editable = true
		v.Title = `File filter`
		v.Editor = &defaultEditor
		setViewTextAndCursor(v, "*")
	}

	if v, err := g.SetView(StringFilterView, 0, 6, maxX-1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		setViewDefaults(v)
		v.Editable = true
		v.Title = "String filter"
		v.Editor = &defaultEditor
	}
	if v, err := g.SetView(ResultView, 0, 9, maxX-1, maxY-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = true
		v.Wrap = true
		v.Title = "Results"
		v.Editable = true
		v.Editor = &ViewEditor{a, g, false, gocui.EditorFunc(func(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
			return
		})}
	}

	return nil
}

func (a *App) ParseArgs(g *gocui.Gui) error {
	a.Layout(g)
	g.SetCurrentView(VIEWS[a.viewIndex])
	return nil
}

func initApp(a *App, g *gocui.Gui) {
	g.Cursor = true
	g.InputEsc = false
	g.BgColor = gocui.ColorDefault
	g.FgColor = gocui.ColorDefault
	g.SetManagerFunc(a.Layout)
}

func main() {
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}

	if runtime.GOOS == WINDOWS_OS && runewidth.IsEastAsian() {
		g.ASCII = true
	}

	app := &App{}

	// overwrite default editor
	defaultEditor = ViewEditor{app, g, false, gocui.DefaultEditor}

	initApp(app, g)

	err = app.ParseArgs(g)

	if err != nil {
		g.Close()
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	err = app.SetKeys(g)
	if err != nil {
		g.Close()
		fmt.Println("Error!", err)
		os.Exit(1)
	}
	defer g.Close()

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
