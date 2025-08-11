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

	const outputDir = "output/"
	err = generator.GeneratorMain(outputDir, data)

}
