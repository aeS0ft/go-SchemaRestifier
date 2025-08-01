package generator

import (
	"fmt"
	"go-SchemaRestifier/internal/parser"
	"os"
)

type TableNameList []string

func generator(filePath string, content []parser.Schema) error {
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
	return nil

}

func GenerateRunner(filePath string, content []string) error {
	exists, err := checkFileExists(filePath)
	contents := "package main \n\nimport (\n\t\"fmt\"\n\"net/http\")\n\nfunc main() {\nmux := http.NewServeMux()\n"
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
		contents += fmt.Sprintf("\tcontrollers.%s.RegisterRoutes(mux)\n}", tableName)

	}
	contents += "\n\tfmt.Println(\"Server is running on port 8080\")\n\thttp.ListenAndServe(\":8080\", mux)\n}\n"
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
