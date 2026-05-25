package winplugin

import (
	"fmt"
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

	return &Loader{
		rootDir:     rootDir,
		packageName: filepath.Base(rootDir),
		environment: env,
	}, nil
}

func (l *Loader) Build(
	fileName string,
) error {

	source := filepath.Join(
		l.rootDir,
		fileName,
	)

	file, err := compiler.ParseFile(
		source,
	)

	if err != nil {
		return err
	}

	functions := compiler.AnalyzeFunctions(
		file,
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

	l.workspace = workspace

	if err != nil {
		return err
	}

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

	err = builder.CopyFile(
		source,
		filepath.Join(
			pluginDir,
			fileName,
		),
	)

	if err != nil {
		return err
	}

	wrapper := compiler.GenerateWrapper(
		l.packageName,
		file.Name.Name,
		l.packageName+"/plugin",
		functions,
	)

	err = os.WriteFile(
		filepath.Join(
			bridgeDir,
			"wrapper.go",
		),

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

	_, err = os.Stat(
		goMod,
	)

	if err == nil {

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

	_, err = os.Stat(
		goSum,
	)

	if err == nil {

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

	err = builder.BuildDLL(
		bridgeDir,
		filepath.Join(
			l.rootDir,
			artifact,
		),
	)

	if err != nil {
		return err
	}

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
	fmt.Println("Workspace:", l.rootDir)
	dllPath := filepath.Join(
		l.workspace,
		"bridge",
		l.rootDir,
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
	err := l.load(artifactName)
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
	if err != nil && err.Error() != "The operation completed successfully." {
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
