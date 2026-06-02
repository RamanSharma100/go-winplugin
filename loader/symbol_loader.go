package loader

import "syscall"

func (d *DLL) Symbol(
	name string,
) *syscall.LazyProc {
	if proc, ok := d.cache[name]; ok {
		return proc
	}

	proc := d.Handle.NewProc(
		name,
	)

	d.cache[name] = proc

	return proc
}
