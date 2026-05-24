package compiler

import (
	"path/filepath"
	"strings"
)

func PluginArtifactName(
	fileName string,
) string {
	base := filepath.Base(fileName)

	return strings.TrimSuffix(
		base,
		filepath.Ext(base),
	)
}
