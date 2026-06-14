package compiler

import (
	"fmt"
	"go/ast"
)

var blockedImports = map[string]bool{
	"os/exec": true,
	"syscall": true,
	"unsafe":  true,
}

func ValidateFunctions(functions []Function) error {
	visited := map[string]bool{}

	for _, fn := range functions {
		if !fn.Exported {
			continue
		}

		if visited[fn.Name] {
			return fmt.Errorf("duplicate export: %s", fn.Name)
		}

		visited[fn.Name] = true

		if fn.ReturnType != "" && fn.ReturnType != "void" &&
			fn.ReturnType != "int" &&
			fn.ReturnType != "string" &&
			fn.ReturnType != "float" &&
			fn.ReturnType != "bool" &&
			fn.ReturnType != "interface{}" {
			return fmt.Errorf("unsupported return type: %s in %s", fn.ReturnType, fn.Name)
		}

		for _, p := range fn.Params {
			switch p.Type {
			case "int",
				"string",
				"float",
				"bool",
				"struct",
				"interface{}":
				continue
			default:
				if len(p.Type) > 0 && p.Type[0] == '*' {
					continue
				}

				return fmt.Errorf(
					"unsupported parameter type: %s in %s",
					p.Type,
					fn.Name,
				)
			}
		}
	}

	return nil
}

func ValidateCall(fn Function, args []any) error {
	if len(args) != len(fn.Params) {
		return fmt.Errorf(
			"expected %d arguments got %d",
			len(fn.Params),
			len(args),
		)
	}

	for i, p := range fn.Params {
		switch p.Type {
		case "int":
			if _, ok := args[i].(int); !ok {
				return fmt.Errorf(
					"parameter %s expects int",
					p.Name,
				)
			}
		case "string":
			if _, ok := args[i].(string); !ok {
				return fmt.Errorf(
					"parameter %s expects string",
					p.Name,
				)
			}
		case "float":
			if _, ok := args[i].(float64); !ok {
				return fmt.Errorf(
					"parameter %s expects float64",
					p.Name,
				)
			}
		case "bool":
			if _, ok := args[i].(bool); !ok {
				return fmt.Errorf(
					"parameter %s expects bool",
					p.Name,
				)
			}
		}
	}

	return nil
}

func ValidateSandbox(
	pkg *ast.Package,
) error {
	for _, file := range pkg.Files {
		for _, imp := range file.Imports {
			path := imp.Path.Value
			if len(path) >= 2 {
				path = path[1 : len(path)-1]
			}

			if blockedImports[path] {
				return fmt.Errorf(
					"sandbox violation: %s import is not allowed",
					path,
				)
			}
		}
	}

	return nil
}
