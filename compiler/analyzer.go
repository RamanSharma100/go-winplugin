package compiler

import (
	"go/ast"
	"unicode"
)

type Param struct {
	Name string
	Type string
}

type Function struct {
	Name       string
	Exported   bool
	Params     []Param
	ReturnType string
}

func AnalyzeFunctions(file *ast.File) []Function {
	var functions []Function

	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		name := fn.Name.Name

		exported := unicode.IsUpper(rune(name[0]))

		var params []Param
		if fn.Type.Params != nil {
			for i, p := range fn.Type.Params.List {
				paramType := exprToString(p.Type)

				for _, nameIdent := range p.Names {
					params = append(params, Param{
						Name: nameIdent.Name,
						Type: paramType,
					})
				}

				if len(p.Names) == 0 {
					params = append(params, Param{
						Name: "p" + string(rune('0'+i)),
						Type: paramType,
					})
				}
			}
		}

		returnType := "void"
		if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 {
			returnType = exprToString(fn.Type.Results.List[0].Type)
		}

		functions = append(functions, Function{
			Name:       name,
			Exported:   exported,
			Params:     params,
			ReturnType: returnType,
		})
	}

	return functions
}

func AnalyzePackageFunctions(
	pkg *ast.Package,
) []Function {
	functions := []Function{}

	for _, file := range pkg.Files {
		fileFunctions := AnalyzeFunctions(file)

		functions = append(
			functions,
			fileFunctions...,
		)
	}

	return functions
}

func exprToString(e ast.Expr) string {
	switch v := e.(type) {
	case *ast.Ident:
		return v.Name
	default:
		return "unknown"
	}
}
