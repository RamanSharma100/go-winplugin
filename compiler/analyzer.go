package compiler

import (
	"go/ast"
	"unicode"
)

type Function struct {
	Name string
}

func AnalyzeFunctions(
	file *ast.File,
) []Function {
	var functions []Function

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := fn.Name.Name

		if unicode.IsLower(
			rune(name[0]),
		) {
			continue
		}

		functions = append(
			functions,
			Function{
				Name: name,
			},
		)
	}

	return functions
}
