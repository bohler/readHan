package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jroimartin/gocui"
)

func TestReadDir(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(path)
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, d := range dir {
		fmt.Println("name:", d.Name())
		if strings.HasPrefix(d.Name(), ".") {
			continue
		}
		dfile := filepath.Join(path, d.Name())
		if d.IsDir() {
			fmt.Println("dir:", dfile)
			readDir(dfile)
		} else {
			fmt.Println("file:", dfile)
			readFile(dfile)
		}
	}
}

func TestReadHan(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}
	reader.fileFilter = []string{"*"}
	fmt.Println(path)
	readHan(path)
}

func TestMain(t *testing.T) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("colors", maxX/2-7, maxY/2-12, maxX/2+7, maxY/2+13); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		for i := 0; i <= 7; i++ {
			for _, j := range []int{1, 4, 7} {
				fmt.Fprintf(v, "Hello \033[3%d;%dm 颜色!\033[0m\n", i, j)
			}
		}
	}
	return nil
}
