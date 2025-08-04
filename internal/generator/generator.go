package generator

import (
	"fmt"
	"go-SchemaRestifier/internal/parser"
	"os"
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

// GenarateModel generates the model layer with all the structs from the schema
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
	for _, schema := range content {
		modelContent += fmt.Sprintf("type %s struct {\n", schema.Name)
		for _, column := range *schema.Columns {
			s := fmt.Sprintf("%s", column.Type)
			a, _ := ParseTypes(s)
			modelContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", column.Name, a.String(), column.Name)
			if len(column.Nestedcolumns) > 0 {

				for nestedName, nestedColumns := range column.Nestedcolumns {
					for _, nestedColumn := range nestedColumns {
						s := fmt.Sprintf("%s", nestedColumn)
						nestedtype, _ := ParseTypes(s)
						modelContent += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", nestedName, nestedtype.String(), nestedName)
						if nestedColumn == "PrimaryKey" {
							modelContent += fmt.Sprintf("\t// Nested column is a primary key\n")
						} else {
							modelContent += fmt.Sprintf("\t// Nested column: %s\n", nestedColumn)
						}

					}

					modelContent += fmt.Sprintf("\t// Nested columns: %v\n", nestedColumns)
				}
				modelContent += fmt.Sprintf("\t// Nested columns: %v\n", column.Nestedcolumns)
			}
		}
		modelContent += "}\n\n"
	}
	return nil
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
