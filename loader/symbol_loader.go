package loader

import "syscall"

func (d *DLL) Symbol(
	name string,
) *syscall.LazyProc {
	return d.Handle.NewProc(name)
}
