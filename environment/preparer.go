package environment

import (
	"os"
	"os/exec"
)

func PrepareBuildEnvironment(
	cmd *exec.Cmd,
) {
	msysGCC := `C:\msys64\ucrt64\bin`

	path := msysGCC + ";" + os.Getenv("PATH")

	cmd.Env = append(
		os.Environ(),
		"CGO_ENABLED=1",
		"CC=gcc",
		"PATH="+path,
	)
}
