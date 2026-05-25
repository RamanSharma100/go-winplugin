package main

import (
	"fmt"

	winplugin "github.com/RamanSharma100/go-winplugin"
)

func main() {
	pluginLoader, err := winplugin.NewLoader(
		"./example",
	)

	if err != nil {
		panic(err)
	}

	err = pluginLoader.Build(
		"plugin.go",
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	output, err := pluginLoader.Call(
		"example",
		"Execute",
		uintptr(10),
		uintptr(20),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Output:", output)
}
