package compiler

import (
	"go/ast"
	"unicode"
)

type Param struct {
	Name        string
	Type        string
	isInterface bool
}

type Function struct {
	Name        string
	Exported    bool
	Params      []Param
	ReturnTypes []string
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
						Name:        nameIdent.Name,
						Type:        paramType,
						isInterface: paramType == "interface{}",
					})
				}

				if len(p.Names) == 0 {
					params = append(params, Param{
						Name:        "p" + string(rune('0'+i)),
						Type:        paramType,
						isInterface: paramType == "interface{}",
					})
				}
			}
		}

		returnTypes := []string{}

		if fn.Type.Results != nil {
			for _, r := range fn.Type.Results.List {
				returnTypes = append(
					returnTypes,
					exprToString(r.Type),
				)
			}
		}

		functions = append(functions, Function{
			Name:        name,
			Exported:    exported,
			Params:      params,
			ReturnTypes: returnTypes,
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

func exprToString(
	e ast.Expr,
) string {
	switch v := e.(type) {
	case *ast.ArrayType:
		return "[]" + exprToString(v.Elt)
	case *ast.MapType:
		return "map[" + exprToString(v.Key) + "]" + exprToString(v.Value)
	case *ast.Ellipsis:
		return "..." + exprToString(v.Elt)
	case *ast.Ident:
		return v.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.StarExpr:
		return "*" + exprToString(
			v.X,
		)
	case *ast.StructType:
		return "struct"
	case *ast.SelectorExpr:
		return exprToString(
			v.X,
		) + "." + v.Sel.Name
	default:
		return "unknown"
	}
}
