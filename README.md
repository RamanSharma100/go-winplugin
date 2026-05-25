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
- Automatic wrapper generation
- Automatic GCC/MSYS2 setup
- Temporary isolated workspaces
- Primitive type support
- Dynamic function execution

## Supported Types

| Type   | Support |
| ------ | ------- |
| int    | Yes     |
| string | Yes     |
| bool   | Yes     |
| float  | Yes     |
| void   | Yes     |

## How It Works

```
User Plugin → AST Parser → Wrapper Generator → CGO Bridge → go build -buildmode=c-shared → DLL Generation → Windows DLL Loader → Runtime Symbol Execution
```

## Installation

```bash
go get github.com/RamanSharma100/go-winplugin
```

## Quick Start

### Create a plugin

```go
package example

func Execute(a int, b int) int {
	return a + b
}
```

### Use the plugin

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

	fmt.Println(result) // Output: 30
}
```

## Build Process

The library automatically handles:

1. Go source parsing
2. Exported function detection
3. CGO wrapper generation
4. Isolated workspace creation
5. DLL compilation
6. Dynamic DLL loading
7. Symbol resolution
8. Runtime function execution

## Environment Setup

The library automatically:

- Detects GCC installation
- Detects MSYS2 installation
- Validates PATH configuration
- Installs missing dependencies
- Prepares the CGO environment

## Known Limitations

- Struct support not implemented
- Interface support not implemented
- String memory cleanup pending
- Windows-first implementation
- Requires CGO enabled environment
- No hot reload support

## Roadmap

- Struct support
- Reflection support
- Plugin sandboxing
- Hot reload support
- Symbol cache
- Linux support
- macOS support
- Runtime type validation
- Automatic memory cleanup

## Testing

```bash
go test ./...
```

## Documentation

- [Contributing](CONTRIBUTING.md)
- [Changelog](CHANGELOG.md)

## License

MIT

## Author

[Raman Sharma](https://github.com/RamanSharma100)
