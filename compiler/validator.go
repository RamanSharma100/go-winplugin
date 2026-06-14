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

		for _, returnType := range fn.ReturnTypes {

			switch returnType {

			case "",
				"void",

				"int",
				"int8",
				"int16",
				"int32",
				"int64",

				"uint",
				"uint8",
				"uint16",
				"uint32",
				"uint64",

				"float32",
				"float64",

				"string",
				"bool",

				"error",

				"interface{}",
				"struct":

				continue

			default:

				if len(returnType) > 0 && returnType[0] == '*' {
					continue
				}

				if len(returnType) > 2 && returnType[:2] == "[]" {
					continue
				}

				if len(returnType) > 4 && returnType[:4] == "map[" {
					continue
				}

				return fmt.Errorf(
					"unsupported return type: %s in %s",
					returnType,
					fn.Name,
				)
			}
		}

		for _, p := range fn.Params {
			switch p.Type {
			case "int",
				"int8",
				"int16",
				"int32",
				"int64",

				"uint",
				"uint8",
				"uint16",
				"uint32",
				"uint64",

				"float32",
				"float64",

				"string",
				"bool",

				"error",

				"struct",
				"interface{}":

				continue
			default:
				if len(p.Type) > 0 && p.Type[0] == '*' {
					continue
				}

				if len(p.Type) > 2 && p.Type[:2] == "[]" {
					continue
				}

				if len(p.Type) > 4 && p.Type[:4] == "map[" {
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

		case "int8":
			if _, ok := args[i].(int8); !ok {
				return fmt.Errorf(
					"parameter %s expects int8",
					p.Name,
				)
			}

		case "int16":
			if _, ok := args[i].(int16); !ok {
				return fmt.Errorf(
					"parameter %s expects int16",
					p.Name,
				)
			}

		case "int32":
			if _, ok := args[i].(int32); !ok {
				return fmt.Errorf(
					"parameter %s expects int32",
					p.Name,
				)
			}

		case "int64":
			if _, ok := args[i].(int64); !ok {
				return fmt.Errorf(
					"parameter %s expects int64",
					p.Name,
				)
			}

		case "uint":
			if _, ok := args[i].(uint); !ok {
				return fmt.Errorf(
					"parameter %s expects uint",
					p.Name,
				)
			}

		case "uint8":
			if _, ok := args[i].(uint8); !ok {
				return fmt.Errorf(
					"parameter %s expects uint8",
					p.Name,
				)
			}

		case "uint16":
			if _, ok := args[i].(uint16); !ok {
				return fmt.Errorf(
					"parameter %s expects uint16",
					p.Name,
				)
			}

		case "uint32":
			if _, ok := args[i].(uint32); !ok {
				return fmt.Errorf(
					"parameter %s expects uint32",
					p.Name,
				)
			}

		case "uint64":
			if _, ok := args[i].(uint64); !ok {
				return fmt.Errorf(
					"parameter %s expects uint64",
					p.Name,
				)
			}

		case "float32":
			if _, ok := args[i].(float32); !ok {
				return fmt.Errorf(
					"parameter %s expects float32",
					p.Name,
				)
			}

		case "float64":
			if _, ok := args[i].(float64); !ok {
				return fmt.Errorf(
					"parameter %s expects float64",
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

		case "bool":
			if _, ok := args[i].(bool); !ok {
				return fmt.Errorf(
					"parameter %s expects bool",
					p.Name,
				)
			}

		case "interface{}":
			continue

		default:

			if IsPointerType(p.Type) {
				continue
			}

			if IsSliceType(p.Type) {
				continue
			}

			if IsMapType(p.Type) {
				continue
			}

			if p.Type == "struct" {
				continue
			}

			if p.Type == "interface{}" {
				continue
			}

			return fmt.Errorf(
				"unsupported parameter type %s for %s",
				p.Type,
				p.Name,
			)
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

func IsPointerType(t string) bool {
	return len(t) > 0 && t[0] == '*'
}

func IsSliceType(t string) bool {
	return len(t) > 2 && t[:2] == "[]"
}

func IsMapType(t string) bool {
	return len(t) > 4 && t[:4] == "map["
}

func IsSupportedPrimitive(t string) bool {
	switch t {

	case "int",
		"int8",
		"int16",
		"int32",
		"int64",

		"uint",
		"uint8",
		"uint16",
		"uint32",
		"uint64",

		"float32",
		"float64",

		"string",
		"bool",

		"error",

		"interface{}":

		return true
	}

	return false
}
