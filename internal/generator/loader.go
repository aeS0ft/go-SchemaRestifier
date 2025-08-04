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

// writeFile writes the given content to a file at the specified path, creating the file if it does not exist.
func writeFile(filePath string, content []byte) error {
	return util.WriteFile(filePath, content)
}

func GetModuleRoot() string {
	// This function should return the root directory of the module.
	// For simplicity, we can assume it returns a hardcoded path.
	// In a real application, you might want to determine this dynamically.
	return "go-SchemaRestifier"
}
