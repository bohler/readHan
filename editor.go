package main

import (
	"github.com/jroimartin/gocui"
)

var defaultEditor ViewEditor

// The singleLineEditor removes multi lines capabilities
type singleLineEditor struct {
	wuzzEditor gocui.Editor
}

type ViewEditor struct {
	app           *App
	g             *gocui.Gui
	backTabEscape bool
	origEditor    gocui.Editor
}

func (e *ViewEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	// handle back-tab (\033[Z) sequence
	if e.backTabEscape {
		if ch == 'Z' {
			e.app.PrevView(e.g, nil)
			e.backTabEscape = false
			return
		} else {
			e.origEditor.Edit(v, 0, '[', gocui.ModAlt)
		}
	}
	if ch == '[' && mod == gocui.ModAlt {
		e.backTabEscape = true
		return
	}

	// disable infinite down scroll
	if key == gocui.KeyArrowDown && mod == gocui.ModNone {
		_, cY := v.Cursor()
		_, err := v.Line(cY)
		if err != nil {
			return
		}
	}

	e.origEditor.Edit(v, key, ch, mod)
}

// The singleLineEditor removes multi lines capabilities
func (e singleLineEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case (ch != 0 || key == gocui.KeySpace) && mod == 0:
		e.wuzzEditor.Edit(v, key, ch, mod)
		// At the end of the line the default gcui editor adds a whitespace
		// Force him to remove
		ox, _ := v.Cursor()
		if ox > 1 && ox >= len(v.Buffer())-2 {
			v.EditDelete(false)
		}
		return
	case key == gocui.KeyEnter:
		return
	case key == gocui.KeyArrowRight:
		ox, _ := v.Cursor()
		if ox >= len(v.Buffer())-1 {
			return
		}
	case key == gocui.KeyHome || key == gocui.KeyArrowUp:
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
		return
	case key == gocui.KeyEnd || key == gocui.KeyArrowDown:
		width, _ := v.Size()
		lineWidth := len(v.Buffer()) - 1
		if lineWidth > width {
			v.SetOrigin(lineWidth-width, 0)
			lineWidth = width - 1
		}
		v.SetCursor(lineWidth, 0)
		return
	}
	e.wuzzEditor.Edit(v, key, ch, mod)
}

type viewProperties struct {
	title    string
	frame    bool
	editable bool
	wrap     bool
	editor   gocui.Editor
	text     string
}
