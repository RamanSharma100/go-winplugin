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
- Struct type detection
- Dynamic function execution

## Supported Types

| Type     | Support |
| -------- | ------- |
| int      | Yes     |
| string   | Yes     |
| bool     | Yes     |
| float    | Yes     |
| void     | Yes     |
| struct   | Partial |
| \*struct | Partial |

### Struct Support

Version 0.2.0 introduces struct parsing and validation support.

Supported:

- Struct type discovery
- Struct pointer discovery
- Wrapper generation compatibility
- Validation support

Not yet supported:

- Struct marshalling across DLL boundaries
- Automatic struct serialization
- Interface marshalling

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

	result, err := loader.Call(
		"example",
		"Execute",
		uintptr(10),
		uintptr(20),
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(result)
}
```

Output:

```text
30
```

## Multi-File Plugin Example

plugin.go

```go
package example

func Execute(a, b int) int {
	return AddInternal(a, b)
}
```

testing.go

```go
package example

func AddInternal(a, b int) int {
	return a + b
}
```

The loader automatically analyzes all Go files within the package.

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
9. Runtime function execution

## Environment Setup

The library automatically:

- Detects Go installation
- Detects GCC installation
- Detects MSYS2 installation
- Validates PATH configuration
- Installs missing dependencies
- Prepares the CGO environment

## Security

Version 0.2.0 introduces basic sandbox validation.

Blocked imports:

- os
- os/exec
- syscall
- unsafe

Plugins using blocked imports will fail validation before compilation.

## Known Limitations

- Interface support not implemented
- String memory cleanup pending
- Windows-first implementation
- Requires CGO enabled environment
- No hot reload support

## Roadmap

### v0.3.0

- Interface support
- Runtime type validation
- Plugin sandboxing improvements

### v0.4.0

- Hot reload support
- Symbol cache enhancements
- Cross-platform runtime loaders

## Testing

```bash
go test ./...
```

Current test coverage includes:

- Symbol mangling
- Function analysis
- Package analysis
- Validation
- Wrapper generation
- Sandbox validation
- Struct parsing
- Multi-file plugins
- DLL compilation
- DLL execution

## Documentation

- Contributing Guide
- Changelog

## License

MIT

## Author

Raman Sharma

GitHub:
https://github.com/RamanSharma100
