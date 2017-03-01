package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	log "github.com/bohler/lib/dlog"
	"github.com/jroimartin/gocui"
)

func init() {
	log.Set(
		log.SetConsoleEnable(false),
		log.SetFileEnable(true),
	)
}

type Reader struct {
	w          io.Writer
	file       string
	fileFilter []string
	prefix     []string
	suffix     []string
	result     string
}

var reader Reader

func readDir(file string) {
	dir, err := ioutil.ReadDir(file)
	if err != nil {
		return
	}
	for _, d := range dir {
		if strings.HasPrefix(d.Name(), ".") {
			continue
		}
		dfile := filepath.Join(file, d.Name())
		if d.IsDir() {
			readDir(dfile)
		} else {
			readFile(dfile)
		}
	}
}

func isHan(file, s string, num int) {
	s = strings.TrimSpace(s)
	//log.Log.Info(s)
	expect := true
	for _, pre := range reader.prefix {
		if pre == "" {
			continue
		}

		if strings.HasPrefix(s, pre) {
			expect = false
			break
		}
	}
	if !expect {
		return
	}

	for _, suf := range reader.suffix {
		if suf == "" {
			continue
		}
		if strings.HasSuffix(s, suf) {
			expect = false
			break
		}
	}

	if !expect {
		return
	}

	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			lineStr := fmt.Sprintf("file:[%s] line:[%d] > [%s]\n", file, num, strings.TrimRight(s, "\n"))
			reader.result += lineStr
			return
		}
	}
}

//无聊

func readFile(file string) {
	expect := false

	//log.Log.Debug(file)

	for _, fix := range reader.fileFilter {
		if strings.HasSuffix(file, fix) || fix == "*" {
			expect = true
			break
		}
	}
	if !strings.Contains(file, ".") {
		expect = false
	}

	if expect {
		absFile, err := filepath.Abs(file)
		if err != nil {
			log.Log.Error(err)
			return
		}

		f, err := os.Open(absFile)
		if err != nil {
			log.Log.Error(err)
			return
		}
		defer f.Close()
		rd := bufio.NewReader(f)
		num := 0
		for {
			line, err := rd.ReadString('\n')
			if io.EOF == err {
				num++
				isHan(file, line, num)
				break
			}
			if err != nil {
				break
			}
			num++
			isHan(file, line, num)
		}
	}
}

func readHan(file string) {
	reader.result = ""
	//log.Log.Info(reader)
	fileInfo, err := os.Stat(file)
	//log.Log.Info(len(reader.prefix), len(reader.suffix))
	if err != nil {
		log.Log.Error(err)
		return
	}
	if fileInfo.IsDir() {
		readDir(file)

	} else {
		readFile(file)
	}
	buf := bytes.NewBufferString(reader.result)
	reader.w.Write(buf.Bytes())
	if _, err := reader.w.(*gocui.View).Line(0); err != nil {
		reader.w.(*gocui.View).SetOrigin(0, 0)
	}
}
