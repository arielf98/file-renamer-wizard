package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type JsonStruct struct {
	Extension string `json:"extension"`
	FilePath  string `json:"dir_path"`
	OldName   string `json:"old_name"`
	Name      string `json:"name"`
	User1     string `json:"user_name_1"`
	User2     string `json:"user_name_2"`
}

func main() {

	// locale the month to indonesian month
	idnMonth := map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}

	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(exec)

	config, err := readJsonFile(exPath)
	if err != nil {
		fmt.Print("Error", err)
	}

	currentTime := time.Now()
	currentYear := currentTime.Year()
	currentMonth := currentTime.Local().Month().String()
	newFileName1 := fmt.Sprintf("Time Sheet %s %s %d %s", config.Name, currentMonth, currentYear, config.User1)
	newFileName2 := fmt.Sprintf("Time Sheet %s %s %d %s", config.Name, currentMonth, currentYear, config.User2)
	pathTarget := exPath + "/" + config.OldName + config.Extension

	found := findFile(config, exPath)

	if found {

		destFolder := fmt.Sprintf("Time Sheet %s %d", idnMonth[currentMonth], currentYear)

		err := os.Mkdir(exPath+"/"+destFolder, 0775)

		if err != nil {
			fmt.Println("Error creating destination folder:", err)
			return
		}

		destination1 := filepath.Join(destFolder, newFileName1)
		destination2 := filepath.Join(destFolder, newFileName2)
		fullDestinationPath1 := fmt.Sprintf("%s/%s", exPath, destination1)
		fullDestinationPath2 := fmt.Sprintf("%s/%s", exPath, destination2)

		// Copy the file to the first destination within the new folder
		err = copyFileToNewFolder(pathTarget, fullDestinationPath1)
		if err != nil {
			fmt.Println("Error copying file to destination 1:", err)
			return
		}
		fmt.Println("File copied to destination 1 successfully.")

		// Copy the file to the second destination within the new folder
		err = copyFileToNewFolder(pathTarget, fullDestinationPath2)
		if err != nil {
			fmt.Println("Error copying file to destination 2:", err)
			return
		}
		fmt.Println("File copied to destination 2 successfully.")
	}

}

func findFile(fileData JsonStruct, execPath string) bool {

	isFileFound := false

	fileName := fileData.OldName + fileData.Extension

	err := filepath.Walk(execPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && filepath.Base(path) == fileName {
			isFileFound = true
		}

		return nil
	})

	if err != nil {
		fmt.Print("Error :", err)
	}
	return isFileFound
}

func readJsonFile(filePath string) (config JsonStruct, err error) {

	// read json
	fileStrPath := fmt.Sprintf("%s/config.json", filePath)

	data, err := os.ReadFile(fileStrPath)

	if err != nil {
		fmt.Print("Error", err)
		return config, err
	}

	var fileConfig JsonStruct

	err = json.Unmarshal(data, &fileConfig)

	if err != nil {
		fmt.Print("Error Parsing Json", err)
		return config, err
	}

	return fileConfig, nil
}

func copyFileToNewFolder(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}

	defer source.Close()

	destination, err := os.Create(dst)

	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	if err != nil {
		return err
	}

	// Flush any buffered data to ensure the file is written completely
	err = destination.Sync()
	if err != nil {
		return err
	}

	return nil
}