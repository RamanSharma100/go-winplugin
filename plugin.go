package winplugin

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/RamanSharma100/go-winplugin/builder"
	"github.com/RamanSharma100/go-winplugin/compiler"
	"github.com/RamanSharma100/go-winplugin/environment"
	"github.com/RamanSharma100/go-winplugin/executor"
	"github.com/RamanSharma100/go-winplugin/loader"
	"github.com/RamanSharma100/go-winplugin/sandbox"
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

	err = compiler.ValidateSandbox(
		parsedPackage,
	)

	if err != nil {
		return err
	}

	for _, file := range parsedPackage.Files {
		err = sandbox.ValidateSandbox(file)

		if err != nil {
			return err
		}
	}

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
	args ...any,
) (any, error) {
	err := l.load(artifactName)
	if err != nil {
		return nil, err
	}

	symbol := compiler.MangleSymbol(l.packageName, function)
	proc := l.dll.Symbol(symbol)

	callArgs := []uintptr{}

	for _, arg := range args {
		switch v := arg.(type) {
		case uintptr:
			callArgs = append(callArgs, v)
		case int:
			callArgs = append(callArgs, uintptr(v))
		case bool:
			if v {
				callArgs = append(callArgs, 1)
			} else {
				callArgs = append(callArgs, 0)
			}
		case string:
			ptr := uintptr(unsafe.Pointer(syscall.StringBytePtr(v)))
			callArgs = append(callArgs, ptr)
		default:
			payload, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			ptr := uintptr(unsafe.Pointer(syscall.StringBytePtr(string(payload))))
			callArgs = append(callArgs, ptr)
		}
	}

	result, _, err := executor.Call(proc, callArgs...)

	if err != nil && err.Error() != "The operation completed successfully." {
		return nil, err
	}

	if result == 0 {
		return nil, nil
	}

	str, err := cStringToGoString(result)
	if err != nil {
		return nil, fmt.Errorf("failed to read DLL return value: %w", err)
	}

	if str == "" {
		return nil, nil
	}

	return parseEnvelope(strings.TrimSpace(str))
}

func parseEnvelope(str string) (any, error) {
	var envelope struct {
		Type  string          `json:"type"`
		Value json.RawMessage `json:"value"`
	}

	if err := json.Unmarshal([]byte(str), &envelope); err != nil {
		return str, nil
	}

	switch envelope.Type {

	case "scalar":
		var n float64
		if err := json.Unmarshal(envelope.Value, &n); err != nil {
			return nil, err
		}
		if n == float64(int64(n)) {
			return int64(n), nil
		}
		return n, nil

	case "string":
		var s string
		if err := json.Unmarshal(envelope.Value, &s); err != nil {
			return nil, err
		}
		return s, nil

	case "bool":
		var b bool
		if err := json.Unmarshal(envelope.Value, &b); err != nil {
			return nil, err
		}
		return b, nil

	case "error":
		if string(envelope.Value) == "null" {
			return nil, nil
		}
		var s string
		if err := json.Unmarshal(envelope.Value, &s); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", s)

	case "bytes":
		var b []byte
		if err := json.Unmarshal(envelope.Value, &b); err != nil {
			return nil, err
		}
		return b, nil

	case "json":
		var v any
		if err := json.Unmarshal(envelope.Value, &v); err != nil {
			return nil, err
		}
		return v, nil

	case "multi":
		var raw []byte
		if err := json.Unmarshal(envelope.Value, &raw); err != nil {
			return nil, err
		}
		var v map[string]any
		if err := json.Unmarshal(raw, &v); err != nil {
			return nil, err
		}
		if errVal, ok := v["error"]; ok && errVal != nil {
			if errStr, ok := errVal.(string); ok {
				delete(v, "error")
				if len(v) == 0 {
					return nil, fmt.Errorf("%s", errStr)
				}
				return nil, fmt.Errorf("%s", errStr)
			}
		}
		if len(v) == 1 {
			for _, val := range v {
				return val, nil
			}
		}
		return v, nil

	default:
		return str, nil
	}
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

func cStringToGoString(ptr uintptr) (string, error) {
	if ptr == 0 {
		return "", nil
	}

	const maxLen = 10 * 1024 * 1024

	var buf []byte
	for i := 0; i < maxLen; i++ {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		if b == 0 {
			return string(buf), nil
		}
		buf = append(buf, b)
	}

	return "", fmt.Errorf("cStringToGoString: string exceeds max length of %d bytes", maxLen)
}
