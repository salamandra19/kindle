package main

import (
	"log"
	"os"
)

func main() {
	files := os.Args[1:]
	if len(files) != 1 {
		log.Printf("введите %v\t и корректный адрес каталога", os.Args[0])
		return
	}

	checkAdress(files[0])
	checkAdress(files[0] + "/documents")
	checkAdress(files[0] + "/system")
}

func checkAdress(s string) bool {
	fileInfo, err := os.Stat(s)
	if err != nil {
		log.Printf("введите %v\t и корректный адрес каталога", os.Args[0])
		os.Exit(1)

	}
	return fileInfo.IsDir()
}
