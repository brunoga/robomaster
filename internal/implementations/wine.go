//go:build linux && amd64

package implementations

import (
	"bytes"
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"github.com/brunoga/unitybridge/event"

	internal_event "github.com/brunoga/unitybridge/internal/event"
)

const (
	dllHostExe = "dllhost.exe"
)

var (
	UnityBridgeImpl *wineUnityBridgeImpl = &wineUnityBridgeImpl{}

	localReadPipe  *os.File
	localWritePipe *os.File
	localEventPipe *os.File
)

func init() {
	// Check if wine is available.
	winePath, err := getWinePath()
	if err != nil {
		panic(err)
	}

	// Check if dllhost.exe is available and is a Windows executable.
	dllHostPath, err := getDLLHostPath()
	if err != nil {
		panic(err)
	}

	var remoteWritePipe *os.File
	localReadPipe, remoteWritePipe, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	var remoteReadPipe *os.File
	remoteReadPipe, localWritePipe, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	var remoteEventPipe *os.File
	localEventPipe, remoteEventPipe, err = os.Pipe()
	if err != nil {
		panic(err)
	}

	err = startDllHost(winePath, dllHostPath, remoteReadPipe, remoteWritePipe,
		remoteEventPipe)
	if err != nil {
		panic(err)
	}

	go loop()
}

type wineUnityBridgeImpl struct{}

func sendRequest(function byte, data *bytes.Buffer) ([]byte, error) {
	// Write function identifier
	_, err := localWritePipe.Write([]byte{function})
	if err != nil {
		return nil, err
	}

	if data != nil {
		// Write total data len.
		err = binary.Write(localWritePipe, binary.BigEndian, uint16(data.Len()))
		if err != nil {
			return nil, err
		}
		// Write data.
		_, err = localWritePipe.Write(data.Bytes())
		if err != nil {
			return nil, err
		}
	} else {
		err = binary.Write(localWritePipe, binary.BigEndian, uint16(0))
		if err != nil {
			return nil, err
		}
	}

	// Read response header.
	headerBuf := make([]byte, 1+2)
	_, err = io.ReadFull(localReadPipe, headerBuf)
	if err != nil {
		return nil, err
	}

	// Check function identifier.
	if headerBuf[0] != function {
		return nil, fmt.Errorf("unexpected function identifier: %d",
			headerBuf[0])
	}

	// Read response length.
	length := binary.BigEndian.Uint16(headerBuf[1:3])

	if length > 0 {
		// Read response data.
		response := make([]byte, length)
		_, err = io.ReadFull(localReadPipe, response)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, nil
}

func (ub wineUnityBridgeImpl) Create(name string, debuggable bool,
	logPath string) {
	var b bytes.Buffer

	if debuggable {
		b.WriteByte(1)
	} else {
		b.WriteByte(0)
	}

	binary.Write(&b, binary.BigEndian, uint16(len(name)))
	b.WriteString(name)

	binary.Write(&b, binary.BigEndian, uint16(len(logPath)))
	b.WriteString(logPath)

	_, err := sendRequest(0x00, &b)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) Destroy() {
	_, err := sendRequest(0x01, nil)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) Initialize() bool {
	res, err := sendRequest(0x02, nil)
	if err != nil {
		panic(err)
	}

	return res[0] != 0
}

func (ub wineUnityBridgeImpl) Uninitialize() {
	_, err := sendRequest(0x03, nil)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) SendEvent(e *event.Event, data uintptr,
	tag uint64) {
	var b bytes.Buffer

	if data != 0 {
		panic("SendEvent only supports data == 0 on Wine.")
	}

	binary.Write(&b, binary.BigEndian, e.Code())
	binary.Write(&b, binary.BigEndian, tag)

	_, err := sendRequest(0x04, &b)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) SendEventWithString(e *event.Event, data string,
	tag uint64) {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, e.Code())
	binary.Write(&b, binary.BigEndian, tag)
	binary.Write(&b, binary.BigEndian, uint16(len(data)))
	b.WriteString(data)

	_, err := sendRequest(0x05, &b)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) SendEventWithNumber(e *event.Event, data uint64,
	tag uint64) {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, e.Code())
	binary.Write(&b, binary.BigEndian, tag)
	binary.Write(&b, binary.BigEndian, data)

	_, err := sendRequest(0x06, &b)
	if err != nil {
		panic(err)
	}
}

func (ub wineUnityBridgeImpl) SetEventCallback(t event.Type,
	callback event.Callback) {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, t)
	binary.Write(&b, binary.BigEndian, callback != nil)

	_, err := sendRequest(0x07, &b)
	if err != nil {
		panic(err)
	}

	internal_event.SetEventCallback(t, callback)
}

func (ub wineUnityBridgeImpl) GetSecurityKeyByKeyChainIndex(index int) string {
	var b bytes.Buffer

	binary.Write(&b, binary.BigEndian, uint64(index))

	res, err := sendRequest(0x08, &b)
	if err != nil {
		panic(err)
	}

	return string(res)
}

func getWinePath() (string, error) {
	return exec.LookPath("wine")
}

func getDLLHostPath() (string, error) {
	dllHostPath, err := exec.LookPath(dllHostExe)
	if err != nil {
		// Try current directory.
		dllHostPath, err = exec.LookPath("./" + dllHostExe)
		if err != nil {
			return "", err
		}
	}

	peFile, err := pe.Open(dllHostPath)
	if err != nil {
		return "", fmt.Errorf("%q does not look like a Windows executable: %w",
			dllHostPath, err)
	}
	peFile.Close()

	return dllHostPath, nil
}

func startDllHost(winePath, dllHostPath string, remoteReadPipe,
	remoteWritePipe, remoteEventPipe *os.File) error {
	argv := []string{
		winePath,
		dllHostPath,
		"-read-fd",
		fmt.Sprintf("%d", getFd(remoteReadPipe)),
		"-write-fd",
		fmt.Sprintf("%d", getFd(remoteWritePipe)),
		"-event-fd",
		fmt.Sprintf("%d", getFd(remoteEventPipe)),
	}

	// Disable close on exec for the pipes.
	disableCloseOnExec(remoteReadPipe)
	disableCloseOnExec(remoteWritePipe)
	disableCloseOnExec(remoteEventPipe)

	_, err := syscall.ForkExec(winePath, argv,
		&syscall.ProcAttr{
			Files: []uintptr{
				getFd(os.Stdin),
				getFd(os.Stdout),
				getFd(os.Stderr),
			},
			Sys: &syscall.SysProcAttr{
				Foreground: false,
			},
			Env: []string{
				"WINEDEBUG=-all",
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error executing windows program: %w", err)
	}

	remoteReadPipe.Close()
	remoteWritePipe.Close()
	remoteEventPipe.Close()

	return nil
}

func disableCloseOnExec(file *os.File) {
	_, _, err := syscall.Syscall(syscall.SYS_FCNTL, getFd(file),
		syscall.F_SETFD, 0)
	if err != syscall.Errno(0) {
		panic(fmt.Sprintf("Error disabling close on exec: %s", err))
	}
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

func loop() {
	headerBuf := make([]byte, 18)
	for {
		if _, err := io.ReadFull(localEventPipe, headerBuf); err != nil {
			panic(fmt.Sprintf("Error reading data: %s", err))
		}

		eventCode := binary.BigEndian.Uint64(headerBuf[0:8])
		tag := binary.BigEndian.Uint64(headerBuf[8:16])
		length := binary.BigEndian.Uint16(headerBuf[16:18])

		data := make([]byte, length)
		if _, err := io.ReadFull(localEventPipe, data); err != nil {
			panic(fmt.Sprintf("Error reading data: %s", err))
		}

		internal_event.RunEventCallback(eventCode, data, tag)
	}
}
