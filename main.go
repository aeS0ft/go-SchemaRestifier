package main

import (
	"fmt"
	"go-SchemaRestifier/internal/parser"
)

func main() {
	data, err := parser.ParseSchema("testdata/data.json")

	if err != nil {
		fmt.Println("Error:", err)

		return
	}

	fmt.Println(data)
}
