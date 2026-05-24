package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ParsePackage(
	dir string,
) (
	map[string]*ast.Package,
	error,
) {
	set := token.NewFileSet()

	return parser.ParseDir(
		set,
		dir,
		nil,
		parser.ParseComments,
	)
}
