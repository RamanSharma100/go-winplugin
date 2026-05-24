package environment

import (
	"fmt"
	"os/exec"
	"runtime"
)

type Environment struct {
	GoInstalled  bool
	GCCInstalled bool
}

func CheckEnvironment() (
	*Environment,
	error,
) {
	if runtime.GOOS != "windows" {
		return nil, fmt.Errorf(
			"currently windows only",
		)
	}

	env := &Environment{}

	goCmd := exec.Command(
		"go",
		"version",
	)
	if err := goCmd.Run(); err == nil {
		env.GoInstalled = true
	}

	gccCmd := exec.Command(
		"gcc",
		"--version",
	)
	if err := gccCmd.Run(); err == nil {
		env.GCCInstalled = true
	}

	if !env.GoInstalled {
		return nil, fmt.Errorf(
			"go is not installed",
		)
	}

	return env, nil
}
