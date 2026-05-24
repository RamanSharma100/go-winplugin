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
	environment *environment.Environment
	dll         *loader.DLL
}

func NewLoader(rootDir string) (*Loader, error) {
	env, err := environment.CheckEnvironment()
	if err != nil {
		return nil, err
	}

	err = environment.InstallDependencies(env)
	if err != nil {
		return nil, err
	}

	return &Loader{
		rootDir:     rootDir,
		packageName: filepath.Base(rootDir),
		environment: env,
	}, nil
}

func (l *Loader) Build(fileName string) error {
	source := filepath.Join(
		l.rootDir,
		fileName,
	)

	file, err := compiler.ParseFile(source)
	if err != nil {
		return err
	}

	functions := compiler.AnalyzeFunctions(file)

	err = compiler.ValidateFunctions(functions)
	if err != nil {
		return err
	}

	workspace, err := builder.CreateWorkspace(
		l.packageName,
	)

	if err != nil {
		return err
	}

	wrapper := compiler.GenerateWrapper(
		l.packageName,
		functions,
	)

	wrapperPath := filepath.Join(workspace, "wrapper.go")

	err = os.WriteFile(wrapperPath, []byte(wrapper), 0644)
	if err != nil {
		return err
	}

	err = EnsureGoMod(workspace, l.packageName)
	if err != nil {
		return err
	}

	err = builder.CopyDirectory(
		l.rootDir,
		workspace,
	)

	if err != nil {
		return err
	}

	goMod := filepath.Join(
		l.rootDir,
		"go.mod",
	)

	workspaceMod := filepath.Join(
		workspace,
		"go.mod",
	)

	_, err = os.Stat(goMod)

	if err == nil {
		err = builder.CopyFile(
			goMod,
			workspaceMod,
		)

		if err != nil {
			return err
		}
	}

	goSum := filepath.Join(
		l.rootDir,
		"go.sum",
	)

	workspaceSum := filepath.Join(
		workspace,
		"go.sum",
	)

	_, err = os.Stat(goSum)

	if err == nil {

		err = builder.CopyFile(
			goSum,
			workspaceSum,
		)

		if err != nil {
			return err
		}
	}

	err = builder.BuildDLL(
		workspace,
		filepath.Join(workspace, builder.PlatformLibraryName(l.packageName)),
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

	artifact := builder.PlatformLibraryName(
		artifactName,
	)

	path := filepath.Join(
		l.rootDir,
		artifact,
	)

	dll, err := loader.Load(path)
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
) error {

	err := l.load(artifactName)
	if err != nil {
		return err
	}

	symbol := compiler.MangleSymbol(
		l.packageName,
		function,
	)

	proc := l.dll.Symbol(symbol)

	_, _, err = executor.Call(
		proc,
		args...,
	)

	return err
}

func EnsureGoMod(dir string, moduleName string) error {
	goModPath := filepath.Join(dir, "go.mod")

	if _, err := os.Stat(goModPath); err == nil {
		return nil
	}

	cmd := exec.Command("go", "mod", "init", moduleName)
	cmd.Dir = dir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod init failed: %s, err: %w", string(out), err)
	}

	return nil
}
