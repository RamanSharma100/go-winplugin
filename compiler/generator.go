package compiler

import (
	"strings"
)

func GenerateWrapper(packageName string, functions []Function) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import \"C\"\n\n")

	sb.WriteString("import (\n")
	sb.WriteString("\t\"" + packageName + "\"\n")
	sb.WriteString(")\n\n")

	for _, fn := range functions {
		if !fn.Exported {
			continue
		}

		exportName := MangleSymbol(packageName, fn.Name)

		cParams, goArgs, pre, post, retType := buildCBridge(fn)

		sb.WriteString("//export " + exportName + "\n")
		sb.WriteString("func " + exportName + "(" + cParams + ") " + retType + " {\n")

		if pre != "" {
			sb.WriteString(pre)
		}

		if len(fn.Params) > 0 {
			sb.WriteString("\tresult := " + packageName + "." + fn.Name + "(" + goArgs + ")\n")
		} else {
			sb.WriteString("\tresult := " + packageName + "." + fn.Name + "()\n")
		}

		if post != "" {
			sb.WriteString(post)
		}

		if retType == "" {
			sb.WriteString("\t_ = result\n")
		} else {
			sb.WriteString("\treturn result\n")
		}

		sb.WriteString("}\n\n")
	}

	return sb.String()
}

func buildCBridge(fn Function) (string, string, string, string, string) {
	var cParams []string
	var goArgs []string
	var pre strings.Builder
	var post strings.Builder
	retType := ""

	for i, p := range fn.Params {
		cName := "p" + string(rune('0'+i))
		goName := p.Name

		switch p.Type {
		case "int":
			cParams = append(cParams, cName+" C.int")
			goArgs = append(goArgs, "int("+cName+")")

		case "string":
			cParams = append(cParams, cName+" *C.char")
			pre.WriteString("\t" + goName + " := C.GoString(" + cName + ")\n")
			goArgs = append(goArgs, goName)

		case "float":
			cParams = append(cParams, cName+" C.double")
			goArgs = append(goArgs, "float64("+cName+")")

		case "bool":
			cParams = append(cParams, cName+" C.int")
			goArgs = append(goArgs, cName+" != 0")

		default:
			cParams = append(cParams, cName+" C.int")
			goArgs = append(goArgs, "int("+cName+")")
		}
	}

	switch fn.ReturnType {
	case "int":
		retType = "C.int"
		post.WriteString("\treturn C.int(result)\n")

	case "string":
		retType = "*C.char"
		post.WriteString("\treturn C.CString(result)\n")

	case "float":
		retType = "C.double"
		post.WriteString("\treturn C.double(result)\n")

	case "bool":
		retType = "C.int"
		post.WriteString("\tif result { return 1 } else { return 0 }\n")

	case "void", "":
		retType = ""
	}

	return strings.Join(cParams, ", "), strings.Join(goArgs, ", "), pre.String(), post.String(), retType
}
