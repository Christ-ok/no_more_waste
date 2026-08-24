package i18n

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

var (
	translations = make(map[string]map[string]string)
	mu           sync.RWMutex
)

func Init() error {
	languages := []string{"fr", "en"}

	for _, language := range languages {
		filePath := filepath.Join("i18n", language+".json")

		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("erreur lecture fichier %s : %w", filePath, err)
		}

		var translation map[string]string

		if err := json.Unmarshal(data, &translation); err != nil {
			return fmt.Errorf("erreur parsing fichier %s : %w", filePath, err)
		}

		mu.Lock()
		translations[language] = translation
		mu.Unlock()
	}

	return nil
}

func Traduction(language string, key string) string {
	mu.RLock()
	defer mu.RUnlock()

	if _, exists := translations[language]; !exists {
		language = "fr"
	}

	translation, exists := translations[language][key]

	if !exists {
		translation = translations["fr"][key]
	}

	if translation == "" {
		return key
	}

	return translation
}
