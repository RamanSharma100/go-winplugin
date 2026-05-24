package loader

import (
	"fmt"
	"syscall"
	"unsafe"
)

const ucrtPath = `C:\msys64\ucrt64\bin`

type DLL struct {
	Handle *syscall.LazyDLL
}

func Load(
	path string,
) (*DLL, error) {
	kernel32 := syscall.NewLazyDLL(
		"kernel32.dll",
	)
	setDllDirectory := kernel32.NewProc(
		"SetDllDirectoryW",
	)

	ptr, err := syscall.UTF16PtrFromString(
		ucrtPath,
	)

	if err == nil {
		setDllDirectory.Call(
			uintptr(
				unsafe.Pointer(ptr),
			),
		)
	}

	dll := syscall.NewLazyDLL(
		path,
	)

	err = dll.Load()

	if err != nil {
		return nil, fmt.Errorf(
			"Failed to load %s: %w",
			path,
			err,
		)
	}

	return &DLL{
		Handle: dll,
	}, nil
}
