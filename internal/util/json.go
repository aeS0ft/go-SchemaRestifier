package util

import (
	"encoding/json"
	"fmt"
	"os"
)

// ReadJSONFile reads a JSON file and unmarshals its content into the provided interface.
func ReadJSONFile(filePath string, v interface{}) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open JSON file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close JSON file %s: %v\n", filePath, err)
		}
	}(file)

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON from file %s: %w", filePath, err)
	}

	return nil
}

// IsJSONFile checks if the given file path has a ".json" extension.
// It returns true if the file path ends with ".json", false otherwise.
func IsJSONFile(filePath string) bool {
	return len(filePath) > 5 && filePath[len(filePath)-5:] == ".json"
}
