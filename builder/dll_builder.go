package builder

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/RamanSharma100/go-winplugin/environment"
)

func PlatformLibraryName(
	name string,
) string {
	switch runtime.GOOS {
	case "windows":
		return name + ".dll"
	case "linux":
		return name + ".so"
	case "darwin":
		return name + ".dylib"
	}

	return name
}

func BuildDLL(
	dir string,
	output string,
) error {

	cmd := exec.Command(
		"go",
		"build",
		"-x",
		"-buildmode=c-shared",
		"-o",
		output,
	)

	cmd.Dir = dir

	environment.PrepareBuildEnvironment(
		cmd,
	)

	out, err := cmd.CombinedOutput()

	fmt.Println(
		string(out),
	)

	if err != nil {
		return err
	}

	return nil
}
