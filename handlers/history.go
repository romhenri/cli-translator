package history

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

var historyFile = "history.json"
var mutex sync.Mutex

type TranslationHistory struct {
	Entries []TranslationEntry `json:"entries"`
}

type TranslationEntry struct {
	Text        string `json:"text"`
	FromLang    string `json:"from_lang"`
	TargetLang  string `json:"target_lang"`
	Translation string `json:"translation"`
}

// Load
func loadHistory() (TranslationHistory, error) {
	var history TranslationHistory

	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		history.Entries = []TranslationEntry{}
		return history, nil
	}

	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.ReadFile(historyFile)
	if err != nil {
		return history, fmt.Errorf("Error to read: %w", err)
	}

	err = json.Unmarshal(file, &history)
	if err != nil {
		return history, fmt.Errorf("Error to decode: %w", err)
	}

	return history, nil
}

// Save
func SaveToHistory(text, fromLang, targetLang, translation string) {
	history, err := loadHistory()
	if err != nil {
		fmt.Println("Error to load:", err)
		return
	}

	newEntry := TranslationEntry{
		Text:        text,
		FromLang:    fromLang,
		TargetLang:  targetLang,
		Translation: translation,
	}
	history.Entries = append(history.Entries, newEntry)

	mutex.Lock()
	defer mutex.Unlock()

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		fmt.Println("Error to convert:", err)
		return
	}

	err = os.WriteFile(historyFile, data, 0644)
	if err != nil {
		fmt.Println("Error to save", err)
	}
}

// Show
func ShowHistory() {
	history, err := loadHistory()
	if err != nil {
		fmt.Println("Error to load:", err)
		return
	}

	if len(history.Entries) == 0 {
		fmt.Println("No translations found in history.")
		return
	}

	fmt.Println("\n= Translations =")
	for _, entry := range history.Entries {
		fmt.Printf("\nText: %s\nFrom: %s\nTarget: %s\nResult: %s\n----------------\n",
			entry.Text, entry.FromLang, entry.TargetLang, entry.Translation)
	}
}

// Clear
func ClearHistory() {
	mutex.Lock()
	defer mutex.Unlock()

	err := os.Remove(historyFile)
	if err != nil {
		fmt.Println("Error to clear:", err)
	} else {
		fmt.Println("History cleared.")
	}
}
