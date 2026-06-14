# go-winplugin

Windows-compatible runtime plugin system for Go using CGO-generated shared libraries and dynamic DLL loading.

## Why This Exists

Go's official `plugin` package does not support Windows, creating major limitations for runtime extensibility, dynamic module systems, AI agents, automation frameworks, game scripting, and runtime service injection.

`go-winplugin` solves this by generating CGO bridges and loading Windows DLLs dynamically at runtime.

## Features

- Windows runtime plugin support
- Dynamic DLL generation
- CGO bridge generation
- Runtime symbol execution
- AST-based exported function parsing
- Package-wide function discovery
- Multi-file plugin support
- Automatic wrapper generation
- Automatic GCC/MSYS2 setup
- Temporary isolated workspaces
- Plugin sandbox validation
- Symbol caching
- Primitive type support
- Struct type support
- Interface type support
- Multi-return value support
- Typed envelope return system
- Dynamic function execution

## Supported Types

| Type            | Support |
| --------------- | ------- |
| int             | Yes     |
| int32           | Yes     |
| int64           | Yes     |
| uint64          | Yes     |
| float32         | Yes     |
| float64         | Yes     |
| string          | Yes     |
| bool            | Yes     |
| void            | Yes     |
| error           | Yes     |
| []byte          | Yes     |
| struct          | Yes     |
| \*struct        | Yes     |
| interface{}     | Yes     |
| []interface{}   | Yes     |
| map[string]any  | Yes     |
| multiple return | Yes     |

## How It Works

```text
User Plugin
    ↓
AST Parser
    ↓
Package Analyzer
    ↓
Sandbox Validator
    ↓
Wrapper Generator
    ↓
CGO Bridge
    ↓
go build -buildmode=c-shared
    ↓
DLL Generation
    ↓
Windows DLL Loader
    ↓
Symbol Cache
    ↓
Runtime Symbol Execution
    ↓
Typed Envelope Decoder
```

## Installation

```bash
go get github.com/RamanSharma100/go-winplugin
```

## Quick Start

### Plugin

```go
package example

func Execute(a int, b int) int {
    return a + b
}
```

### Host

```go
package main

import (
    "fmt"

    winplugin "github.com/RamanSharma100/go-winplugin"
)

func main() {
    loader, err := winplugin.NewLoader("./example")
    if err != nil {
        panic(err)
    }

    err = loader.Build("plugin.go")
    if err != nil {
        panic(err)
    }

    result, err := loader.Call("example", "Execute", uintptr(10), uintptr(20))
    if err != nil {
        panic(err)
    }

    fmt.Println(result) // 30
}
```

## Supported Scenarios

### Primitive Types

```go
// int
func Execute(a, b int) int { return a + b }

// string
func Version() string { return "1.0.0" }

// bool
func IsPositive(n int) bool { return n > 0 }

// float64
func Divide(a, b float64) float64 { return a / b }

// void
func Logger(message string) { fmt.Println(message) }

// []byte
func ReadData() []byte { return []byte("hello") }

// error only
func Validate(name string) error {
    if name == "" {
        return errors.New("name cannot be empty")
    }
    return nil
}
```

### Interface Parameters

```go
// interface{} accepts map, string, number, bool, slice, nil
func ProcessData(data interface{}) {
    fmt.Println(data)
}
```

```go
// map
loader.Call("example", "ProcessData", map[string]any{"key": "value"})

// string
loader.Call("example", "ProcessData", "hello")

// number
loader.Call("example", "ProcessData", 42)

// bool
loader.Call("example", "ProcessData", true)

// slice
loader.Call("example", "ProcessData", []any{1, "two", true})

// nil
loader.Call("example", "ProcessData", nil)
```

### Struct Parameters

```go
type User struct {
    Name string
    Age  int
}

func ProcessUser(user *User) {
    fmt.Println(user.Name, user.Age)
}
```

```go
loader.Call("example", "ProcessUser", map[string]any{
    "Name": "John",
    "Age":  30,
})
```

### Multiple Return Values

```go
func CreateUser(name string, age int) (*User, error) {
    return &User{Name: name, Age: age}, nil
}
```

```go
result, err := loader.Call("example", "CreateUser", "Jane", 25)
if err != nil {
    panic(err)
}

user := result.(map[string]any)
fmt.Println(user["Name"]) // Jane
fmt.Println(user["Age"])  // 25
```

### Error Propagation

```go
func CreateUserWithError(name string, age int) (*User, error) {
    if age < 0 {
        return nil, errors.New("age cannot be negative")
    }
    return &User{Name: name, Age: age}, nil
}
```

```go
result, err := loader.Call("example", "CreateUserWithError", "Bad", -1)
if err != nil {
    fmt.Println(err) // age cannot be negative
}
```

## Multi-File Plugin Example

plugin.go

```go
package example

func Execute(a, b int) int {
    return AddInternal(a, b)
}
```

internal.go

```go
package example

func AddInternal(a, b int) int {
    return a + b
}
```

The loader automatically analyzes all Go files within the package.

## Return Value System

v0.3.0 introduces a typed envelope system for all return values. Every DLL return is wrapped as:

```json
{ "type": "scalar|string|bool|error|bytes|json|multi", "value": <data> }
```

This eliminates all pointer/scalar ambiguity and makes return value decoding fully deterministic regardless of value size.

| Envelope Type | Go Types                                    |
| ------------- | ------------------------------------------- |
| scalar        | int, int32, int64, uint64, float32, float64 |
| string        | string                                      |
| bool          | bool                                        |
| error         | error                                       |
| bytes         | []byte                                      |
| json          | struct, \*struct, map, interface            |
| multi         | any function with 2+ returns                |

## Build Process

The library automatically handles:

1. Package parsing
2. Exported function detection
3. Sandbox validation
4. CGO wrapper generation
5. Workspace creation
6. DLL compilation
7. Dynamic DLL loading
8. Symbol caching
9. Typed envelope decoding
10. Runtime function execution

## Environment Setup

The library automatically:

- Detects Go installation
- Detects GCC installation
- Detects MSYS2 installation
- Validates PATH configuration
- Installs missing dependencies
- Prepares the CGO environment

## Security

Sandbox validation blocks plugins from importing:

- os
- os/exec
- syscall
- unsafe

Plugins using blocked imports will fail validation before compilation.

## Known Limitations

- Windows-only (Linux/macOS not supported)
- String memory cleanup pending
- No hot reload support
- Requires CGO-enabled environment (GCC via MSYS2)

## Roadmap

### v0.4.0

- Hot reload support
- Symbol cache enhancements
- Cross-platform runtime loaders
- Plugin versioning

### v0.5.0

- Async plugin execution
- Plugin isolation improvements
- Performance benchmarks

## Testing

```bash
go test ./...
```

Test coverage includes:

- Symbol mangling
- Function analysis
- Package analysis
- Validation
- Wrapper generation
- Sandbox validation
- Struct parsing
- Interface marshaling
- Multi-return value parsing
- Envelope type verification
- DLL compilation
- DLL execution
- Edge cases: zero values, negative numbers, large integers, empty strings, nil interfaces

## Changelog

### v0.3.0

- Interface parameter support (`interface{}`, `[]interface{}`, `map[string]any`)
- Interface marshaling via JSON across DLL boundary
- Multi-return value support with full error propagation
- Typed envelope return system — eliminates pointer/scalar ambiguity
- `[]byte` return type support
- `error`-only return type support
- Struct and `*struct` parameter marshaling via JSON
- Fixed CGO preprocessor crash on multi-return functions (inline IIFE in generated code)
- `Call` return type changed from `uintptr` to `any`
- Expanded compiler test suite with envelope and interface coverage
- Expanded integration test suite with all type edge cases

### v0.2.0

- Struct type discovery
- Struct pointer discovery
- Wrapper generation compatibility
- Sandbox validation

### v0.1.0

- Initial release
- Windows runtime plugin support
- Dynamic DLL generation
- CGO bridge generation
- Primitive type support

## Documentation

- [Contributing Guide](CONTRIBUTING.md)
- [Changelog](CHANGELOG.md)

## License

MIT

## Author

Raman Sharma

GitHub: https://github.com/RamanSharma100
