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

func getVisualizationsMap(records [][]string) map[int][]Visualization {
	visualizations := make(map[int][]Visualization)

	for _, record := range records[1:] {
		user_id, _ := strconv.Atoi(record[UserIdColumn])
		vis := Visualization{
			UserID:   record[UserIdColumn],
			UserName: record[1],
			Title:    record[2],
			Type:     record[3],
			Genre:    record[4],
		}
		visualizations[user_id] = append(visualizations[user_id], vis)
	}

	return visualizations
}

func getVisualizations(filePath string) (map[int][]Visualization, int) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Could not open file visualizaciones.csv", err)
		return nil, 0
	}

	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return nil, 0
	}

	if len(records) < 2 {
		fmt.Println("No data was found in CSV")
		return nil, 0
	}

	visualizations := getVisualizationsMap(records)

	distinctUsers := len(visualizations)

	return visualizations, distinctUsers
}

func main() {

	root := getRootDir()
	filePath := filepath.Join(root, PreferencesFileName)

	visualizations := getVisualizations(filePath)
}
