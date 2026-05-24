package builder

import (
	"os"
	"path/filepath"
)

func CreateWorkspace(
	name string,
) (
	string,
	error,
) {
	dir := filepath.Join(
		os.TempDir(),
		"go-winplugin",
		name,
	)

	err := os.MkdirAll(
		dir,
		0755,
	)
	if err != nil {
		return "", err
	}

	return dir, nil
}
