package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) != 2 || !isKindle(os.Args[1]) {
		log.Fatalf("Usage: %s path/to/kindle\n", os.Args[0])
	}
	kindleDir := os.Args[1]
	err := filepath.Walk(kindleDir, filePath)
	if err != nil {
		log.Fatal(err)
	}
	matchDirs := []string{".mobi", ".azw", ".pdf", ".prc", ".txt"}
	err = filepath.Walk(kindleDir, matchFiles(matchDirs))
	if err != nil {
		log.Fatal(err)
	}
}

func filePath(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	return nil
}

func matchFiles(matchDirs []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}
		if info.IsDir() {
			dir := filepath.Base(path)
			for _, d := range matchDirs {
				if d != dir {
					return filepath.SkipDir
				}
			}
		}
		fmt.Println(path)
		return nil
	}
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
