package main

import (
	"errors"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s path/to/kindle\n", os.Args[0])
		os.Exit(1)
	}
	kindleDir := os.Args[1]

	err := dirExists(kindleDir)
	if err != nil {
		usage()
	}
	err = dirExists(kindleDir + "/documents")
	if err != nil {
		usage()
	}
	err = dirExists(kindleDir + "/documents")
	if err != nil {
		usage()
	}
	err = dirExists(kindleDir + "/system")
	if err != nil {
		usage()
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

func usage() {
	fmt.Printf("Usage: %s path/to/kindle\n", os.Args[0])
	os.Exit(1)
}
