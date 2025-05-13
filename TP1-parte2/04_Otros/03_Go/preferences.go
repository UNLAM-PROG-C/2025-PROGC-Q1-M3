package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
)

const PreferencesFileName = "visualizaciones.csv"
const UserIdColumn = 0

type Visualization struct {
	UserID   string
	UserName string
	Title    string
	Type     string
	Genre    string
}

type Preference struct {
	UserID          string `json:"user_id"`
	UserName        string `json:"user_name"`
	ChosenGenre     string `json:"chosen_genre"`
	ChosenType      string `json:"chosen_type"`
	Total           int    `json:"total"`
	DifferentGenres int    `json:"different_genres"`
}

func getRootDir() string {
	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Could not get callerFile")
		os.Exit(1)
	}
	return filepath.Dir(callerFile)
}

func getTotalUsers() map[int]struct{} {
	root := getRootDir()
	filePath := filepath.Join(root, PreferencesFileName)
	f, err := os.Open(filePath)

	if err != nil {
		fmt.Println("Could not open file visualizaciones.csv", err)
		return nil
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil
	}

	if len(records) < 2 {
		fmt.Println("No data was found in CSV")
		return nil
	}

	distinct_user := make(map[int]struct{})

	for _, record := range records[1:] {
		user_id, _ := strconv.Atoi(record[UserIdColumn])
		distinct_user[user_id] = struct{}{}
	}

	return distinct_user
}

func main() {
	totalUsers := getTotalUsers()

}
