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

const (
	PreferencesFileName = "visualizaciones.csv"
	UserIdColumn        = 0
	UserNameColumn      = 1
	TitleColumn         = 2
	TypeColumn          = 3
	GenreColumn         = 4
	CurrentGoRoutine    = 1
	Error               = 1
	CurrentExecution    = 0
	SkipHeader          = 1
	Currentuser         = 0
)

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
	_, callerFile, _, ok := runtime.Caller(CurrentExecution)
	if !ok {
		fmt.Println("Could not get callerFile")
		os.Exit(Error)
	}
	return filepath.Dir(callerFile)
}

func getVisualizationsMap(records [][]string) map[string][]Visualization {
	visualizations := make(map[string][]Visualization)

	for _, record := range records[SkipHeader:] {
		userId := record[UserIdColumn]
		vis := Visualization{
			UserID:   userId,
			UserName: record[UserNameColumn],
			Title:    record[TitleColumn],
			Type:     record[TypeColumn],
			Genre:    record[GenreColumn],
		}
		visualizations[userId] = append(visualizations[userId], vis)
	}

	return visualizations
}

func readVisualizationsFile(file *os.File) ([][]string, error) {
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	return records, err
}

func getVisualizations(filePath string) map[string][]Visualization {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Could not open file visualizaciones.csv", err)
		return nil
	}
	defer f.Close()
	records, err := readVisualizationsFile(f)

	if err != nil {
		fmt.Println("Error reading CSV:", err)
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

func getPreference(userId string, userName string, chosenGenre string, chosenType string, totalVisualizations int, genreCount map[string]int) Preference {
	preference := Preference{
		UserID:          userId,
		UserName:        userName,
		ChosenGenre:     chosenGenre,
		ChosenType:      chosenType,
		Total:           totalVisualizations,
		DifferentGenres: len(genreCount),
	}
	return preference
}

func createPreferences(userId string, userName string, visualizations []Visualization, wg *sync.WaitGroup) {
	defer wg.Done()

	genreCount := make(map[string]int)
	typeCount := make(map[string]int)

	for _, visualization := range visualizations {
		genreCount[visualization.Genre]++
		typeCount[visualization.Type]++
	}

	chosenGenre := getPreferred(genreCount)
	chosenType := getPreferred(typeCount)
	preference := getPreference(userId, userName, chosenGenre, chosenType, len(visualizations), genreCount)

	writePreferenceToFile(userId, preference)
}

func createOutputFilePath(userId string) (string, *os.File, error) {
	root := getRootDir()
	outFileName := fmt.Sprintf("user_%s_preferences.csv", userId)
	outFilePath := filepath.Join(root, outFileName)
	outputFile, err := os.Create(outFilePath)
	return outFileName, outputFile, err
}

func writePreferenceToFile(userId string, preference Preference) {
	outFileName, outputFile, err := createOutputFilePath(userId)
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

	allPreferences := loadUsersPreferences(visualizations, root)

	encoder := json.NewEncoder(mergedFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(allPreferences); err != nil {
		fmt.Printf("Error writing merged preferences: %v\n", err)
		return
	}

	fmt.Printf("Merged preferences written to %s\n", mergedFilePath)
}

func loadUsersPreferences(visualizations map[string][]Visualization, root string) []Preference {
	var allPreferences []Preference

	for userId := range visualizations {
		preference, err := loadUserPreference(userId, root)
		if err != nil {
			fmt.Printf("Error processing user %s: %v\n", userId, err)
			continue
		}
		allPreferences = append(allPreferences, preference)
		fmt.Printf("Loaded preferences for user %s\n", userId)
	}

	return allPreferences
}

func loadUserPreference(userId string, root string) (Preference, error) {
	var preference Preference
	userFileName, userFilePath, file, err := createUserFilePath(userId, root)

	if err != nil {
		return preference, fmt.Errorf("opening file %s: %w", userFileName, err)
	}

	if err := json.NewDecoder(file).Decode(&preference); err != nil {
		file.Close()
		return preference, fmt.Errorf("decoding file %s: %w", userFileName, err)
	}
	file.Close()

	if err := os.Remove(userFilePath); err != nil {
		fmt.Printf("Warning: failed to delete file %s: %v\n", userFileName, err)
	}
	return preference, nil
}

func createUserFilePath(userId string, root string) (string, string, *os.File, error) {
	userFileName := fmt.Sprintf("user_%s_preferences.csv", userId)
	userFilePath := filepath.Join(root, userFileName)
	file, err := os.Open(userFilePath)
	return userFileName, userFilePath, file, err
}

func main() {

	root := getRootDir()
	filePath := filepath.Join(root, PreferencesFileName)

	visualizations := getVisualizations(filePath)

	var wg sync.WaitGroup

	for userId, visualizations := range visualizations {
		wg.Add(CurrentGoRoutine)
		go createPreferences(userId, visualizations[Currentuser].UserName, visualizations, &wg)
	}

	wg.Wait()

	mergePreferencesFiles(visualizations, root)

}
