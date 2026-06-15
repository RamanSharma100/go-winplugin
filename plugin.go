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
	functions   []compiler.Function
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

	absRoot, err := filepath.Abs(rootDir)
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

func (l *Loader) Build(fileName string) error {
	packages, err := compiler.ParsePackage(l.rootDir)
	if err != nil {
		return err
	}

	var parsedPackage *ast.Package
	for _, pkg := range packages {
		parsedPackage = pkg
		break
	}

	if parsedPackage == nil {
		return fmt.Errorf("no package found in %s", l.rootDir)
	}

	functions := compiler.AnalyzePackageFunctions(parsedPackage)

	err = compiler.ValidateSandbox(parsedPackage)
	if err != nil {
		return err
	}

	for _, file := range parsedPackage.Files {
		err = sandbox.ValidateSandbox(file)
		if err != nil {
			return err
		}
	}

	err = compiler.ValidateFunctions(functions)
	if err != nil {
		return err
	}

	l.functions = functions

	workspace, err := builder.CreateWorkspace(l.packageName)
	if err != nil {
		return err
	}

	l.workspace = workspace

	err = EnsureGoMod(workspace, l.packageName)
	if err != nil {
		return err
	}

	pluginDir := filepath.Join(workspace, "plugin")
	bridgeDir := filepath.Join(workspace, "bridge")

	err = os.MkdirAll(pluginDir, 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(bridgeDir, 0755)
	if err != nil {
		return err
	}

	goFiles, err := filepath.Glob(filepath.Join(l.rootDir, "*.go"))
	if err != nil {
		return err
	}

	if len(goFiles) == 0 {
		return fmt.Errorf("no go files found in %s", l.rootDir)
	}

	for _, goFile := range goFiles {
		fileBase := filepath.Base(goFile)
		err = builder.CopyFile(goFile, filepath.Join(pluginDir, fileBase))
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

	wrapperPath := filepath.Join(bridgeDir, "wrapper.go")
	err = os.WriteFile(wrapperPath, []byte(wrapper), 0644)
	if err != nil {
		return err
	}

	goMod := filepath.Join(l.rootDir, "go.mod")
	if _, err := os.Stat(goMod); err == nil {
		err = builder.CopyFile(goMod, filepath.Join(workspace, "go.mod"))
		if err != nil {
			return err
		}
	}

	goSum := filepath.Join(l.rootDir, "go.sum")
	if _, err := os.Stat(goSum); err == nil {
		err = builder.CopyFile(goSum, filepath.Join(workspace, "go.sum"))
		if err != nil {
			return err
		}
	}

	artifact := builder.PlatformLibraryName(l.packageName)
	outputPath := filepath.Join(bridgeDir, artifact)

	err = builder.BuildDLL(bridgeDir, outputPath)
	if err != nil {
		return err
	}

	finalOutput := filepath.Join(l.rootDir, artifact)
	err = builder.CopyFile(outputPath, finalOutput)
	if err != nil {
		return err
	}

	fmt.Println("DLL built successfully at:", finalOutput)

	return nil
}

func (l *Loader) load(artifactName string) error {
	if l.dll != nil {
		return nil
	}

	if l.workspace == "" {
		return fmt.Errorf("workspace not initialized")
	}

	artifact := builder.PlatformLibraryName(artifactName)
	dllPath := filepath.Join(l.workspace, "bridge", artifact)

	fmt.Println("Loading DLL:", dllPath)

	_, err := os.Stat(dllPath)
	if err != nil {
		return fmt.Errorf("dll not found: %s", dllPath)
	}

	dll, err := loader.Load(dllPath)
	if err != nil {
		return err
	}

	l.dll = dll
	return nil
}

func (l *Loader) findFunction(name string) *compiler.Function {
	for _, fn := range l.functions {
		if fn.Name == name {
			return &fn
		}
	}
	return nil
}

func marshalArg(v any) (uintptr, error) {
	payload, err := json.Marshal(v)
	if err != nil {
		return 0, err
	}
	ptr := uintptr(unsafe.Pointer(syscall.StringBytePtr(string(payload))))
	return ptr, nil
}

func stringArg(s string) uintptr {
	return uintptr(unsafe.Pointer(syscall.StringBytePtr(s)))
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

	fn := l.findFunction(function)

	callArgs := []uintptr{}

	for i, arg := range args {
		// Determine the declared param type if available
		paramType := ""
		if fn != nil && i < len(fn.Params) {
			paramType = fn.Params[i].Type
		}

		ptr, err := marshalCallArg(arg, paramType)
		if err != nil {
			return nil, fmt.Errorf("arg %d: %w", i, err)
		}
		callArgs = append(callArgs, ptr)
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

// marshalCallArg converts a Go value to a uintptr for the DLL call.
// paramType is the declared type from the plugin's function signature.
// Rules:
//   - int, int32, int64, uint, uint32, uint64 with param type int/int32/int64 → raw uintptr (C.int)
//   - float32, float64 with param type float32/float64 → raw uintptr bits (C.double/C.float)
//   - bool with param type bool → 0 or 1
//   - string with param type string → C string pointer
//   - anything going to interface{}, map, slice, *struct, struct → JSON C string pointer
//   - unknown param type → JSON C string pointer (safe default)
func marshalCallArg(arg any, paramType string) (uintptr, error) {
	// Nil always marshals to JSON "null"
	if arg == nil {
		return marshalArg(nil)
	}

	switch paramType {

	case "int", "int32", "int64", "uint", "uint32", "uint64":
		switch v := arg.(type) {
		case int:
			return uintptr(v), nil
		case int8:
			return uintptr(v), nil
		case int16:
			return uintptr(v), nil
		case int32:
			return uintptr(v), nil
		case int64:
			return uintptr(v), nil
		case uint:
			return uintptr(v), nil
		case uint8:
			return uintptr(v), nil
		case uint16:
			return uintptr(v), nil
		case uint32:
			return uintptr(v), nil
		case uint64:
			return uintptr(v), nil
		case uintptr:
			return v, nil
		case float32:
			return uintptr(int64(v)), nil
		case float64:
			return uintptr(int64(v)), nil
		default:
			return marshalArg(v)
		}

	case "float32":
		switch v := arg.(type) {
		case float32:
			return uintptr(v), nil
		case float64:
			return uintptr(float32(v)), nil
		case int:
			return uintptr(float32(v)), nil
		default:
			return marshalArg(v)
		}

	case "float64":
		switch v := arg.(type) {
		case float64:
			return uintptr(v), nil
		case float32:
			return uintptr(float64(v)), nil
		case int:
			return uintptr(float64(v)), nil
		default:
			return marshalArg(v)
		}

	case "bool":
		switch v := arg.(type) {
		case bool:
			if v {
				return 1, nil
			}
			return 0, nil
		case int:
			if v != 0 {
				return 1, nil
			}
			return 0, nil
		default:
			return marshalArg(v)
		}

	case "string":
		switch v := arg.(type) {
		case string:
			return stringArg(v), nil
		default:
			s := fmt.Sprintf("%v", v)
			return stringArg(s), nil
		}

	default:
		// interface{}, *struct, struct, map, slice, []byte,
		// unknown, or no type info — always JSON marshal
		return marshalArg(arg)
	}
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
				return nil, fmt.Errorf("%s", errStr)
			}
		}
		delete(v, "error")
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
