package parser

import (
	"go-SchemaRestifier/internal/util"
)

// loadSchemaFromUtil LoadSchemaFromFile loads a schema from a JSON file and returns a interface object.
func loadSchemaFromUtil(filePath string) (interface{}, error) {
	var schema interface{}
	err := util.ReadJSONFile(filePath, &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}
