# Changelog

All notable changes to go-winplugin are documented here.

## [v0.3.0] - 2026-06-14

### Added

- Interface parameter support — `interface{}`, `[]interface{}`, `map[string]any` can now be passed to plugin functions and are marshaled via JSON across the DLL boundary
- Multi-return value support — functions returning multiple values (e.g. `(*User, error)`) are fully supported; values are returned as `map[string]any` and errors are propagated as Go `error` values
- Typed envelope return system — every DLL return is now wrapped in a typed envelope `{"type": "...", "value": ...}` eliminating all pointer/scalar ambiguity that existed in v0.2.0
- `[]byte` return type support
- `error`-only return type support with nil/non-nil propagation
- Struct and `*struct` parameter marshaling via JSON
- `CreateUserWithError`, `IsPositive`, `Divide`, `ReadData`, `Validate` added to example plugin covering all new type scenarios
- Expanded compiler test suite — envelope wrapping verified for all types, multi-return AST parsing, interface marshaling codegen, slice/map parameter parsing
- Expanded integration test suite — zero values, negative numbers, large integers, empty strings, nil interfaces, bool true/false, float64 fractional, void nil return, struct pointer param, multi-return with and without error

### Fixed

- CGO preprocessor crash (`expected ';', found 'return'`) on multi-return functions caused by inline anonymous functions (IIFEs) in generated wrapper code — replaced with hoisted variables
- `Call` return type changed from `uintptr` to `any` — raw pointer values are no longer returned to callers
- Multi-return error detection no longer requires the error to be the only key in the response map

### Changed

- `buildReturnBridge` now emits typed envelopes for all return types instead of raw C scalars or untyped JSON
- `Call` now reads all results as C strings and decodes via `parseEnvelope` — the `0xFFFF` pointer threshold heuristic is removed
- `looksLikeError` heuristic removed — error detection is now fully deterministic via envelope type field

---

## [v0.2.0] - 2026-06-02

### Added

- Multi-file plugin support
- Package-wide function analysis
- Struct type detection
- Struct pointer detection
- Plugin sandbox validation
- Symbol caching
- Package parser improvements
- Sandbox security checks for dangerous imports
- Additional compiler test coverage
- Additional integration test coverage

### Changed

- Build system now copies all Go files from plugin packages
- Package analysis now scans entire packages instead of a single source file
- DLL symbol resolution now uses runtime caching
- Improved workspace generation workflow
- Improved wrapper generation pipeline
- Improved plugin compilation reliability

### Fixed

- Cross-file function resolution during plugin compilation
- Temporary workspace DLL path handling
- Runtime DLL loading path issues
- Wrapper generation inconsistencies
- Exported function discovery across multiple files
- Symbol lookup stability improvements

### Security

- Added plugin sandbox validation
- Blocked dangerous imports:
  - `os`
  - `os/exec`
  - `syscall`
  - `unsafe`

### Testing

Added tests for:

- Struct parsing
- Struct validation
- Sandbox validation
- Multi-file plugin compilation
- Symbol generation consistency
- Runtime plugin execution

### Notes

Struct support in v0.2.0 includes:

- Struct detection
- Struct pointer detection
- Compiler validation
- Wrapper generation compatibility

Struct marshalling across DLL boundaries is planned for a future release.

---

## [v0.1.1] - 2026-05-26

### Added

- Windows CI support via `windows-latest` GitHub Actions runner
- Explicit Go test workflow separation
- Environment variable support for `GOOS=windows` and `GOARCH=amd64`

### Changed

- Improved CI pipeline structure (separated build/test jobs conceptually)
- Updated Go version handling to support `1.26.2`

### Fixed

- Go module version format issues (`go.mod` validation)
- Call of function from another file in plugin
- Temporary folder handling improvements (preparation for cleanup support)

### Notes

- This release improves stability for Windows plugin development workflow

---

## v0.1.0 - 2026-05-25

### Added

- Windows runtime plugin system
- Dynamic DLL loading
- CGO wrapper generation
- Runtime symbol execution
- AST function parsing
- Automatic wrapper generation
- Automatic MSYS2 setup
- Automatic GCC detection
- PATH validation
- Temporary isolated workspaces
- Primitive type support
- Integration testing
