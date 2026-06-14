package tests

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/RamanSharma100/go-winplugin/compiler"
)

func TestMangleSymbol(t *testing.T) {
	symbol := compiler.MangleSymbol("example/plugin", "Execute")
	expected := "go_plugin_example/plugin_Execute"
	if symbol != expected {
		t.Fatalf("expected %s got %s", expected, symbol)
	}
}

func TestMangleSymbolConsistency(t *testing.T) {
	first := compiler.MangleSymbol("example", "Execute")
	second := compiler.MangleSymbol("example", "Execute")
	if first != second {
		t.Fatal("symbol generation not deterministic")
	}
}

func TestAnalyzeFunctions(t *testing.T) {
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
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)

	expected := map[string]bool{
		"Execute": true,
		"Version": true,
		"logger":  false,
	}

	if len(functions) != len(expected) {
		t.Fatalf("expected %d functions got %d", len(expected), len(functions))
	}

	for _, fn := range functions {
		if !fn.Exported && expected[fn.Name] {
			t.Fatalf("function %s should be exported", fn.Name)
		}
	}

	if functions[0].Name != "Execute" {
		t.Fatalf("expected Execute got %s", functions[0].Name)
	}

	if functions[1].Name != "Version" {
		t.Fatalf("expected Version got %s", functions[1].Name)
	}
}

func TestValidateFunctions(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Execute",
			Exported:    true,
			ReturnTypes: []string{"int"},
			Params: []compiler.Param{
				{Name: "a", Type: "int"},
			},
		},
	}
	if err := compiler.ValidateFunctions(functions); err != nil {
		t.Fatal(err)
	}
}

func TestSandboxValidation(t *testing.T) {
	source := `
package plugin

import "os/exec"

func Execute() {}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	pkg := &ast.Package{
		Name:  "plugin",
		Files: map[string]*ast.File{"plugin.go": file},
	}

	if err = compiler.ValidateSandbox(pkg); err == nil {
		t.Fatal("expected sandbox validation failure")
	}
}

func TestStructParameterAnalysis(t *testing.T) {
	source := `
package example

type User struct {
	Name string
}

func ProcessUser(user *User) {}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if len(functions) != 1 {
		t.Fatalf("expected 1 function got %d", len(functions))
	}

	if functions[0].Params[0].Type != "*User" {
		t.Fatalf("expected *User got %s", functions[0].Params[0].Type)
	}
}

func TestStructValidation(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:     "ProcessUser",
			Exported: true,
			Params: []compiler.Param{
				{Name: "user", Type: "*User"},
			},
		},
	}
	if err := compiler.ValidateFunctions(functions); err != nil {
		t.Fatal(err)
	}
}

func TestInterfaceParsing(t *testing.T) {
	source := `
package plugin

func Execute(v interface{}) {}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if functions[0].Params[0].Type != "interface{}" {
		t.Fatalf("expected interface{} got %s", functions[0].Params[0].Type)
	}
}

func TestInterfaceMarshaling(t *testing.T) {
	source := `
package plugin

func ProcessData(data interface{}) interface{} {
	return data
}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if len(functions) != 1 {
		t.Fatalf("expected 1 function got %d", len(functions))
	}

	fn := functions[0]
	if fn.Params[0].Type != "interface{}" {
		t.Fatalf("expected interface{} param got %s", fn.Params[0].Type)
	}
	if fn.ReturnTypes[0] != "interface{}" {
		t.Fatalf("expected interface{} return got %s", fn.ReturnTypes[0])
	}

	wrapper := compiler.GenerateWrapper("plugin", "plugin", "plugin/plugin", functions)

	checks := []string{
		"json.Unmarshal",
		"interface{}",
		"C.GoString",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing generated code for interface marshaling: %s\n\n%s", check, wrapper)
		}
	}
}

func TestInterfaceSliceParsing(t *testing.T) {
	source := `
package plugin

func ProcessList(items []interface{}) {}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if functions[0].Params[0].Type != "[]interface{}" {
		t.Fatalf("expected []interface{} got %s", functions[0].Params[0].Type)
	}
}

func TestMapParameterParsing(t *testing.T) {
	source := `
package plugin

func ProcessMap(data map[string]interface{}) {}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if functions[0].Params[0].Type != "map[string]interface{}" {
		t.Fatalf("expected map[string]interface{} got %s", functions[0].Params[0].Type)
	}
}

func TestMultipleReturnTypesParsing(t *testing.T) {
	source := `
package plugin

type User struct {
	Name string
	Age  int
}

func CreateUser(name string, age int) (*User, error) {
	return &User{Name: name, Age: age}, nil
}
`
	set := token.NewFileSet()
	file, err := parser.ParseFile(set, "plugin.go", source, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	functions := compiler.AnalyzeFunctions(file)
	if len(functions) != 1 {
		t.Fatalf("expected 1 function got %d", len(functions))
	}

	fn := functions[0]
	if len(fn.ReturnTypes) != 2 {
		t.Fatalf("expected 2 return types got %d", len(fn.ReturnTypes))
	}
	if fn.ReturnTypes[0] != "*User" {
		t.Fatalf("expected *User got %s", fn.ReturnTypes[0])
	}
	if fn.ReturnTypes[1] != "error" {
		t.Fatalf("expected error got %s", fn.ReturnTypes[1])
	}
}

func TestMultipleReturnTypesWrapperGeneration(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "CreateUser",
			Exported:    true,
			ReturnTypes: []string{"*User", "error"},
			Params: []compiler.Param{
				{Name: "name", Type: "string"},
				{Name: "age", Type: "int"},
			},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"result0, result1 :=",
		"errStr1",
		"json.Marshal",
		"\"error\": errStr1",
		"\"result0\": result0",
		"envelope",
		"\"type\": \"multi\"",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing multi-return generated code: %s\n\n%s", check, wrapper)
		}
	}
}

func TestMultipleReturnTypesOnlyError(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "DoSomething",
			Exported:    true,
			ReturnTypes: []string{"string", "error"},
			Params:      []compiler.Param{},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"result0, result1 :=",
		"errStr1",
		"\"result0\": result0",
		"\"error\": errStr1",
		"\"type\": \"multi\"",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing string+error return code: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingInt(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Execute",
			Exported:    true,
			ReturnTypes: []string{"int"},
			Params: []compiler.Param{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "int"},
			},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"scalar\"",
		"\"value\": result0",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for int return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingString(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Version",
			Exported:    true,
			ReturnTypes: []string{"string"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"string\"",
		"\"value\": result0",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for string return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingBool(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "IsReady",
			Exported:    true,
			ReturnTypes: []string{"bool"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"bool\"",
		"\"value\": result0",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for bool return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingError(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Validate",
			Exported:    true,
			ReturnTypes: []string{"error"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"error\"",
		"errVal",
		"result0.Error()",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for error return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingBytes(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "ReadData",
			Exported:    true,
			ReturnTypes: []string{"[]byte"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"bytes\"",
		"\"value\": result0",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for []byte return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestEnvelopeWrappingStruct(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "GetUser",
			Exported:    true,
			ReturnTypes: []string{"User"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example", "plugin/plugin", functions)

	checks := []string{
		"\"type\": \"json\"",
		"\"value\": result0",
		"C.CString(string(envelope))",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing envelope for struct return: %s\n\n%s", check, wrapper)
		}
	}
}

func TestGenerateWrapper_Int(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Execute",
			Exported:    true,
			ReturnTypes: []string{"int"},
			Params: []compiler.Param{
				{Name: "a", Type: "int"},
				{Name: "b", Type: "int"},
			},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example/plugin", "plugin", functions)

	checks := []string{
		"package main",
		"import \"C\"",
		"plugin \"plugin\"",
		"//export go_plugin_plugin_Execute",
		"func go_plugin_plugin_Execute",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing generated code: %s\n\n%s", check, wrapper)
		}
	}
}

func TestGenerateWrapper_String(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:        "Version",
			Exported:    true,
			ReturnTypes: []string{"string"},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example/plugin", "plugin", functions)

	checks := []string{
		"func go_plugin_plugin_Version",
		"plugin.Version()",
		"C.CString",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing generated code: %s", check)
		}
	}
}

func TestGenerateWrapper_Void(t *testing.T) {
	functions := []compiler.Function{
		{
			Name:     "Logger",
			Exported: true,
			Params: []compiler.Param{
				{Name: "message", Type: "string"},
			},
		},
	}

	wrapper := compiler.GenerateWrapper("plugin", "example/plugin", "plugin", functions)

	checks := []string{
		"func go_plugin_plugin_Logger",
		"C.GoString",
		"example/plugin.Logger",
	}
	for _, check := range checks {
		if !strings.Contains(wrapper, check) {
			t.Fatalf("missing generated code: %s", check)
		}
	}
}
