package winebridge

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"
)

type WineBridge struct {
	files []*os.File
}

func New(numFds int, fdArgs []string) (*WineBridge, error) {
	ntDll, err := syscall.LoadDLL("ntdll.dll")
	if err != nil {
		return nil, fmt.Errorf("error loading ntdll.dll: %s", err)
	}
	defer ntDll.Release()

	wineServerFdToHandleProc, err := ntDll.FindProc(
		"wine_server_fd_to_handle")
	if err != nil {
		return nil, fmt.Errorf(
			"error finding wine_server_fd_to_handle: %w", err)
	}

	if len(fdArgs) != numFds {
		return nil, fmt.Errorf("wrong number of file descriptors: "+
			"expected %d, got %d", numFds, len(fdArgs))
	}

	files := make([]*os.File, len(fdArgs))
	for i, fdArg := range fdArgs {
		fd, err := strconv.Atoi(fdArg)
		if err != nil {
			return nil, fmt.Errorf("error converting file "+
				"descriptor to integer: %s", err)
		}

		file := fdToFile(wineServerFdToHandleProc, uintptr(fd),
			syscall.GENERIC_READ|syscall.GENERIC_WRITE,
			fmt.Sprintf("linux%d", i))
		files[i] = file
	}

	return &WineBridge{
		files,
	}, nil
}

func (w *WineBridge) File(i int) *os.File {
	return w.files[i]
}

func fdToFile(proc *syscall.Proc, fd uintptr, flags uintptr,
	name string) *os.File {
	var fdHandle uintptr
	ntStatus, _, _ := proc.Call(fd, flags|syscall.SYNCHRONIZE,
		2 /*OBJ_INHERIT*/, uintptr(unsafe.Pointer(&fdHandle)))
	if ntStatus != 0 {
		panic(ntStatus)
	}

	return os.NewFile(fdHandle, name)
}
