package main

import (
	"fmt"
	"go-SchemaRestifier/internal/generator"
	"go-SchemaRestifier/internal/parser"
)

func main() {

	data, err := parser.ParseSchema("testdata/")

	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(data)

	err = generator.GeneratorMain("testdata", data)

}
