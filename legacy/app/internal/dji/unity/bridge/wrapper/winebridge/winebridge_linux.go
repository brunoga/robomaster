package winebridge

import (
	"debug/pe"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

type WineBridge struct {
	winePath       string
	windowsExePath string
	files          []*os.File
	pid            int
}

func New(windowsExePath string, files ...*os.File) (*WineBridge, error) {
	// Check Wine is available.
	winePath, err := exec.LookPath("wine")
	if err != nil {
		return nil, fmt.Errorf("wine not found in path: %w", err)
	}

	// Check Windows executable.
	peFile, err := pe.Open(windowsExePath)
	if err != nil {
		return nil, fmt.Errorf("%s does not look like a windows "+
			"executable: %w", windowsExePath, err)
	}
	peFile.Close()

	return &WineBridge{
		winePath,
		windowsExePath,
		files,
		-1,
	}, nil
}

func (w *WineBridge) Start() error {
	argv := []string{
		w.winePath,
		w.windowsExePath,
	}

	for _, file := range w.files {
		disableCloseOnExec(file)
		argv = append(argv, fmt.Sprintf("%d", getFd(file)))
	}

	pid, err := syscall.ForkExec(w.winePath, argv,
		&syscall.ProcAttr{
			Files: []uintptr{
				getFd(os.Stdin),
				getFd(os.Stdout),
				getFd(os.Stderr),
			},
			Env: []string{"WINEDEBUG=-all"},
			Sys: &syscall.SysProcAttr{
				Foreground: false,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error executing windows program: %w", err)
	}

	w.pid = pid

	for _, file := range w.files {
		file.Close()
	}

	return nil
}

func (w *WineBridge) Stop() error {
	err := syscall.Kill(w.pid, syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("error terminating child process: %w", err)
	}

	return nil
}

func getFd(file *os.File) uintptr {
	rawConn, err := file.SyscallConn()
	if err != nil {
		panic(fmt.Sprintf("Error getting raw connection "+
			"for file: %s", err))
	}

	var fileFd uintptr
	err = rawConn.Control(func(fd uintptr) {
		fileFd = fd
	})
	if err != nil {
		panic(fmt.Sprintf("Error controlling raw "+
			"connection: %s", err))
	}

	return fileFd
}

func disableCloseOnExec(file *os.File) {
	_, _, err := syscall.Syscall(syscall.SYS_FCNTL, getFd(file),
		syscall.F_SETFD, 0)
	if err != syscall.Errno(0) {
		panic(fmt.Sprintf("Error disabling close on exec: %s", err))
	}
}
