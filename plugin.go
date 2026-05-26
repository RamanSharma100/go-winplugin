package winplugin

import (
	"fmt"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/RamanSharma100/go-winplugin/builder"
	"github.com/RamanSharma100/go-winplugin/compiler"
	"github.com/RamanSharma100/go-winplugin/environment"
	"github.com/RamanSharma100/go-winplugin/executor"
	"github.com/RamanSharma100/go-winplugin/loader"
)

type Loader struct {
	rootDir     string
	packageName string
	workspace   string
	environment *environment.Environment
	dll         *loader.DLL
}

func NewLoader(
	rootDir string,
) (*Loader, error) {
	env, err := environment.CheckEnvironment()

	if err != nil {
		return nil, err
	}

	err = environment.InstallDependencies(
		env,
	)

	if err != nil {
		return nil, err
	}

	absRoot, err := filepath.Abs(
		rootDir,
	)

	if err != nil {
		return nil, err
	}

	tempDir := filepath.Join(os.TempDir(), "go-winplugin")

	defer os.RemoveAll(tempDir)

	return &Loader{
		rootDir:     absRoot,
		packageName: filepath.Base(absRoot),
		environment: env,
	}, nil
}

func (l *Loader) Build(
	fileName string,
) error {
	packages, err := compiler.ParsePackage(
		l.rootDir,
	)

	if err != nil {
		return err
	}

	var parsedPackage *ast.Package

	for _, pkg := range packages {
		parsedPackage = pkg
		break
	}

	if parsedPackage == nil {
		return fmt.Errorf(
			"no package found in %s",
			l.rootDir,
		)
	}

	functions := compiler.AnalyzePackageFunctions(
		parsedPackage,
	)

	err = compiler.ValidateFunctions(
		functions,
	)

	if err != nil {
		return err
	}

	workspace, err := builder.CreateWorkspace(
		l.packageName,
	)

	if err != nil {
		return err
	}

	l.workspace = workspace

	err = EnsureGoMod(
		workspace,
		l.packageName,
	)

	if err != nil {
		return err
	}

	pluginDir := filepath.Join(
		workspace,
		"plugin",
	)

	bridgeDir := filepath.Join(
		workspace,
		"bridge",
	)

	err = os.MkdirAll(
		pluginDir,
		0755,
	)

	if err != nil {
		return err
	}

	err = os.MkdirAll(
		bridgeDir,
		0755,
	)

	if err != nil {
		return err
	}

	goFiles, err := filepath.Glob(
		filepath.Join(
			l.rootDir,
			"*.go",
		),
	)

	if err != nil {
		return err
	}

	if len(goFiles) == 0 {
		return fmt.Errorf(
			"no go files found in %s",
			l.rootDir,
		)
	}

	for _, goFile := range goFiles {
		fileBase := filepath.Base(
			goFile,
		)

		err = builder.CopyFile(
			goFile,
			filepath.Join(
				pluginDir,
				fileBase,
			),
		)

		if err != nil {
			return err
		}
	}

	packageName := parsedPackage.Name

	wrapper := compiler.GenerateWrapper(
		l.packageName,
		packageName,
		l.packageName+"/plugin",
		functions,
	)

	wrapperPath := filepath.Join(
		bridgeDir,
		"wrapper.go",
	)

	err = os.WriteFile(
		wrapperPath,
		[]byte(wrapper),
		0644,
	)

	if err != nil {
		return err
	}

	goMod := filepath.Join(
		l.rootDir,
		"go.mod",
	)

	if _, err := os.Stat(
		goMod,
	); err == nil {
		err = builder.CopyFile(
			goMod,
			filepath.Join(
				workspace,
				"go.mod",
			),
		)

		if err != nil {
			return err
		}
	}

	goSum := filepath.Join(
		l.rootDir,
		"go.sum",
	)

	if _, err := os.Stat(
		goSum,
	); err == nil {
		err = builder.CopyFile(
			goSum,
			filepath.Join(
				workspace,
				"go.sum",
			),
		)

		if err != nil {
			return err
		}
	}

	artifact := builder.PlatformLibraryName(
		l.packageName,
	)

	outputPath := filepath.Join(
		bridgeDir,
		artifact,
	)

	err = builder.BuildDLL(
		bridgeDir,
		outputPath,
	)

	if err != nil {
		return err
	}

	finalOutput := filepath.Join(
		l.rootDir,
		artifact,
	)

	err = builder.CopyFile(
		outputPath,
		finalOutput,
	)

	if err != nil {
		return err
	}

	fmt.Println(
		"DLL built successfully at:",
		finalOutput,
	)

	return nil
}

func (l *Loader) load(
	artifactName string,
) error {
	if l.dll != nil {
		return nil
	}

	if l.workspace == "" {
		return fmt.Errorf(
			"workspace not initialized",
		)
	}

	artifact := builder.PlatformLibraryName(
		artifactName,
	)

	dllPath := filepath.Join(
		l.workspace,
		"bridge",
		artifact,
	)

	fmt.Println(
		"Loading DLL:",
		dllPath,
	)

	_, err := os.Stat(
		dllPath,
	)

	if err != nil {
		return fmt.Errorf(
			"dll not found: %s",
			dllPath,
		)
	}

	dll, err := loader.Load(
		dllPath,
	)

	if err != nil {
		return err
	}

	l.dll = dll

	return nil
}

func (l *Loader) Call(
	artifactName string,
	function string,
	args ...uintptr,
) (uintptr, error) {
	err := l.load(
		artifactName,
	)

	if err != nil {
		return 0, err
	}

	symbol := compiler.MangleSymbol(
		l.packageName,
		function,
	)

	proc := l.dll.Symbol(
		symbol,
	)

	result, _, err := executor.Call(
		proc,
		args...,
	)

	if err != nil &&
		err.Error() != "The operation completed successfully." {
		return 0, err
	}

	return result, nil
}

func EnsureGoMod(
	dir string,
	moduleName string,
) error {
	goModPath := filepath.Join(
		dir,
		"go.mod",
	)

	if _, err := os.Stat(
		goModPath,
	); err == nil {
		return nil
	}

	cmd := exec.Command(
		"go",
		"mod",
		"init",
		moduleName,
	)

	cmd.Dir = dir

	out, err := cmd.CombinedOutput()

	if err != nil {
		return fmt.Errorf(
			"go mod init failed: %s, err: %w",
			string(out),
			err,
		)
	}

	return nil
}
