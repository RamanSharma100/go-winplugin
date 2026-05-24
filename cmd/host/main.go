package main

import (
	"fmt"

	"github.com/RamanSharma100/go-winplugin/compiler"
)

func main() {
	file, err := compiler.ParseFile(
		"./example/main.go",
	)
	if err != nil {
		panic(err)
	}

	functions := compiler.AnalyzeFunctions(file)

	for _, fn := range functions {
		fmt.Println(fn.Name)
	}
}
