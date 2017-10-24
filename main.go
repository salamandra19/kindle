package main

import (
	"errors"
	"fmt"
	"os"
)

var usage = fmt.Errorf("Usage: %s path/to/kindle\n", os.Args[0])

func main() {
	if len(os.Args) != 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
	kindleDir := os.Args[1]

	err := dirExists(kindleDir)
	if err != nil {
		fmt.Print(usage)
		os.Exit(1)
	}
	err = dirExists(kindleDir + "/documents")
	if err != nil {
		fmt.Print(usage)
		os.Exit(1)
	}
	err = dirExists(kindleDir + "/documents")
	if err != nil {
		fmt.Print(usage)
		os.Exit(1)
	}
	err = dirExists(kindleDir + "/system")
	if err != nil {
		fmt.Print(usage)
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
