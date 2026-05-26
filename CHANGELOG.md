# Changelog

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

## Planned

- Struct support
- Reflection support
- Hot reload
- Linux/macOS support
- Symbol cache
- Runtime validation
- Memory cleanup improvements

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
