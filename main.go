package main

import (
	"fmt"
	"go-SchemaRestifier/internal/generator"
	"go-SchemaRestifier/internal/parser"

	"go.uber.org/dig"
)

type APIName string
type FilePath string

func NewGenerator(apiName APIName, schemas []parser.Schema, Filepath FilePath) *generator.Generator {
	return &generator.Generator{
		APIName:  string(apiName),
		Schemas:  schemas,
		FilePath: string(Filepath),
	}
}

func main() {

	container := dig.New()

	// Provide the parsed schema data to the container
	err := container.Provide(func() ([]parser.Schema, error) {
		return parser.ParseSchema("testdata/")
	})
	if err != nil {
		fmt.Println("Error providing schema data:", err)
		return
	}
	// Provide metadata like api name

	// Provide metadata like api name
	err = container.Provide(func() APIName {
		return "SampleAPI"
	})
	if err != nil {
		fmt.Println("Error providing API name:", err)
		return
	}
	err = container.Provide(func() FilePath {
		return "/output"
	})
	if err != nil {
		fmt.Println("Error providing file path:", err)
		return
	}

	// Invoke the generator with the parsed schema data
	container.Provide(NewGenerator)
	err = container.Invoke(func(gen *generator.Generator) {
		err := gen.GeneratorMain()
		if err != nil {
			return
		}
	})

}
