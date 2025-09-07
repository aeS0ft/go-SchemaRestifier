package generator

import (
	"fmt"
	"go-SchemaRestifier/internal/datastructures"
	"go-SchemaRestifier/internal/parser"
	"os"

	"github.com/iancoleman/strcase"
)

type TableNameList []string

// GeneratorMain is the main function for generating the Go code based on the parsed schema.
// It takes the file path where the generated code will be saved and the content which is a slice of parser.Schema.
func GeneratorMain(filePath string, content []parser.Schema) error {
	err := GenerateRunner(filePath+"/runner.go", content)
	if err != nil {
		return err
	}
	err = GenerateModel(filePath+"model/", content)
	if err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}
	return nil

	//TODO: Create function to generate the API controller layer which uses a object instance to make controllers for the api routes. Using Mux router.

	// TODO: Create a function to generate the logic/service layer which uses the crud values from the schema,
	// the types to be used for validation i.e length of strings.
	// nullability of fields, uniqueness etc.

	//TODO: Database layer using sqlx with crud operations. and potentially automatic middleware for authentication. for postgresql.

}

func GenerateRunner(filePath string, content []parser.Schema) error {
	exists, err := checkFileExists(filePath)
	tableNames := make(TableNameList, 0, len(content))
	for _, value := range content {
		tableNames = append(tableNames, value.Name)
	}
	contents := "package main \n\nimport (\n\t\"fmt\"\n\"net/http\"\n)\n\nfunc main() {\nmux := http.NewServeMux()\n"
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", filePath, err)
	}
	if exists {
		// TODO: add handling for existing file
		fmt.Printf("File %s already exists, skipping generation.\n", filePath)
		return nil
	}
	for _, tableName := range tableNames {
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

func GenerateDTO(filepath string, content []parser.Schema) error {

	return nil
}

// GenerateModel generates the model layer with all the structs from the schema
func GenerateModel(filePath string, content []parser.Schema) error {

	// TODO: add handling for existing file
	for _, schema := range content {
		modelContent := "package model\n\n"
		modelContent += fmt.Sprintf("type %s struct {\n", strcase.ToCamel(schema.Name))
		nestedStructContent := ""

		for _, column := range *schema.Columns {
			s := fmt.Sprintf("%s", column.Type)
			a, _ := ParseTypes(s)
			if s != "json" {
				modelContent += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", strcase.ToCamel(column.Name), a.String(), column.Name)
			} else {
				modelContent += fmt.Sprintf("\t%s %sOBJ `json:\"%s\"`\n", strcase.ToCamel(column.Name), strcase.ToCamel(column.Name), column.Name)

			}
			if column.Nestedcolumns != nil {
				// Your logic for handling nested columns
				p := new(string)
				traTree, _ := TraverseTreeModel(column.Nestedcolumns, p)
				nestedStructContent += traTree
			}
		}
		modelContent += "}\n\n"
		modelContent += nestedStructContent

		err := writeFile(filePath+""+schema.Name+".go", []byte(modelContent))
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}
	return nil
}

// ConvertNestedMapToNodeTree processes a nested map structure and converts it into a tree-like structure using the datastructures.Node type.
// It recursively traverses the map, creating nodes for each key-value pair and handling nested maps appropriately.
func TraverseTreeModel(n *datastructures.Node, p *string) (string, error) {
	if p == nil {

		p = new(string)

	}
	if n == nil {
		return "", fmt.Errorf("node is nil")
	} else {
		*p += fmt.Sprintf("type %sOBJ struct {\n", strcase.ToCamel(n.Name))
	}
	fmt.Println(n.Name)

	for _, field := range n.Fields {
		*p += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", strcase.ToCamel(field.Name), field.Type, field.Name)
	}
	// needed to build the struct with the fields for each child so encapsulation works
	for _, child := range n.Children {
		*p += fmt.Sprintf("\t%s %sOBJ `json:\"%s\"`\n", strcase.ToCamel(child.Name), strcase.ToCamel(child.Name), child.Name)

	}
	*p += "}\n\n"

	for _, child := range n.Children {

		pString, err := TraverseTreeModel(child, new(string))
		if err != nil {
			return "", err
		}
		*p += pString

	}

	return *p, nil
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
