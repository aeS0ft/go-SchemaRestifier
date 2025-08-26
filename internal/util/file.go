package util

import (
	"fmt"
	"io"
	"os"
)

// ReadFile reads the content of a file to check if an http server has already been made and concretes if so and returns it as a byte slice.
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file %s: %v\n", filePath, err)
		}
	}(file)

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return content, nil
}

// WriteFile writes the given content to a file at the specified path, creating the file if it does not exist.
func WriteFile(filePath string, content []byte) error {

	file, err := os.Create(filePath)
	newfilepath := StripGOFileFromPath(filePath)
	println(newfilepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("failed to close file %s: %v\n", filePath, err)
		}
	}(file)

	_, err = file.Write(content)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}

	return nil
}

// CheckFile checks if a file exists at the given path.
func CheckFile(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil // File does not exist
	} else if err != nil {
		return false, fmt.Errorf("failed to check file %s: %w", filePath, err)
	}
	return true, nil // File exists
}
func ListFilesInDirectory(directoryPath string) ([]string, error) {
	files, err := os.ReadDir(directoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %s: %w", directoryPath, err)
	}

	var filePaths []string
	for _, file := range files {
		if !file.IsDir() {
			filePaths = append(filePaths, directoryPath+"/"+file.Name())
		}
	}

	return filePaths, nil
}

func StripGOFileFromPath(filePath string) string {
	result := ""
	runes := []rune(filePath)
	for i := len(runes) - 1; i >= 0; i-- {
		if filePath[i] == '/' {
			result = string(runes[0:i])
			break
		}
	}
	return result
}
