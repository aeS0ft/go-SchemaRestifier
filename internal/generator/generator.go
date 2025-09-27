package generator

import (
	"fmt"
	"go-SchemaRestifier/internal/datastructures"
	"go-SchemaRestifier/internal/parser"
	"os"
	"strings"

	"github.com/iancoleman/strcase"
)

type Generator struct {
	Schemas    []parser.Schema
	APIName    string
	FilePath   string
	TableNames []string
}

// GeneratorMain is the main function for generating the Go code based on the parsed schema.
// It takes the file path where the generated code will be saved and the content which is a slice of parser.Schema.
func (g *Generator) GeneratorMain() error {
	err := g.GenerateRunner()
	if err != nil {
		return err
	}
	err = g.GenerateModel()
	if err != nil {
		return fmt.Errorf("failed to generate model: %w", err)
	}
	err = g.GenerateDTO()
	if err != nil {
		return fmt.Errorf("failed to generate dto: %w", err)
	}
	err = g.GenerateGoMod()
	if err != nil {
		return fmt.Errorf("failed to generate go.mod: %w", err)
	}
	return nil

	//TODO: Create function to generate the API controller layer which uses a object instance to make controllers for the api routes. Using Mux router.

	// TODO: Create a function to generate the logic/service layer which uses the crud values from the schema,
	// the types to be used for validation i.e length of strings.
	// nullability of fields, uniqueness etc.

	//TODO: Database layer using sqlx with crud operations. and potentially automatic middleware for authentication. for postgresql.

}

func (g *Generator) GenerateRunner() error {
	exists, err := checkFileExists(g.FilePath + "/runner.go")
	filePath := g.FilePath + "/runner.go"

	import_content := "package main \n\nimport (\n\t\"fmt\"\n\"net/http\"\n"
	contents := "\n\nfunc main() {\n\n\tmux := http.NewServeMux()\n"
	if err != nil {
		return fmt.Errorf("failed to check if file exists %s: %w", filePath, err)
	}
	if exists {
		// TODO: add handling for existing file
		fmt.Printf("File %s already exists, skipping generation.\n", filePath)
		return nil
	}
	for _, tableName := range g.TableNames {
		fmt.Printf("Generating runner/api-routing for table: %s\n", tableName)
		import_content += fmt.Sprintf("\t\"%s/controller/%s\"\n", g.APIName, tableName)
		contents += fmt.Sprintf("\n\t%s.RegisterRoutes(mux)\n", tableName)

	}
	import_content += ")\n"
	contents = import_content + contents
	contents += "\n\tfmt.Println(\"Server is running on port 8080\")\n\thttp.ListenAndServe(\":8080\", mux)\n}\n"

	err = writeFile(filePath, []byte(contents))
	if err != nil {
		return err
	}

	return nil
}

func (g *Generator) GenerateGoMod() error {
	modContent := fmt.Sprintf("module %s\n\ngo 1.20\n\nrequire (\n\tgithub.com/gorilla/mux v1.8.0\n)\n", g.APIName) // Replace %s with the module name
	exists, err := checkFileExists(g.FilePath + "/go.mod")
	filepath := g.FilePath + "/go.mod"
	if err != nil {
		return fmt.Errorf("failed to check if go.mod file exists: %w", err)
	}
	if exists {
		// TODO: add handling for existing file
		// For now, we will just skip generation if the file already exists.
		fmt.Printf("go.mod file already exists at %s, skipping generation.\n", filepath)
		return nil
	}

	err = writeFile(filepath, []byte(modContent))
	if err != nil {
		return fmt.Errorf("failed to write go.mod file: %w", err)
	}
	return nil
}

func (g *Generator) GenerateDTO() error {

	for _, schema := range g.Schemas {
		depenencies := new(string)
		*depenencies = "package dto\n\nimport (\n"
		modelContent := fmt.Sprintf("type %s struct {\n", strcase.ToCamel(schema.Name))
		nestedStructContent := ""
		for _, column := range *schema.Columns {
			if column.Hidden {
				continue
			} else {
				s := fmt.Sprintf("%s", column.Type)
				a, _ := ParseTypes(s)
				if s != "json" {
					modelContent += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", strcase.ToCamel(column.Name), a.String(), column.Name)
					switch a.String() {
					case "time.Time":
						if !strings.Contains(*depenencies, "\"time\"") {
							*depenencies += "\"time\"\n"
						}
					}
				} else {
					modelContent += fmt.Sprintf("\t%s %s_%sOBJ `json:\"%s\"`\n", strcase.ToCamel(column.Name), strcase.ToCamel(schema.Name), strcase.ToCamel(column.Name), column.Name)
				}
				if column.Nestedcolumns != nil {
					// Your logic for handling nested columns
					result, _ := TraverseTree(
						column.Nestedcolumns,
						nil,
						func(field datastructures.Field) string {
							if field.Hidden {
								return ""
							}
							return fmt.Sprintf("\t%s %s `json:\"%s\"`\n", strcase.ToCamel(field.Name), field.Type, field.Name)
						},
						func(child *datastructures.Node) string {
							if child.Hidden {
								return ""
							}
							return fmt.Sprintf("\t%s %s_%sOBJ `json:\"%s\"`\n",
								strcase.ToCamel(child.Name),
								strcase.ToCamel(schema.Name),
								strcase.ToCamel(child.Name),
								child.Name)
						},
						func(node *datastructures.Node) (Value string, isAborted bool) {

							if node.Hidden {
								return "", true
							}
							return fmt.Sprintf("type %s_%sOBJ struct {\n",
								strcase.ToCamel(schema.Name),
								strcase.ToCamel(node.Name)), false
						},
						func(name string) string {
							return strcase.ToCamel(name)
						},
						"json",
					)
					*depenencies, _ = TraverseTree(
						column.Nestedcolumns,
						depenencies,
						func(field datastructures.Field) string {
							switch field.Type {
							//add more dependencies here
							case "time.Time":
								if !strings.Contains(*depenencies, "\"time\"") {
									return "\"time\"\n"
								}

							}
							return ""
						},
						func(child *datastructures.Node) string {
							return ""
						},
						func(node *datastructures.Node) (Value string, isAborted bool) {
							return "", false
						},
						func(name string) string {
							return ""
						},
						"dependencies",
					)
					nestedStructContent += result
				}
			}

		}
		modelContent += "}\n\n"
		modelContent += nestedStructContent
		*depenencies += ")\n\n"
		modelContent = *depenencies + modelContent

		err := writeFile(g.FilePath+"/"+"dto"+"/"+schema.Name+".go", []byte(modelContent))
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", g.FilePath+""+schema.Name+".go", err)
		}
	}
	return nil
}

// GenerateModel generates the model layer with all the structs from the schema
func (g *Generator) GenerateModel() error {

	// TODO: add handling for existing file
	for _, schema := range g.Schemas {
		depenencies := new(string)
		*depenencies = "package model\n\nimport (\n"
		modelContent := fmt.Sprintf("type %s struct {\n", strcase.ToCamel(schema.Name))
		nestedStructContent := ""

		for _, column := range *schema.Columns {
			s := fmt.Sprintf("%s", column.Type)
			a, _ := ParseTypes(s)
			if s != "json" {
				modelContent += fmt.Sprintf("\t%s %s `db:\"%s\"`\n", strcase.ToCamel(column.Name), a.String(), column.Name)
				switch a.String() {
				case "time.Time":
					if !strings.Contains(*depenencies, "\"time\"") {
						*depenencies += "\"time\"\n"
					}
				}
			} else {
				modelContent += fmt.Sprintf("\t%s %s_%sOBJ `json:\"%s\"`\n", strcase.ToCamel(column.Name), strcase.ToCamel(schema.Name), strcase.ToCamel(column.Name), column.Name)

			}
			if column.Nestedcolumns != nil {
				// Your logic for handling nested columns
				result, _ := TraverseTree(
					column.Nestedcolumns,
					nil,
					func(field datastructures.Field) string {
						return fmt.Sprintf("\t%s %s `json:\"%s\"`\n", strcase.ToCamel(field.Name), field.Type, field.Name)
					},
					func(child *datastructures.Node) string {
						return fmt.Sprintf("\t%s %s_%sOBJ `json:\"%s\"`\n",
							strcase.ToCamel(child.Name),
							strcase.ToCamel(schema.Name),
							strcase.ToCamel(child.Name),
							child.Name)
					},
					func(node *datastructures.Node) (Value string, isAborted bool) {
						return fmt.Sprintf("type %s_%sOBJ struct {\n",
							strcase.ToCamel(schema.Name),
							strcase.ToCamel(node.Name)), false
					},
					func(name string) string {
						return strcase.ToCamel(name)
					},
					"json",
				)
				// call to build the import dependencies for nested structs

				*depenencies, _ = TraverseTree(
					column.Nestedcolumns,
					depenencies,
					func(field datastructures.Field) string {
						switch field.Type {
						//add more dependencies here
						case "time.Time":
							if !strings.Contains(*depenencies, "\"time\"") {
								return "\"time\"\n"
							}

						}
						return ""
					},
					func(child *datastructures.Node) string {
						return ""
					},
					func(node *datastructures.Node) (Value string, isAborted bool) {
						return "", false
					},
					func(name string) string {
						return ""
					},
					"dependencies",
				)

				nestedStructContent += result

			}

		}
		modelContent += "}\n\n"
		modelContent += nestedStructContent
		*depenencies += ")\n\n"
		modelContent = *depenencies + modelContent

		err := writeFile(g.FilePath+"/"+"model"+"/"+schema.Name+".go", []byte(modelContent))
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", g.FilePath+""+schema.Name+".go", err)
		}
	}
	return nil
}

// TraverseTree traverses the tree and generates code based on the formatters provided.
func TraverseTree(
	n *datastructures.Node,
	p *string,
	formatField func(field datastructures.Field) string,
	formatChild func(child *datastructures.Node) string,
	formatNode func(node *datastructures.Node) (Value string, isAborted bool),
	typeName func(name string) string,
	tag string,
) (string, error) {
	if n == nil {
		return "", fmt.Errorf("node is nil")
	}
	if p == nil {
		p = new(string)
	}

	// check if the formatNode sends an abort signal and if so, return early. implemented in the call.
	nodeStr, isAborted := formatNode(n)
	if isAborted {
		return "", nil
	}
	*p += nodeStr
	for _, field := range n.Fields {
		*p += formatField(*field)
	}
	for _, child := range n.Children {
		*p += formatChild(child)
	}
	// TODO: make a non implementation specific function for this
	if tag == "json" {
		*p += "}\n\n"
	}
	for _, child := range n.Children {
		pString, err := TraverseTree(child, new(string), formatField, formatChild, formatNode, typeName, tag)
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
