package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 || !isKindle(os.Args[1]) {
		fmt.Fprintf(os.Stderr, "Usage: %s path/to/kindle\n", os.Args[0])
		os.Exit(1)
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
