package generator

import (
	"go-SchemaRestifier/internal/util"
)

// checkfileExists checks if a file exists at the given path.
func checkFileExists(filePath string) (bool, error) {
	return util.CheckFile(filePath)
}

// readFile reads the content of a file and returns it as a byte slice.
func readFile(filePath string) ([]byte, error) {
	content, err := util.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return content, nil
}
