package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
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

func getVisualizationsMap(records [][]string) map[string][]Visualization {
	visualizations := make(map[string][]Visualization)

	for _, record := range records[1:] {
		userId := record[UserIdColumn]
		vis := Visualization{
			UserID:   userId,
			UserName: record[1],
			Title:    record[2],
			Type:     record[3],
			Genre:    record[4],
		}
		visualizations[userId] = append(visualizations[userId], vis)
	}

	return visualizations
}

func getVisualizations(filePath string) map[string][]Visualization {
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

	visualizations := getVisualizationsMap(records)

	return visualizations
}

func getPreferred(field map[string]int) string {
	var maxKeyCount int
	var chosenFieldValue string

	for key, count := range field {
		if count > maxKeyCount {
			maxKeyCount = count
			chosenFieldValue = key
		}
	}

	return chosenFieldValue
}

func getPreferences(userId string, userName string, visualizations []Visualization, wg *sync.WaitGroup) {
	defer wg.Done()

	genreCount := make(map[string]int)
	typeCount := make(map[string]int)

	for _, visualization := range visualizations {
		genreCount[visualization.Genre]++
		typeCount[visualization.Type]++
	}

	chosenGenre := getPreferred(genreCount)
	chosenType := getPreferred(typeCount)

	preference := Preference{
		UserID:          userId,
		UserName:        userName,
		ChosenGenre:     chosenGenre,
		ChosenType:      chosenType,
		Total:           len(visualizations),
		DifferentGenres: len(genreCount),
	}

	writePreferenceToFile(userId, preference)

}

func writePreferenceToFile(userId string, preference Preference) {
	root := getRootDir()
	outFileName := fmt.Sprintf("user_%s_preferences.csv", userId)
	outFilePath := filepath.Join(root, outFileName)
	outputFile, err := os.Create(outFilePath)
	if err != nil {
		fmt.Printf("Error creating preferences file for user %s: %v", userId, err)
		return
	}

	defer outputFile.Close()

	encoder := json.NewEncoder(outputFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(preference); err != nil {
		fmt.Printf("Error writing preferences to file for user %s: %v", userId, err)
		return
	}
	fmt.Printf("Preferences for user %s written to %s\n", userId, outFileName)
}

func mergePreferencesFiles(visualizations map[string][]Visualization, root string) {
	mergedFilePath := filepath.Join(root, "preferencias.json")
	mergedFile, err := os.Create(mergedFilePath)
	if err != nil {
		fmt.Printf("Error creating merged preferences file: %v", err)
		return
	}
	defer mergedFile.Close()

	var allPreferences []Preference

	for userId := range visualizations {
		userFileName := fmt.Sprintf("user_%s_preferences.csv", userId)
		userFilePath := filepath.Join(root, userFileName)
		userFile, err := os.Open(userFilePath)
		if err != nil {
			fmt.Printf("Error opening user preferences file %s: %v", userFileName, err)
			continue
		}
		defer userFile.Close()

		var preference Preference

		if err := json.NewDecoder(userFile).Decode(&preference); err != nil {
			fmt.Printf("Error decoding user preferences file %s: %v", userFileName, err)
			userFile.Close()
			continue
		}
		userFile.Close()

		if err := os.Remove(userFilePath); err != nil {
			fmt.Printf("Error deleting user preferences file %s: %v\n", userFileName, err)
			continue
		}

		allPreferences = append(allPreferences, preference)
		fmt.Printf("Loaded preferences for user %s\n", userId)
	}

	encoder := json.NewEncoder(mergedFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allPreferences); err != nil {
		fmt.Printf("Error writing merged preferences: %v\n", err)
		return
	}

	fmt.Printf("Merged preferences written to %s\n", mergedFilePath)
}

func main() {

	root := getRootDir()
	filePath := filepath.Join(root, PreferencesFileName)

	visualizations := getVisualizations(filePath)

	var wg sync.WaitGroup

	for userId, visualizations := range visualizations {
		wg.Add(1)
		go getPreferences(userId, visualizations[0].UserName, visualizations, &wg)
	}

	wg.Wait()

	mergePreferencesFiles(visualizations, root)

}
