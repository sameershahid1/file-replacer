package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

func main() {
	urlList := []string{}
	f, err := excelize.OpenFile("sheet.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet 1")
	if err != nil {
		log.Fatal(err)
	}

	for _, row := range rows {
		urlList = append(urlList, row...)
	}

	totalFiles := len(urlList)
	fileChannels := make([]chan string, totalFiles)

	for i := 0; i < totalFiles; i++ {
		fileChannels[i] = make(chan string)
	}

	dirPath := "/home/sameer/Documents/allsydneytowtruck.com.au2/_websiteroot/"
	desPath := "/home/sameer/Documents/all-sydney"
	files, err := os.ReadDir(dirPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	for _, url := range urlList {
		splitStr := strings.Split(url, "https://allsydneytowtruck.com.au/")

		if len(splitStr) == 2 {
			for _, file := range files {
				extension := filepath.Ext(file.Name())
				if extension == ".php" && splitStr[1] == file.Name() {
					fileName := dirPath + file.Name()
					_ = os.Rename(fileName, filepath.Join(desPath, filepath.Base(fileName)))
				}
			}
		}
	}

	for i := 0; i < totalFiles; i++ {
		go func(url string, ch chan string) {
			ch <- url
		}(urlList[i], fileChannels[i])
	}

	for i := 0; i < totalFiles; i++ {
		<-fileChannels[i]
	}
}
