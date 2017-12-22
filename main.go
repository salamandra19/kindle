package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var kindleDir string

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 || !isKindle(os.Args[1]) {
		log.Fatalf("Usage: %s path/to/kindle", os.Args[0])
	}
	kindleDir = os.Args[1]
	err := filepath.Walk(kindleDir, filePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, path := range base {
		fmt.Println(path)
	}
}

var base []string

func filePath(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}

	switch filepath.Ext(path) {
	case (".mobi"), (".pdf"), (".prc"), (".txt"):
		base = append(base, path)
	}
	if match(path) {
		base = append(base, path)
	}
	return nil
}

func match(s string) bool {
	re := regexp.MustCompile(".azw")
	if !re.MatchString(s) {
		return false
	}
	return true
}

var ErrNotADir = errors.New("not a directory")

func dirExists(dir string) error {
	fileInfo, err := os.Stat(dir)
	if err != nil {
		return err
	}
	if !fileInfo.IsDir() {
		return ErrNotADir
	}
	return nil
}

func isKindle(dir string) bool {
	return dirExists(dir) == nil &&
		dirExists(dir+"/documents") == nil &&
		dirExists(dir+"/system") == nil
}
