package main

import (
	"fmt"

	winplugin "github.com/RamanSharma100/go-winplugin"
)

func main() {
	pluginLoader, err := winplugin.NewLoader("./example")
	if err != nil {
		panic(err)
	}

	err = pluginLoader.Build("plugin.go")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("\n--- int return ---")
	output, err := pluginLoader.Call("example", "Execute", uintptr(10), uintptr(20))
	if err != nil {
		panic(err)
	}
	fmt.Println("Execute(10, 20):", output)

	output, err = pluginLoader.Call("example", "Execute", uintptr(0), uintptr(0))
	if err != nil {
		panic(err)
	}
	fmt.Println("Execute(0, 0):", output)

	output, err = pluginLoader.Call("example", "Execute", uintptr(1000000), uintptr(2000000))
	if err != nil {
		panic(err)
	}
	fmt.Println("Execute(1000000, 2000000):", output)

	fmt.Println("\n--- string return ---")
	output, err = pluginLoader.Call("example", "Version")
	if err != nil {
		panic(err)
	}
	fmt.Println("Version():", output)

	fmt.Println("\n--- void function ---")
	pluginLoader.Call("example", "Logger", "hello from plugin")
	pluginLoader.Call("example", "Logger", "")
	pluginLoader.Call("example", "Logger", "special chars: !@#$%^&*()")

	fmt.Println("\n--- bool return ---")
	output, err = pluginLoader.Call("example", "IsPositive", uintptr(5))
	if err != nil {
		panic(err)
	}
	fmt.Println("IsPositive(5):", output)

	output, err = pluginLoader.Call("example", "IsPositive", uintptr(0))
	if err != nil {
		panic(err)
	}
	fmt.Println("IsPositive(0):", output)

	fmt.Println("\n--- float64 return ---")
	output, err = pluginLoader.Call("example", "Divide", 7.0, 2.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Divide(7.0, 2.0):", output)

	output, err = pluginLoader.Call("example", "Divide", 10.0, 0.0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Divide(10.0, 0.0):", output)

	fmt.Println("\n--- []byte return ---")
	output, err = pluginLoader.Call("example", "ReadData")
	if err != nil {
		panic(err)
	}
	fmt.Println("ReadData():", output)

	fmt.Println("\n--- error-only return (nil error) ---")
	output, err = pluginLoader.Call("example", "Validate", "Jane")
	if err != nil {
		panic(err)
	}
	fmt.Println("Validate(\"Jane\"):", output)

	fmt.Println("\n--- error-only return (non-nil error) ---")
	output, err = pluginLoader.Call("example", "Validate", "")
	if err != nil {
		fmt.Println("Validate(\"\") error (expected):", err)
	} else {
		fmt.Println("Validate(\"\"):", output)
	}

	fmt.Println("\n--- interface{} param: map ---")
	output, err = pluginLoader.Call("example", "ProcessData", map[string]any{
		"message": "Hello, World!",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(map):", output)

	fmt.Println("\n--- interface{} param: string ---")
	output, err = pluginLoader.Call("example", "ProcessData", "plain string value")
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(string):", output)

	fmt.Println("\n--- interface{} param: number ---")
	output, err = pluginLoader.Call("example", "ProcessData", 42)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(int):", output)

	fmt.Println("\n--- interface{} param: bool ---")
	output, err = pluginLoader.Call("example", "ProcessData", true)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(bool):", output)

	fmt.Println("\n--- interface{} param: slice ---")
	output, err = pluginLoader.Call("example", "ProcessData", []any{1, "two", true})
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(slice):", output)

	fmt.Println("\n--- interface{} param: nil ---")
	output, err = pluginLoader.Call("example", "ProcessData", nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessData(nil):", output)

	fmt.Println("\n--- struct pointer param ---")
	output, err = pluginLoader.Call("example", "ProcessUser", map[string]any{
		"Name": "John",
		"Age":  30,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("ProcessUser({Name:John Age:30}):", output)

	fmt.Println("\n--- multi-return (*User, nil error) ---")
	output, err = pluginLoader.Call("example", "CreateUser", "Jane Doe", 25)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateUser(\"Jane Doe\", 25):", output)
	if user, ok := output.(map[string]any); ok {
		fmt.Println("  Name:", user["Name"])
		fmt.Println("  Age:", user["Age"])
	}

	fmt.Println("\n--- multi-return (*User, nil error) empty name ---")
	output, err = pluginLoader.Call("example", "CreateUser", "", 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateUser(\"\", 0):", output)

	fmt.Println("\n--- multi-return (*User, non-nil error) ---")
	output, err = pluginLoader.Call("example", "CreateUserWithError", "Bad", -1)
	if err != nil {
		fmt.Println("CreateUserWithError(\"Bad\", -1) error (expected):", err)
	} else {
		fmt.Println("CreateUserWithError(\"Bad\", -1):", output)
	}

	fmt.Println("\n--- multi-return (*User, nil error) valid ---")
	output, err = pluginLoader.Call("example", "CreateUserWithError", "Alice", 30)
	if err != nil {
		panic(err)
	}
	fmt.Println("CreateUserWithError(\"Alice\", 30):", output)

	fmt.Println("\nAll scenarios completed.")
}
