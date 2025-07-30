package parser

import (
	"fmt"
)

func ParseSchema(schemaFilePath string) (Schema, error) {
	fmt.Println("Parsing schema...")
	// Here you would implement the logic to parse the schema.
	UtilSchema, err := loadSchemaFromUtil(schemaFilePath)

	if err != nil {
		fmt.Println("Error:", err)
		return Schema{}, fmt.Errorf("failed to load schema from file %s: %w", schemaFilePath, err)
	}
	fmt.Println(UtilSchema)
	// UtilSchema is a dynamic interface{} type that can hold any JSON structure.
	schema := Schema{}
	for key, value := range UtilSchema.(map[string]interface{}) {
		fmt.Printf("Key: %s, Value: %v\n", key, value)
		if key == "table" {

			for key, value := range value.(map[string]interface{}) {
				switch key {
				case "name":
					schema.Name = value.(string)
				case "columns":
					columns := value.([]interface{})
					schema.Fields = &[]Column{} // Initialize the Fields slice
					for _, col := range columns {
						if colMap, ok := col.(map[string]interface{}); ok {
							column := Column{
								Name:        colMap["name"].(string),
								Type:        colMap["type"].(string),
								Description: colMap["description"].(string),
							}
							if pk, ok := colMap["primary_key"]; ok && pk.(bool) {
								column.PrimaryKey = true
							}
							*schema.Fields = append(*schema.Fields, column)
						} else {
							fmt.Println("Invalid column format:", col)
						}

					}
				default:
					fmt.Printf("Unknown table key: %s with value: %v\n", key, value)
				}

			}
		}

		if key == "crud" {
			schema.Crud = value.(map[string]interface{})
		} else {
			fmt.Printf("Unknown key: %s with value: %v\n", key, value)
		}
	}
	// This is a placeholder function to demonstrate the structure.
	// You can call other functions or methods to handle the actual parsing logic.
	fmt.Println("Schema parsed successfully.")
	return schema, nil
}
