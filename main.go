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
	Extension string   `json:"extension"`
	OldName   string   `json:"old_file_name"`
	NewName   string   `json:"new_file_name"`
	UsersName []string `json:"users_name"`
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

	// exec, err := os.Executable()
	// if err != nil {
	// 	panic(err)
	// }

	// exPath := filepath.Dir(exec) // for build only
	exPath := "." // for dev only

	config, err := readJsonFile(exPath)
	if err != nil {
		fmt.Print("Error", err)
		panic(err)
	}

	currentTime := time.Now()
	currentYear := currentTime.Year()
	currentMonth := currentTime.Local().Month().String()
	pathTarget := exPath + "/" + config.OldName + config.Extension

	found := findFile(config, exPath)

	if found {

		destFolder := fmt.Sprintf("Time Sheet %s %d", idnMonth[currentMonth], currentYear)

		err := os.Mkdir(exPath+"/"+destFolder, 0775)

		if err != nil {
			fmt.Println("Error creating destination folder:", err)
			return
		}

		users := config.UsersName

		//loop through the users
		for _, user := range users {
			newFileName := fmt.Sprintf("%s %s %d %s%s", config.NewName, idnMonth[currentMonth], currentYear, user, config.Extension)
			destination := filepath.Join(destFolder, newFileName)
			fullDestinationPath := fmt.Sprintf("%s/%s", exPath, destination)

			err := copyFileToNewFolder(pathTarget, fullDestinationPath)
			if err != nil {
				fmt.Println("Error copying file to destination", err)
				return
			}

			fmt.Println("File copied to destination successfully.")
		}
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
