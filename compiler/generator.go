package compiler

import "strings"

func GenerateWrapper(
	moduleName string,
	importAlias string,
	importPath string,
	functions []Function,
) string {
	var sb strings.Builder

	sb.WriteString("package main\n\n")
	sb.WriteString("import \"C\"\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"encoding/json\"\n")
	sb.WriteString("\t")
	sb.WriteString(importAlias)
	sb.WriteString(" \"")
	sb.WriteString(importPath)
	sb.WriteString("\"\n")
	sb.WriteString(")\n\n")

	for _, fn := range functions {
		if !fn.Exported {
			continue
		}

		exportName := MangleSymbol(
			moduleName,
			fn.Name,
		)

		cParams,
			goArgs,
			pre,
			post,
			retType := buildCBridge(
			fn,
			importAlias,
		)

		sb.WriteString("//export ")
		sb.WriteString(exportName)
		sb.WriteString("\n")

		sb.WriteString("func ")
		sb.WriteString(exportName)
		sb.WriteString("(")
		sb.WriteString(cParams)
		sb.WriteString(") ")

		if retType != "" {
			sb.WriteString(retType)
		}

		sb.WriteString(" {\n")

		if pre != "" {
			sb.WriteString(pre)
		}

		call :=
			importAlias +
				"." +
				fn.Name +
				"(" +
				goArgs +
				")"

		if len(fn.Params) == 0 {
			call =
				importAlias +
					"." +
					fn.Name +
					"()"
		}

		if retType == "" {
			sb.WriteString("\t")
			sb.WriteString(call)
			sb.WriteString("\n")
		} else {
			sb.WriteString("\tresult := ")
			sb.WriteString(call)
			sb.WriteString("\n")
		}

		if post != "" {
			sb.WriteString(post)
		}

		sb.WriteString("}\n\n")
	}

	sb.WriteString("func main() {}\n")

	return sb.String()
}

func buildCBridge(
	fn Function,
	importAlias string,
) (
	string,
	string,
	string,
	string,
	string,
) {
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
			cParams = append(
				cParams,
				cName+" C.int",
			)

			goArgs = append(
				goArgs,
				"int("+cName+")",
			)

		case "string":
			cParams = append(
				cParams,
				cName+" *C.char",
			)

			pre.WriteString(
				"\t" +
					goName +
					" := C.GoString(" +
					cName +
					")\n",
			)

			goArgs = append(
				goArgs,
				goName,
			)

		case "float64":
			cParams = append(
				cParams,
				cName+" C.double",
			)

			goArgs = append(
				goArgs,
				"float64("+cName+")",
			)

		case "bool":
			cParams = append(
				cParams,
				cName+" C.int",
			)

			goArgs = append(
				goArgs,
				cName+" != 0",
			)

		default:
			cParams = append(
				cParams,
				cName+" *C.char",
			)

			if strings.HasPrefix(
				p.Type,
				"*",
			) {

				targetType :=
					strings.TrimPrefix(
						p.Type,
						"*",
					)

				if !strings.Contains(
					targetType,
					".",
				) {
					targetType =
						importAlias +
							"." +
							targetType
				}

				pre.WriteString(
					"\t" +
						goName +
						" := &" +
						targetType +
						"{}\n",
				)

				pre.WriteString(
					"\tjson.Unmarshal([]byte(C.GoString(" +
						cName +
						")), " +
						goName +
						")\n",
				)

				goArgs = append(
					goArgs,
					goName,
				)

			} else {

				targetType := p.Type

				if !strings.Contains(
					targetType,
					".",
				) &&
					targetType != "interface{}" &&
					targetType != "map" &&
					targetType != "[]byte" {

					targetType =
						importAlias +
							"." +
							targetType
				}

				pre.WriteString(
					"\tvar " +
						goName +
						" " +
						targetType +
						"\n",
				)

				pre.WriteString(
					"\tjson.Unmarshal([]byte(C.GoString(" +
						cName +
						")), &" +
						goName +
						")\n",
				)

				goArgs = append(
					goArgs,
					goName,
				)
			}
		}
	}

	switch fn.ReturnType {

	case "int":
		retType = "C.int"
		post.WriteString(
			"\treturn C.int(result)\n",
		)

	case "string":
		retType = "*C.char"
		post.WriteString(
			"\treturn C.CString(result)\n",
		)

	case "float64":
		retType = "C.double"
		post.WriteString(
			"\treturn C.double(result)\n",
		)

	case "bool":
		retType = "C.int"
		post.WriteString(
			"\tif result { return 1 }\n\treturn 0\n",
		)

	case "void", "":
		retType = ""

	default:
		retType = "*C.char"

		post.WriteString(
			"\tpayload, _ := json.Marshal(result)\n",
		)

		post.WriteString(
			"\treturn C.CString(string(payload))\n",
		)
	}

	return strings.Join(cParams, ", "),
		strings.Join(goArgs, ", "),
		pre.String(),
		post.String(),
		retType
}
