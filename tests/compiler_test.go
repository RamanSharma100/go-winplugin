package tests

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/RamanSharma100/go-winplugin/compiler"
	"github.com/RamanSharma100/go-winplugin/sandbox"
)

func TestMangleSymbol(
	t *testing.T,
) {

	symbol := compiler.MangleSymbol(
		"example/plugin",
		"Execute",
	)

	expected := "go_plugin_example/plugin_Execute"

	if symbol != expected {
		t.Fatalf(
			"expected %s got %s",
			expected,
			symbol,
		)
	}
}

func TestAnalyzeFunctions(
	t *testing.T,
) {

	source := `
package plugin

func Execute(a int, b int) int {
	return a + b
}

func Version() string {
	return "1.0.0"
}

func logger() {}
`

	set := token.NewFileSet()

	file, err := parser.ParseFile(
		set,
		"plugin.go",
		source,
		parser.ParseComments,
	)

	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(
		file,
	)

	expected := map[string]bool{
		"Execute": true,
		"Version": true,
		"logger":  false,
	}

	if len(functions) != len(expected) {
		t.Fatalf(
			"expected %d exported functions got %d",
			len(expected),
			len(functions),
		)
	}

	for _, fn := range functions {
		if !fn.Exported && expected[fn.Name] {
			t.Fatalf(
				"function %s is not exported",
				fn.Name,
			)
		}
	}

	if functions[0].Name != "Execute" {
		t.Fatalf(
			"expected Execute got %s",
			functions[0].Name,
		)
	}

	if functions[1].Name != "Version" {
		t.Fatalf(
			"expected Version got %s",
			functions[1].Name,
		)
	}
}

func TestValidateFunctions(
	t *testing.T,
) {

	functions := []compiler.Function{
		{
			Name:       "Execute",
			Exported:   true,
			ReturnType: "int",
			Params: []compiler.Param{
				{
					Name: "a",
					Type: "int",
				},
			},
		},
	}

	err := compiler.ValidateFunctions(
		functions,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateWrapper_Int(
	t *testing.T,
) {

	functions := []compiler.Function{
		{
			Name:       "Execute",
			Exported:   true,
			ReturnType: "int",
			Params: []compiler.Param{
				{
					Name: "a",
					Type: "int",
				},
				{
					Name: "b",
					Type: "int",
				},
			},
		},
	}

	wrapper := compiler.GenerateWrapper(
		"plugin",
		"example/plugin",
		"plugin",
		functions,
	)

	checks := []string{
		"package main",
		"import \"C\"",
		"plugin \"plugin\"",
		"//export go_plugin_plugin_Execute",
		"func go_plugin_plugin_Execute",
		"result := example/plugin.Execute",
		"return C.int(result)",
	}

	for _, check := range checks {
		if !strings.Contains(
			wrapper,
			check,
		) {
			t.Fatalf(
				"missing generated code: %s\n\n%s",
				check,
				wrapper,
			)
		}
	}
}

func TestGenerateWrapper_String(
	t *testing.T,
) {

	functions := []compiler.Function{
		{
			Name:       "Version",
			Exported:   true,
			ReturnType: "string",
		},
	}

	wrapper := compiler.GenerateWrapper(
		"plugin",
		"example/plugin",
		"plugin",
		functions,
	)

	checks := []string{
		"func go_plugin_plugin_Version",
		"plugin.Version()",
		"C.CString",
	}

	for _, check := range checks {

		if !strings.Contains(
			wrapper,
			check,
		) {
			t.Fatalf(
				"missing generated code: %s",
				check,
			)
		}
	}
}

func TestGenerateWrapper_Void(
	t *testing.T,
) {

	functions := []compiler.Function{
		{
			Name:     "Logger",
			Exported: true,
			Params: []compiler.Param{
				{
					Name: "message",
					Type: "string",
				},
			},
		},
	}

	wrapper := compiler.GenerateWrapper(
		"plugin",
		"example/plugin",
		"plugin",
		functions,
	)

	checks := []string{
		"func go_plugin_plugin_Logger",
		"plugin.Logger",
		"C.GoString",
		"example/plugin.Logger",
	}

	for _, check := range checks {
		if !strings.Contains(
			wrapper,
			check,
		) {
			t.Fatalf(
				"missing generated code: %s",
				check,
			)
		}
	}
}

func TestSandboxValidation(
	t *testing.T,
) {
	source := `
package example

import "os"

func Execute() {}
`

	set := token.NewFileSet()

	file, err := parser.ParseFile(
		set,
		"plugin.go",
		source,
		parser.ParseComments,
	)

	if err != nil {
		t.Fatal(err)
	}

	err = sandbox.ValidateSandbox(
		file,
	)

	if err == nil {
		t.Fatal(
			"expected sandbox violation",
		)
	}
}

func TestStructParameterAnalysis(
	t *testing.T,
) {
	source := `
package example

type User struct {
	Name string
}

func ProcessUser(
	user *User,
) {
}
`

	set := token.NewFileSet()

	file, err := parser.ParseFile(
		set,
		"plugin.go",
		source,
		parser.ParseComments,
	)

	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(
		file,
	)

	if len(functions) != 1 {
		t.Fatalf(
			"expected 1 function got %d",
			len(functions),
		)
	}

	if functions[0].Params[0].Type != "*User" {
		t.Fatalf(
			"expected *User got %s",
			functions[0].Params[0].Type,
		)
	}
}

func TestStructValidation(
	t *testing.T,
) {
	functions := []compiler.Function{
		{
			Name:     "ProcessUser",
			Exported: true,
			Params: []compiler.Param{
				{
					Name: "user",
					Type: "*User",
				},
			},
		},
	}

	err := compiler.ValidateFunctions(
		functions,
	)

	if err != nil {
		t.Fatal(err)
	}
}

func TestMangleSymbolConsistency(
	t *testing.T,
) {
	first := compiler.MangleSymbol(
		"example",
		"Execute",
	)

	second := compiler.MangleSymbol(
		"example",
		"Execute",
	)

	if first != second {
		t.Fatal(
			"symbol generation not deterministic",
		)
	}
}
