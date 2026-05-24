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

	err = pluginLoader.Call(
		"plugin",
		"Execute",
	)
	if err != nil {
		panic(err)
	}
}
