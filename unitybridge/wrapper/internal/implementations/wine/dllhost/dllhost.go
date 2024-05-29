//go:build windows && amd64

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"syscall"
	"unsafe"

	"github.com/brunoga/robomaster/support/logger"
	"github.com/brunoga/robomaster/unitybridge/wrapper"
)

var (
	// Command line flags.
	readFd   = flag.Int("read-fd", -1, "file descriptor to read from")
	writeFd  = flag.Int("write-fd", -1, "file descriptor to write to")
	eventFd  = flag.Int("event-fd", -1, "file descriptor to write events to")
	logLevel = flag.Int("loglevel", int(slog.LevelError), "log level")

	callbackHandler *CallbackHandler

	ub wrapper.UnityBridge
)

func main() {
	flag.Parse()

	if *readFd < 0 || *writeFd < 0 || *eventFd < 0 {
		fmt.Fprintln(os.Stderr, "Flags -read-fd, -write-fd  and -events-fs "+
			"must be provided and non-negative")
		os.Exit(1)
	}

	files, err := fdsToFiles([]int{*readFd, *writeFd, *eventFd})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error converting file descriptors to "+
			"files: %s\n", err)
		os.Exit(1)
	}

	level := slog.Level(*logLevel)
	l := logger.New(level)

	l.Info("Starting Unity Bridge DLL Host")

	ub = wrapper.Get(l)

	callbackHandler = &CallbackHandler{eventFile: files[2]}

	err = loop(files[0], files[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error in loop: %s\n", err)
		os.Exit(1)
	}
}

func fdsToFiles(fds []int) ([]*os.File, error) {
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

	files := make([]*os.File, len(fds))
	for i, fd := range fds {
		file := fdToFile(wineServerFdToHandleProc, uintptr(fd),
			syscall.GENERIC_READ|syscall.GENERIC_WRITE,
			fmt.Sprintf("dllhost%d", i))
		files[i] = file
	}

	return files, nil
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

func loop(readFile, writeFile *os.File) error {
	headerBuf := make([]byte, 1+4)

	for {
		if _, err := io.ReadFull(readFile, headerBuf); err != nil {
			if err != io.EOF {
				return err
			} else {
				break
			}
		}

		function := headerBuf[0]

		length := binary.BigEndian.Uint32(headerBuf[1:5])

		var data []byte
		if length != 0 {
			data = make([]byte, length)
			_, err := io.ReadFull(readFile, data)
			if err != nil {
				return err
			}
		}

		process(writeFile, function, data)
	}

	return nil
}

func process(writeFile *os.File, function byte, data []byte) {
	var b bytes.Buffer

	b.WriteByte(function)

	switch function {
	case 0x00:
		runCreateUnityBridge(data, &b)
	case 0x01:
		runDestroyUnityBridge(data, &b)
	case 0x02:
		runInitializeUnityBridge(data, &b)
	case 0x03:
		runUnitializeUnityBridge(data, &b)
	case 0x04:
		runUnitySendEvent(data, &b)
	case 0x05:
		runUnitySendEventWithString(data, &b)
	case 0x06:
		runUnitySendEventWithNumber(data, &b)
	case 0x07:
		runUnitySetEventCallback(data, &b)
	case 0x08:
		runGetSecurityKeyByKeyChainIndex(data, &b)
	}

	_, err := b.WriteTo(writeFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %s\n", err)
		return
	}
}

func runCreateUnityBridge(data []byte, b *bytes.Buffer) {
	debuggable := data[0] != 0
	nameLength := binary.BigEndian.Uint32(data[1:5])
	name := string(data[5 : 5+nameLength])

	// No need to parse the logPath size because it will be whatever is left of
	// the buffer. So we just make sure we skip the size.
	logPath := string(data[5+nameLength+4:])

	ub.Create(name, debuggable, logPath)

	// Write data size.
	writeSize(b, 0)
}

func runDestroyUnityBridge(data []byte, b *bytes.Buffer) {
	ub.Destroy()

	// Write data size.
	writeSize(b, 0)
}

func runInitializeUnityBridge(data []byte, b *bytes.Buffer) {
	res := ub.Initialize()

	// Write data size.
	writeSize(b, 1)

	if res {
		b.WriteByte(0x01)
	} else {
		b.WriteByte(0x00)
	}
}

func runUnitializeUnityBridge(data []byte, b *bytes.Buffer) {
	ub.Uninitialize()

	// Write data size.
	writeSize(b, 0)
}

func runUnitySendEvent(data []byte, b *bytes.Buffer) {
	eventCode := binary.BigEndian.Uint64(data[0:8])
	tag := binary.BigEndian.Uint64(data[8:16])
	outputLen := binary.BigEndian.Uint32(data[16:20])

	output := make([]byte, outputLen)

	ub.SendEvent(eventCode, output, tag)

	// Write data size.
	writeSize(b, outputLen)

	b.Write(output)
}

func runUnitySendEventWithString(data []byte, b *bytes.Buffer) {
	eventCode := binary.BigEndian.Uint64(data[0:8])
	tag := binary.BigEndian.Uint64(data[8:16])
	length := binary.BigEndian.Uint32(data[16:20])
	data2 := string(data[20 : 20+length])

	ub.SendEventWithString(eventCode, data2, tag)

	// Write data size.
	writeSize(b, 0)
}

func runUnitySendEventWithNumber(data []byte, b *bytes.Buffer) {
	eventCode := binary.BigEndian.Uint64(data[0:8])
	tag := binary.BigEndian.Uint64(data[8:16])
	data2 := binary.BigEndian.Uint64(data[16:24])

	ub.SendEventWithNumber(eventCode, data2, tag)

	// Write data size.
	writeSize(b, 0)
}

func runUnitySetEventCallback(data []byte, b *bytes.Buffer) {
	eventTypeCode := binary.BigEndian.Uint64(data[0:8])
	add := data[8] != 0

	if add {
		ub.SetEventCallback(eventTypeCode, callbackHandler.HandleCallback)
	} else {
		ub.SetEventCallback(eventTypeCode, nil)
	}

	// Write data size.
	writeSize(b, 0)
}

func runGetSecurityKeyByKeyChainIndex(data []byte, b *bytes.Buffer) {
	index := binary.BigEndian.Uint64(data[0:8])

	key := ub.GetSecurityKeyByKeyChainIndex(int(index))

	// Write data size.
	writeSize(b, uint32(len(key)))

	b.WriteString(key)
}

func writeSize(w io.Writer, size uint32) error {
	return binary.Write(w, binary.BigEndian, size)
}
