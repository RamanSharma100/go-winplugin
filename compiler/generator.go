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

		if len(fn.ReturnTypes) == 0 {

			sb.WriteString("\t")
			sb.WriteString(call)
			sb.WriteString("\n")

		} else {

			sb.WriteString("\t")

			for i := range fn.ReturnTypes {

				if i > 0 {
					sb.WriteString(", ")
				}

				sb.WriteString("result")
				sb.WriteString(string(rune('0' + i)))
			}

			sb.WriteString(" := ")
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

			pre.WriteString("\t")
			pre.WriteString(goName)
			pre.WriteString(" := C.GoString(")
			pre.WriteString(cName)
			pre.WriteString(")\n")

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

			if p.isInterface {
				pre.WriteString("\tvar ")
				pre.WriteString(goName)
				pre.WriteString(" interface{}\n")

				pre.WriteString("\tjson.Unmarshal([]byte(C.GoString(")
				pre.WriteString(cName)
				pre.WriteString(")), &")
				pre.WriteString(goName)
				pre.WriteString(")\n")

				goArgs = append(goArgs, goName)

				continue
			}

			if isSlice(p.Type) {

				pre.WriteString("\tvar ")
				pre.WriteString(goName)
				pre.WriteString(" ")
				pre.WriteString(p.Type)
				pre.WriteString("\n")

				pre.WriteString("\tjson.Unmarshal([]byte(C.GoString(")
				pre.WriteString(cName)
				pre.WriteString(")), &")
				pre.WriteString(goName)
				pre.WriteString(")\n")

				goArgs = append(goArgs, goName)

				continue
			}

			if isMap(p.Type) {

				pre.WriteString("\tvar ")
				pre.WriteString(goName)
				pre.WriteString(" ")
				pre.WriteString(p.Type)
				pre.WriteString("\n")

				pre.WriteString("\tjson.Unmarshal([]byte(C.GoString(")
				pre.WriteString(cName)
				pre.WriteString(")), &")
				pre.WriteString(goName)
				pre.WriteString(")\n")

				goArgs = append(goArgs, goName)

				continue
			}

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

	retType, postCode := buildReturnBridge(fn)
	post.WriteString(postCode)

	return strings.Join(cParams, ", "),
		strings.Join(goArgs, ", "),
		pre.String(),
		post.String(),
		retType
}

func buildReturnBridge(
	fn Function,
) (
	string,
	string,
) {
	var post strings.Builder

	if len(fn.ReturnTypes) == 0 {
		return "", ""
	}

	if len(fn.ReturnTypes) == 1 {

		switch fn.ReturnTypes[0] {

		case "int", "int32":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"scalar\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "int64":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"scalar\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "uint64":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"scalar\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "float32", "float64":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"scalar\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "string":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"string\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "bool":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"bool\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "error":
			post.WriteString("\tvar errVal interface{}\n")
			post.WriteString("\tif result0 != nil {\n")
			post.WriteString("\t\terrVal = result0.Error()\n")
			post.WriteString("\t}\n")
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"error\", \"value\": errVal})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		case "[]byte":
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"bytes\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()

		default:
			post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"json\", \"value\": result0})\n")
			post.WriteString("\treturn C.CString(string(envelope))\n")
			return "*C.char", post.String()
		}
	}

	for i, returnType := range fn.ReturnTypes {
		if returnType == "error" {
			varName := "errStr" + string(rune('0'+i))
			resultName := "result" + string(rune('0'+i))
			post.WriteString("\tvar " + varName + " interface{}\n")
			post.WriteString("\tif " + resultName + " != nil {\n")
			post.WriteString("\t\t" + varName + " = " + resultName + ".Error()\n")
			post.WriteString("\t}\n")
		}
	}

	post.WriteString("\tpayload, _ := json.Marshal(map[string]interface{}{")

	for i, returnType := range fn.ReturnTypes {
		if i > 0 {
			post.WriteString(", ")
		}
		if returnType == "error" {
			post.WriteString("\"error\": errStr" + string(rune('0'+i)))
		} else {
			post.WriteString("\"result" + string(rune('0'+i)) + "\": result" + string(rune('0'+i)))
		}
	}

	post.WriteString("})\n")
	post.WriteString("\tenvelope, _ := json.Marshal(map[string]interface{}{\"type\": \"multi\", \"value\": payload})\n")
	post.WriteString("\treturn C.CString(string(envelope))\n")

	return "*C.char", post.String()
}

func isSlice(t string) bool {
	return strings.HasPrefix(t, "[]")
}

func isMap(t string) bool {
	return strings.HasPrefix(t, "map[")
}
