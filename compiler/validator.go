package compiler

import "fmt"

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

		if fn.ReturnType != "" && fn.ReturnType != "void" &&
			fn.ReturnType != "int" &&
			fn.ReturnType != "string" &&
			fn.ReturnType != "float" &&
			fn.ReturnType != "bool" {
			return fmt.Errorf("unsupported return type: %s in %s", fn.ReturnType, fn.Name)
		}

		for _, p := range fn.Params {
			switch p.Type {
			case "int", "string", "float", "bool":
				continue
			default:
				return fmt.Errorf("unsupported parameter type: %s in %s", p.Type, fn.Name)
			}
		}
	}

	return nil
}
