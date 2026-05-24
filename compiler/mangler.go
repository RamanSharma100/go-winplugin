package compiler

func MangleSymbol(
	packageName string,
	functionName string,
) string {
	return "go_plugin_" + packageName + "_" + functionName
}
