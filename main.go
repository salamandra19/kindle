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

type Books struct {
	Items      []string
	LastAccess int64
}

var collection = make(map[string]*Books)

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
	for k, v := range collection {
		fmt.Printf("%v\t%v\n", k, v)
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
		coll := strings.Replace(Abs2KindlePath(filepath.Dir(path)), "/", "-", -1) + "@en-US"
		sha := fmt.Sprintf("*%X", sha1.Sum([]byte("/mnt/us/documents/"+Abs2KindlePath(path))))
		if collection[coll] == nil {
			collection[coll] = &Books{}
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		fmt.Println(fileInfo.ModTime().Unix())
		collection[coll].LastAccess = (fileInfo.ModTime().Unix() * 1000)
		if collection[coll].LastAccess < (fileInfo.ModTime().Unix() * 1000) {
			collection[coll].LastAccess = (fileInfo.ModTime().Unix() * 1000)
		}
		collection[coll].Items = append(collection[coll].Items, sha)
	default:
		if match(path) {
			coll := strings.Replace(Abs2KindlePath(filepath.Dir(path)), "/", "-", -1) + "@en-US"
			sha := fmt.Sprintf("*%X", sha1.Sum([]byte("/mnt/us/documents/"+Abs2KindlePath(path))))
			if collection[coll] == nil {
				collection[coll] = &Books{}
			}
			fileInfo, err := os.Stat(path)
			if err != nil {
				return err
			}
			fmt.Println(fileInfo.ModTime().Unix())
			if collection[coll].LastAccess < (fileInfo.ModTime().Unix() * 1000) {
				collection[coll].LastAccess = (fileInfo.ModTime().Unix() * 1000)
			}
			collection[coll].Items = append(collection[coll].Items, sha)
		}
	}
	return nil
}

func Abs2KindlePath(path string) string {
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
