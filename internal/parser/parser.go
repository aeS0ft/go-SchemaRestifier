package parser

import (
	"fmt"
	"go-SchemaRestifier/internal/datastructures"
	"sync"
)

func JsonDataAlgorithm(content map[string]interface{}, n *datastructures.Node, root *datastructures.Node) (datastructures.Node, error) {
	var wg sync.WaitGroup
	if datastructures.IsNodeEmpty(n) {
		n.Name = root.Name
		root = n
	}
	if len(content) == 0 {
		return datastructures.Node{}, fmt.Errorf("empty map")
	}

	for key, value := range content {
		if m, ok := value.(map[string]interface{}); ok {
			typ, _ := m["type"].(string)
			switch typ {
			case "object":
				fmt.Println("Nested map reached with value:", value)
				newN := new(datastructures.Node)
				newN.Name = key
				newN.Hidden = value.(map[string]interface{})["hidden"].(bool)
				n.Mu.Lock()
				n.Children = append(n.Children, newN)
				n.Mu.Unlock()

				wg.Add(1)
				go func(key string, value interface{}) {
					defer wg.Done()
					_, _ = JsonDataAlgorithm(value.(map[string]interface{}), newN, root)

				}(key, value)
			default:
				parsedType, _ := ParseTypes(typ)
				typeName := parsedType.String()
				field := datastructures.Field{
					Name: key,
					Type: typeName,
				}
				n.Fields = append(n.Fields, &field)
			}
		}

	}
	wg.Wait()
	return *root, nil
}

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
								pk, _ := colMap["primary_key"].(bool)
								hidden, _ := colMap["hidden"].(bool)
								struct_field, _ := colMap["struct"].(map[string]interface{})
								field_name := struct_field["field_name"].(string)
								var nested *datastructures.Node
								if jsonData, ok := colMap["json_data"].(map[string]interface{}); ok {
									prev_node := new(datastructures.Node)
									prev_node.Name = field_name
									node, _ := JsonDataAlgorithm(jsonData, &datastructures.Node{}, prev_node)
									nested = &node
								} else {
									nested = nil
								}
								queries := make(map[string]bool)
								for queryType, IsTrue := range colMap["query"].(map[string]interface{}) {

									queries[queryType] = IsTrue.(bool)
								}

								column := Column{
									Name:          name,
									Type:          typ,
									Description:   desc,
									PrimaryKey:    pk,
									Hidden:        hidden,
									Nestedcolumns: nested,
									Capabilities:  queries,
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
