package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
