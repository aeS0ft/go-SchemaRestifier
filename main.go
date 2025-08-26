package main

import (
	"fmt"
	"go-SchemaRestifier/internal/generator"
	"go-SchemaRestifier/internal/parser"
	"path/filepath"
)

func main() {

	data, err := parser.ParseSchema("testdata/")

	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(data)

	filedir, _ := filepath.Abs("./")

	const outputDir = "/output/"
	err = generator.GeneratorMain(filedir+outputDir, data)

}
