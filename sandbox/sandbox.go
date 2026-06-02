package sandbox

import (
	"fmt"
	"go/ast"
	"strings"
)

var blockedImports = map[string]bool{
	"os":      true,
	"os/exec": true,
	"syscall": true,
	"unsafe":  true,
}

func ValidateSandbox(
	file *ast.File,
) error {
	for _, imp := range file.Imports {
		name := strings.Trim(
			imp.Path.Value,
			"\"",
		)

		if blockedImports[name] {
			return fmt.Errorf(
				"sandbox violation: %s import not allowed",
				name,
			)
		}
	}

	return nil
}
