package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

func readDir(file string, w io.Writer) {
	dir, err := ioutil.ReadDir(file)
	if err != nil {
		fmt.Println(err)
	}
	for _, d := range dir {
		if strings.HasPrefix(d.Name(), ".") {
			continue
		}
		dfile := filepath.Join(file, d.Name())
		if d.IsDir() {
			readDir(dfile, w)
		} else {
			readFile(dfile, w)
		}
	}
}

//你好

func isHan(file, s string, num int, w io.Writer) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "//") || strings.HasPrefix(s, "/*") {
		return
	}
	for _, r := range s {
		if unicode.Is(unicode.Scripts["Han"], r) {
			fmt.Fprintf(w, "file:[%s] line:[%d] > [%s]\n", file, num, strings.TrimRight(s, "\n"))
			return
		}
	}
}

func readFile(file string, w io.Writer) {
	if strings.HasSuffix(file, ".m") || strings.HasSuffix(file, ".h") {
		absFile, err := filepath.Abs(file)
		if err != nil {
			fmt.Println(err)
			return
		}

		f, err := os.Open(absFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		rd := bufio.NewReader(f)
		num := 0
		for {
			line, err := rd.ReadString('\n')
			if io.EOF == err {
				num++
				isHan(file, line, num, w)
				break
			}
			if err != nil {
				break
			}
			num++
			isHan(file, line, num, w)
		}
	}
}

func readHan(file string, w io.Writer) {
	fileInfo, err := os.Stat(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	if fileInfo.IsDir() {
		readDir(file, w)
	}
	readFile(file, w)

}
