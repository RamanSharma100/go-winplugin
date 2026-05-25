package builder

import (
	"os"
	"path/filepath"

	"github.com/google/uuid"
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
		name+"-"+uuid.NewString(),
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
