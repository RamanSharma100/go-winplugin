package compiler

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func ParseFile(path string) (*ast.File, error) {
	fs := token.NewFileSet()

	return parser.ParseFile(
		fs,
		path,
		nil,
		parser.ParseComments,
	)
}
