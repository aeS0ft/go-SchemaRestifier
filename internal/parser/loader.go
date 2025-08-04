package parser

import (
	"go-SchemaRestifier/internal/util"
)

// readSchemaFile LoadSchemaFromFile loads a schema from a JSON file and returns a interface object.
func readSchemaFile(filePath string) (interface{}, error) {
	var schema interface{}
	err := util.ReadJSONFile(filePath, &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// LoadSchemasDirectory loads all JSON schema files from a specified directory and returns them as a map.
// The keys of the map are the file names, and the values are the parsed schema objects.
// It returns an error if there is an issue reading the directory or parsing any of the files
func LoadSchemasDirectory(directoryPath string) (map[string]interface{}, error) {
	schemas := make(map[string]interface{})

	files, err := util.ListFilesInDirectory(directoryPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if util.IsJSONFile(file) {
			schema, err := readSchemaFile(file)
			if err != nil {
				return nil, err
			}
			schemas[file] = schema
		}
	}

	return schemas, nil
}
