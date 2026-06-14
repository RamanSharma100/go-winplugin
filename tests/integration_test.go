package tests

import (
	"testing"

	winplugin "github.com/RamanSharma100/go-winplugin"
)

func TestPluginExecution(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}

	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Execute", uintptr(10), uintptr(20))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64 got %T: %v", result, result)
	}
	if val != 30 {
		t.Fatalf("expected 30 got %d", val)
	}

	_, err = loader.Call("sample_plugin", "Version")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiFilePluginBuild(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}

	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}
}

func TestIntReturnValue(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Execute", uintptr(5), uintptr(7))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64 got %T: %v", result, result)
	}
	if val != 12 {
		t.Fatalf("expected 12 got %d", val)
	}
}

func TestStringReturnValue(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Version")
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(string)
	if !ok {
		t.Fatalf("expected string got %T: %v", result, result)
	}
	if val == "" {
		t.Fatal("expected non-empty version string")
	}
}

func TestVoidFunction(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Logger", "hello test")
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Fatalf("expected nil result for void function got %v", result)
	}
}

func TestInterfaceParameter(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "ProcessData", map[string]any{
		"message": "Hello, World!",
	})
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

func TestInterfaceParameterString(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "ProcessData", "plain string value")
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

func TestInterfaceParameterNumber(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "ProcessData", 42)
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

func TestInterfaceParameterSlice(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "ProcessData", []any{1, "two", true})
	if err != nil {
		t.Fatal(err)
	}
	_ = result
}

func TestMultipleReturnValues(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "CreateUser", "Jane Doe", 25)
	if err != nil {
		t.Fatal(err)
	}

	user, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any got %T: %v", result, result)
	}
	if user["Name"] != "Jane Doe" {
		t.Fatalf("expected Name=Jane Doe got %v", user["Name"])
	}
	age, ok := user["Age"].(float64)
	if !ok {
		t.Fatalf("expected float64 Age got %T", user["Age"])
	}
	if int(age) != 25 {
		t.Fatalf("expected Age=25 got %v", age)
	}
}

func TestMultipleReturnValuesWithError(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "CreateUserWithError", "Bad", -1)
	if err == nil {
		t.Fatalf("expected error got result: %v", result)
	}
}

func TestStructParameter(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "ProcessUser", map[string]any{
		"Name": "John",
		"Age":  30,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result != nil {
		t.Fatalf("expected nil result for void function got %v", result)
	}
}

func TestZeroIntReturn(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Execute", uintptr(0), uintptr(0))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64 got %T: %v", result, result)
	}
	if val != 0 {
		t.Fatalf("expected 0 got %d", val)
	}
}

func TestNegativeIntReturn(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Execute", uintptr(5), uintptr(3))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64 got %T: %v", result, result)
	}
	if val != 2 {
		t.Fatalf("expected 2 got %d", val)
	}
}

func TestBoolReturnTrue(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "IsPositive", uintptr(5))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(bool)
	if !ok {
		t.Fatalf("expected bool got %T: %v", result, result)
	}
	if !val {
		t.Fatal("expected true got false")
	}
}

func TestBoolReturnFalse(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "IsPositive", uintptr(0))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(bool)
	if !ok {
		t.Fatalf("expected bool got %T: %v", result, result)
	}
	if val {
		t.Fatal("expected false got true")
	}
}

func TestFloat64Return(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Divide", 7.0, 2.0)
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(float64)
	if !ok {
		t.Fatalf("expected float64 got %T: %v", result, result)
	}
	if val != 3.5 {
		t.Fatalf("expected 3.5 got %f", val)
	}
}

func TestEmptyStringParameter(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "CreateUser", "", 0)
	if err != nil {
		t.Fatal(err)
	}

	user, ok := result.(map[string]any)
	if !ok {
		t.Fatalf("expected map got %T: %v", result, result)
	}
	if user["Name"] != "" {
		t.Fatalf("expected empty Name got %v", user["Name"])
	}
}

func TestLargeIntValue(t *testing.T) {
	loader, err := winplugin.NewLoader("./fixtures/sample_plugin")
	if err != nil {
		t.Fatal(err)
	}
	if err = loader.Build("plugin.go"); err != nil {
		t.Fatal(err)
	}

	result, err := loader.Call("sample_plugin", "Execute", uintptr(1000000), uintptr(2000000))
	if err != nil {
		t.Fatal(err)
	}

	val, ok := result.(int64)
	if !ok {
		t.Fatalf("expected int64 got %T: %v", result, result)
	}
	if val != 3000000 {
		t.Fatalf("expected 3000000 got %d", val)
	}
}
