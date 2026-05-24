package executor

import "syscall"

func Call(
	proc *syscall.LazyProc,
	args ...uintptr,
) (uintptr, uintptr, error) {
	return proc.Call(args...)
}
