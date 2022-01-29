package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper/winebridge"
)

var (
	writePipe io.Writer

	readBuffer  []byte
	writeBuffer []byte

	w wrapper.Wrapper
)

func main() {
	wineBridge, err := winebridge.New(2, os.Args[1:])
	if err != nil {
		panic(err)
	}

	readPipe := wineBridge.File(0)
	writePipe = wineBridge.File(1)

	// Initialize wrapper.
	w, err = wrapper.New("./unitybridge.dll")
	if err != nil {
		panic(err)
	}

	lengthBuffer := make([]byte, 4)
	for {
		_, err := io.ReadFull(readPipe, lengthBuffer)
		if err == io.ErrUnexpectedEOF {
			panic(err)
		}
		if err == io.EOF {
			break
		}

		length := binary.LittleEndian.Uint32(lengthBuffer)

		err = processRead(readPipe, int(length-4))
		if err != nil {
			panic(err)
		}
	}
}

func processRead(readPipe io.Reader, length int) error {
	if length > len(readBuffer) {
		readBuffer = make([]byte, length)
	}

	sizedRequestBuffer := readBuffer[:length]
	_, err := io.ReadFull(readPipe, sizedRequestBuffer)
	if err != nil {
		return err
	}

	function := sizedRequestBuffer[0]
	switch function {
	case wrapper.FuncCreateUnityBridge:
		runCreateUnityBridge(sizedRequestBuffer[1:])
	case wrapper.FuncDestroyUnityBridge:
		runDestroyUnityBridge()
	case wrapper.FuncUnityBridgeInitialize:
		runUnityBridgeInitialize()
	case wrapper.FuncUnityBridgeUninitialize:
		runUnityBridgeUninitialize()
	case wrapper.FuncUnitySetEventCallback:
		runUnitySetEventCallback(sizedRequestBuffer[1:])
	case wrapper.FuncUnitySendEvent:
		runUnitySendEvent(sizedRequestBuffer[1:])
	case wrapper.FuncUnitySendEventWithNumber:
		runUnitySendEventWithNumber(sizedRequestBuffer[1:])
	case wrapper.FuncUnitySendEventWithString:
		runUnitySendEventWithString(sizedRequestBuffer[1:])
	}

	return nil
}

func runCreateUnityBridge(buffer []byte) {
	debuggable := false
	if buffer[0] != 0 {
		debuggable = true
	}

	name := string(buffer[1:])

	w.CreateUnityBridge(name, debuggable)
}

func runDestroyUnityBridge() {
	w.DestroyUnityBridge()
}

func runUnityBridgeInitialize() {
	// Currently there is no way to return the error to the Linux side.
	// Simply log it here for now.
	if !w.UnityBridgeInitialize() {
		fmt.Println("Unity Bridge failed to initialize.")
	}
}

func runUnityBridgeUninitialize() {
	w.UnityBridgeUninitialize()
}

func runUnitySetEventCallback(buffer []byte) {
	add := false
	if buffer[0] == 1 {
		add = true
	}
	eventCode := binary.LittleEndian.Uint64(buffer[1:])
	if add {
		w.UnitySetEventCallback(eventCode,
			eventCallback)
	} else {
		w.UnitySetEventCallback(eventCode, nil)
	}
}

func runUnitySendEvent(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	info := buffer[16:]

	w.UnitySendEvent(eventCode, info, tag)
}

func runUnitySendEventWithNumber(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	data := binary.LittleEndian.Uint64(buffer[16:])

	w.UnitySendEventWithNumber(eventCode, data, tag)
}

func runUnitySendEventWithString(buffer []byte) {
	eventCode := binary.LittleEndian.Uint64(buffer)
	tag := binary.LittleEndian.Uint64(buffer[8:])
	data := string(buffer[16:])

	w.UnitySendEventWithString(eventCode, data, tag)
}

func eventCallback(eventCode uint64, data []byte, tag uint64) {
	length := 4 + 8 + 8 + len(data)
	if len(writeBuffer) < length {
		writeBuffer = make([]byte, length)
	}

	sizedWriteBuffer := writeBuffer[:length]
	binary.LittleEndian.PutUint32(sizedWriteBuffer, uint32(length))
	binary.LittleEndian.PutUint64(sizedWriteBuffer[4:], eventCode)
	binary.LittleEndian.PutUint64(sizedWriteBuffer[12:], tag)
	copy(sizedWriteBuffer[20:], data)

	_, err := writePipe.Write(sizedWriteBuffer)
	if err != nil {
		panic(err)
	}
}
