package environment

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	msysRoot = `C:\msys64`
	ucrtBin  = `C:\msys64\ucrt64\bin`
	gccPath  = `C:\msys64\ucrt64\bin\gcc.exe`
	bashPath = `C:\msys64\usr\bin\bash.exe`
)

func InstallDependencies(
	env *Environment,
) error {

	// GCC ALREADY AVAILABLE
	if env.GCCInstalled {

		fmt.Println(
			"GCC already installed",
		)

		ensurePath()

		return nil
	}

	fmt.Println(
		"GCC not found",
	)

	// CHECK EXISTING GCC INSTALLATION
	_, err := os.Stat(
		gccPath,
	)

	if err == nil {

		fmt.Println(
			"Existing GCC installation found",
		)

		ensurePath()

		return nil
	}

	fmt.Println(
		"Installing MSYS2...",
	)

	install := exec.Command(
		"winget",
		"install",
		"-e",
		"--id",
		"MSYS2.MSYS2",
	)

	install.Stdout = os.Stdout
	install.Stderr = os.Stderr

	err = install.Run()

	if err != nil {

		fmt.Println(
			"MSYS2 may already be installed",
		)

		fmt.Println(
			"Continuing setup...",
		)
	}

	// VERIFY MSYS2
	_, err = os.Stat(
		bashPath,
	)

	if err != nil {
		return fmt.Errorf(
			"MSYS2 installation not found",
		)
	}

	// CHECK GCC AGAIN
	_, err = os.Stat(
		gccPath,
	)

	if err == nil {

		fmt.Println(
			"GCC already installed inside MSYS2",
		)

		ensurePath()

		return nil
	}

	fmt.Println(
		"Installing GCC...",
	)

	pacman := exec.Command(
		bashPath,
		"-lc",
		"pacman -S --noconfirm mingw-w64-ucrt-x86_64-gcc",
	)

	pacman.Stdout = os.Stdout
	pacman.Stderr = os.Stderr

	err = pacman.Run()

	if err != nil {
		return fmt.Errorf(
			"failed to install gcc: %w",
			err,
		)
	}

	_, err = os.Stat(
		gccPath,
	)

	if err != nil {
		return fmt.Errorf(
			"gcc installed but binary missing",
		)
	}

	ensurePath()

	fmt.Println(
		"GCC installation completed",
	)

	return nil
}

func ensurePath() {
	pathEnv := os.Getenv(
		"PATH",
	)

	if strings.Contains(
		strings.ToLower(pathEnv),
		strings.ToLower(ucrtBin),
	) {

		fmt.Println(
			"GCC already exists in PATH",
		)

		return
	}

	fmt.Println()
	fmt.Println(
		"GCC path not found in PATH",
	)
	fmt.Println()
	fmt.Println(
		"Add this to PATH manually:",
	)

	fmt.Println(
		ucrtBin,
	)

	fmt.Println()
	fmt.Println(
		"Windows Steps:",
	)

	fmt.Println(
		"1. Open Start Menu",
	)

	fmt.Println(
		"2. Search 'Environment Variables'",
	)

	fmt.Println(
		"3. Open 'Edit the system environment variables'",
	)

	fmt.Println(
		"4. Click Environment Variables",
	)

	fmt.Println(
		"5. Edit PATH",
	)

	fmt.Println(
		"6. Add:",
	)

	fmt.Println(
		ucrtBin,
	)

	fmt.Println()
	fmt.Println(
		"Restart terminal after updating PATH",
	)

	fmt.Println()
}
