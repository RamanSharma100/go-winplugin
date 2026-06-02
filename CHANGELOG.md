# Changelog

## v0.2.0 - 2026-06-02

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

## v0.1.1 - 2026-05-26

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
