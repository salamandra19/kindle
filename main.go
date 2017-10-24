package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, os.Args[0], "/path/to/kindle")
		return
	}
	kindleDir := os.Args[1]

	dirExists(kindleDir)
	dirExists(kindleDir + "/documents")
	dirExists(kindleDir + "/system")
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
