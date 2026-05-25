package tests

import (
	"testing"

	winplugin "github.com/RamanSharma100/go-winplugin"
)

func TestPluginExecution(
	t *testing.T,
) {
	loader, err := winplugin.NewLoader(
		"./fixtures/sample_plugin",
	)
	if err != nil {
		t.Fatal(err)
	}

	err = loader.Build(
		"plugin.go",
	)
	if err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call(
		"sample_plugin",
		"Execute",
		uintptr(10),
		uintptr(20),
	)
	if err != nil {
		t.Fatal(err)
	}

	if result != 30 {
		t.Fatalf(
			"expected 30 got %d",
			result,
		)
	}
}
