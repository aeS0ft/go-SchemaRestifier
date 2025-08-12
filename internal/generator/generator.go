package generator

import (
	"fmt"
	"go-SchemaRestifier/internal/datastructures"
	"go-SchemaRestifier/internal/parser"
	"os"
	"sync"
)

type TableNameList []string

// GeneratorMain is the main function for generating the Go code based on the parsed schema.
// It takes the file path where the generated code will be saved and the content which is a slice of parser.Schema.
func GeneratorMain(filePath string, content []parser.Schema) error {
	tableNames := make(TableNameList, 0, len(content))

	for key, value := range content {
		tableNames = append(tableNames, value.Name)
		fmt.Printf("Key: %s, Value: %v\n", key, value)
	}
	runnerpath := filePath + "/runner.go"
	err := GenerateRunner(runnerpath, tableNames)
	if err != nil {
		return err
	}
	err = GenerateModel(filePath+"/model", content)
	if err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}
	return nil

}

func GenerateRunner(filePath string, content []string) error {
	exists, err := checkFileExists(filePath)
	contents := "package main \n\nimport (\n\t\"fmt\"\n\"net/http\"\n)\n\nfunc main() {\nmux := http.NewServeMux()\n"
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", filePath, err)
	}
	if exists {
		// TODO: add handling for existing file
		fmt.Printf("File %s already exists, skipping generation.\n", filePath)
		return nil
	}
	for _, tableName := range content {
		fmt.Printf("Generating runner/api-routing for table: %s\n", tableName)
		contents += fmt.Sprintf("\t%s.RegisterRoutes(mux)\n}", tableName)

	}
	contents += "\n\tfmt.Println(\"Server is running on port 8080\")\n\thttp.ListenAndServe(\":8080\", mux)\n}\n"

	return nil
}

func GenerateGoMod(filePath string, name string) error {
	modContent := fmt.Sprintf("module %s\n\ngo 1.20\n\nrequire (\n\tgithub.com/gorilla/mux v1.8.0\n)\n", name) // Replace %s with the module name
	exists, err := checkFileExists(filePath + "/go.mod")
	if err != nil {
		return fmt.Errorf("failed to check if go.mod file exists: %w", err)
	}
	if exists {
		// TODO: add handling for existing file
		// For now, we will just skip generation if the file already exists.
		fmt.Printf("go.mod file already exists at %s, skipping generation.\n", filePath+"/go.mod")
		return nil
	}

	err = writeFile(filePath+"/go.mod", []byte(modContent))
	if err != nil {
		return fmt.Errorf("failed to write go.mod file: %w", err)
	}
	return nil
}

// GenerateModel generates the model layer with all the structs from the schema
func GenerateModel(filePath string, content []parser.Schema) error {
	exists, err := checkFileExists(filePath)
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", filePath, err)
	}
	// TODO: add handling for existing file
	if exists {
		return fmt.Errorf("file %s already exists, skipping generation.", filePath)
	}
	modelContent := "package model\n\n"
	nestedStructContent := ""
	for _, schema := range content {
		modelContent += fmt.Sprintf("type %s struct {\n", schema.Name)
		for _, column := range *schema.Columns {
			s := fmt.Sprintf("%s", column.Type)
			a, _ := ParseTypes(s)
			if s != "json" {
				modelContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", column.Name, a.String(), column.Name)
			} else {
				modelContent += fmt.Sprintf("\t%s %sOBJ\n", column.Name, column.Name)

			}
			if len(column.Nestedcolumns) > 0 {
				p := new(string)
				nestedMap := make(map[string]interface{})
				for k, v := range column.Nestedcolumns {
					for _, b := range v {
						nestedMap[k] = b
					}

				}
				nestedMap["@__root_name__@"] = column.Name
				nestedthing, _ := NestedMapAlgo(nestedMap, &datastructures.Node{}, nil)
				fmt.Println(nestedthing)

				nestedStructContent = fmt.Sprintf("type %sOBJ struct {\n", column.Name)

				for nestedName, nestedColumns := range column.Nestedcolumns {
					for _, nestedColumna := range nestedColumns {
						// if statement to check if the column is a nested column ie. a interface or map[string]interface{} and create a new struct

						if _, ok := nestedColumna.(map[string]interface{}); ok {
							new_nestedObj := fmt.Sprintf("type %sOBJ struct {\n", nestedName)
							for key, nestedColumnd := range nestedColumna.(map[string]interface{}) {
								if a, ok := nestedColumnd.(map[string]interface{}); ok {
									a["@__root_name__@"] = key
									go func() {
										_, err := NestedMapAlgo(nestedColumnd.(map[string]interface{}), &datastructures.Node{}, nil)
										if err != nil {
											fmt.Println(err)
										} else {
											fmt.Println("Nested map algo reached base case")
										}
									}()
								}

								s := fmt.Sprintf("%s", nestedColumnd)
								nestedtype, _ := ParseTypes(s)
								new_nestedObj += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", key, nestedtype.String(), nestedName)
								*p += new_nestedObj
								new_nestedObj = ""

							}
							*p += "}\n"

						}

						if _, ok := nestedColumna.(map[string]interface{}); !ok {
							s := fmt.Sprintf("%s", nestedColumna)
							nestedtype, _ := ParseTypes(s)
							nestedStructContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", nestedName, nestedtype.String(), nestedName)
						} else {
							nestedStructContent += fmt.Sprintf("\t%s %sOBJ `json:\"%s\"`\n", nestedName, nestedName, nestedName)
						}

					}

				}
				nestedStructContent += "}\n"
				nestedStructContent += *p
			}

		}
		modelContent += "}\n\n"
		modelContent += nestedStructContent

		err = writeFile(filePath+""+schema.Name+".go", []byte(modelContent))
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}
	return nil
}

// NestedMapAlgo processes a nested map structure and converts it into a tree-like structure using the datastructures.Node type.
// It recursively traverses the map, creating nodes for each key-value pair and handling nested maps appropriately.
func NestedMapAlgo(content map[string]interface{}, n *datastructures.Node, root *datastructures.Node) (datastructures.Node, error) {
	var wg sync.WaitGroup
	if len(content) == 0 {
		return datastructures.Node{}, fmt.Errorf("empty map")
	}

	// Initialize the root node of the tree structure
	if root == nil {
		root = new(datastructures.Node)
		root.Name = content["@__root_name__@"].(string)

	}
	if datastructures.IsNodeEmpty(*n) {
		*n = *root
		root = n

	}

	for key, value := range content {
		// it is the root denoter, so we skip it
		if key == "@__root_name__@" {
			continue
		}
		if _, ok := value.(map[string]interface{}); !ok {

			parsedType, _ := ParseTypes(value.(string))
			typeName := parsedType.String()
			field := datastructures.Fields{
				Name: key,
				Type: typeName,
			}
			n.Fields = append(n.Fields, &field)

			fmt.Println("Base case reached with value:", value)
		} else {
			fmt.Println("Nested map reached with value:", value)
			newN := new(datastructures.Node)
			newN.Name = key
			n.Children = append(n.Children, newN)
			wg.Add(1)
			go func(key string, value interface{}) {
				defer wg.Done()
				_, _ = NestedMapAlgo(value.(map[string]interface{}), newN, root)

			}(key, value)

		}
	}
	wg.Wait()
	return *root, nil
}
func GenerateAPIController(filePath string, content []byte) error {
	file, err := os.Create(filePath)
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
