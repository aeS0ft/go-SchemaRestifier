package parser

import (
	"fmt"
)

func ParseSchema(schemaFilePath string) ([]Schema, error) {
	fmt.Println("Parsing schema...")

	schemas, err := LoadSchemasDirectory(schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load schemas from directory %s: %w", schemaFilePath, err)
	}

	schemaList := make([]Schema, len(schemas))
	idx := 0

	for i, schemaData := range schemas {
		idx++
		fmt.Printf("Schema %d: %v\n", i, schemaData)
		schema := Schema{}

		for key, value := range schemaData.(map[string]interface{}) {
			fmt.Printf("Key: %s, Value: %v\n", key, value)
			if key == "table" {
				for key, value := range value.(map[string]interface{}) {
					switch key {
					case "name":
						schema.Name = value.(string)
					case "columns":
						columns := value.([]interface{})
						schema.Columns = &[]Column{}
						for _, col := range columns {
							if colMap, ok := col.(map[string]interface{}); ok {
								name, _ := colMap["name"].(string)
								typ, _ := colMap["type"].(string)
								desc, _ := colMap["description"].(string)
								column := Column{
									Name:        name,
									Type:        typ,
									Description: desc,
								}
								if pk, ok := colMap["primary_key"]; ok && pk.(bool) {
									column.PrimaryKey = true
								} else {
									column.PrimaryKey = false
								}
								if nested, ok := colMap["json_data"]; ok {
									if nestedMap, ok := nested.(map[string]interface{}); ok {
										convertedMap := make(map[string][]interface{})
										for k, v := range nestedMap {
											convertedMap[k] = []interface{}{v}
										}
										column.Nestedcolumns = convertedMap
									} else {
										fmt.Println("Invalid nested columns format:", nested)
									}
								}
								*schema.Columns = append(*schema.Columns, column)
							} else {
								fmt.Println("Invalid column format:", col)
							}
						}
					default:
						fmt.Printf("Unknown table key: %s with value: %v\n", key, value)
					}
				}
			} else if key == "crud" {
				schema.Crud = value.(map[string]interface{})
			} else {
				fmt.Printf("Unknown key: %s with value: %v\n", key, value)
			}

		}
		fmt.Println("Schema parsed successfully.")
		schemaList[idx-1] = schema
	}

	return schemaList, nil
}
