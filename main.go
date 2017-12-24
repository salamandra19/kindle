package main

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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
}

var base []string

func filePath(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Print(err)
		return nil
	}
	if filepath.Dir(path) == kindleDir+"/documents" {
		return nil
	}

	switch filepath.Ext(path) {
	case ".mobi", ".pdf", ".prc", ".txt":
		a := "*" + fmt.Sprintf("%X", sha1.Sum([]byte("/mnt/us/documents"+strings.Replace(CollName(filepath.Dir(path)), "/", "-", -1)+"@en-US")))
		fmt.Println(a)
	default:
		if match(path) {
			b := "*" + fmt.Sprintf("%X", sha1.Sum([]byte("/mnt/us/documents"+strings.Replace(CollName(filepath.Dir(path)), "/", "-", -1)+"@en-US")))

			fmt.Println(b)
		}
	}
	return nil
}

func CollName(path string) string {
	return strings.TrimPrefix(path, kindleDir+"/documents/")
}

func match(s string) bool {
	re := regexp.MustCompile(`[.]azw.*$`)
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
