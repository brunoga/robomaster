package wrapper

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge/wrapper/winebridge"
)

type Linux struct {
	readPipe  io.Reader
	writePipe io.Writer

	wineBridge *winebridge.WineBridge

	m                sync.RWMutex
	eventCallbackMap map[unity.EventType]EventCallback
}

func New(string) (Wrapper, error) {
	localReadPipe, remoteWritePipe, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	remoteReadPipe, localWritePipe, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	wineBridge, err := winebridge.New("./winewrapper.exe",
		remoteReadPipe, remoteWritePipe)
	if err != nil {
		return nil, err
	}

	err = wineBridge.Start()
	if err != nil {
		panic(err)
	}

	l := &Linux{
		readPipe:         localReadPipe,
		writePipe:        localWritePipe,
		wineBridge:       wineBridge,
		eventCallbackMap: make(map[unity.EventType]EventCallback),
	}

	go l.readLoop()

	return l, nil
}

func (l *Linux) CreateUnityBridge(name string, debuggable bool) {
	// size (4 bytes) + function (1 byte) + debuggable (1 byte) + len(name)
	buffer := make([]byte, 4+1+1+len(name))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncCreateUnityBridge
	if debuggable {
		buffer[5] = 1
	} else {
		buffer[5] = 0
	}
	copy(buffer[6:], name)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) DestroyUnityBridge() {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncDestroyUnityBridge

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnityBridgeInitialize() bool {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnityBridgeInitialize

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}

	return true
}

func (l *Linux) UnityBridgeUninitialize() {
	// size (4 bytes) + function (1 byte)
	buffer := make([]byte, 4+1)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnityBridgeUninitialize

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEvent(eventCode uint64, info []byte, tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + len(info)
	buffer := make([]byte, 4+1+8+8+len(info))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEvent
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	copy(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEventWithNumber(eventCode, info, tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + info (8 bytes)
	buffer := make([]byte, 4+1+8+8+8)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEventWithNumber
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	binary.LittleEndian.PutUint64(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySendEventWithString(eventCode uint64, info string,
	tag uint64) {
	// size (4 bytes) + function (1 byte) + eventCode (8 bytes) +
	// tag (8 bytes) + len(info)
	buffer := make([]byte, 4+1+8+8+len(info))
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySendEventWithString
	binary.LittleEndian.PutUint64(buffer[5:], eventCode)
	binary.LittleEndian.PutUint64(buffer[13:], tag)
	copy(buffer[21:], info)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) UnitySetEventCallback(eventCode uint64,
	eventCallback EventCallback) {
	l.m.Lock()
	defer l.m.Unlock()

	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode << 32, eventCode & 0xffffffff)
		return
	}

	_, ok := l.eventCallbackMap[event.Type()]

	add := false
	if eventCallback == nil {
		if ok {
			delete(l.eventCallbackMap, event.Type())
		}
	} else {
		if !ok {
			l.eventCallbackMap[event.Type()] = eventCallback
		}

		add = true
	}

	// size (4 bytes) + function (1 byte) + add (1 byte) +
	// eventCode (8 bytes)
	buffer := make([]byte, 4+1+1+8)
	binary.LittleEndian.PutUint32(buffer, uint32(len(buffer)))
	buffer[4] = FuncUnitySetEventCallback
	if add {
		buffer[5] = 1
	} else {
		buffer[5] = 0
	}
	binary.LittleEndian.PutUint64(buffer[6:], eventCode)

	_, err := l.writePipe.Write(buffer)
	if err != nil {
		panic(err)
	}
}

func (l *Linux) readLoop() {
	readBuffer := make([]byte, 10000)
	lengthBuffer := readBuffer[:4]
	for {
		_, err := io.ReadFull(l.readPipe, lengthBuffer)
		if err != nil {
			panic(err)
		}

		length := binary.LittleEndian.Uint32(lengthBuffer)

		length -= 4

		if len(readBuffer) < int(length) {
			readBuffer = make([]byte, length)
		}

		sizedReadBuffer := readBuffer[:length]

		_, err = io.ReadFull(l.readPipe, sizedReadBuffer)
		if err != nil {
			panic(err)
		}

		eventCode := binary.LittleEndian.Uint64(sizedReadBuffer)
		tag := binary.LittleEndian.Uint64(sizedReadBuffer[8:])
		data := sizedReadBuffer[16:]

		l.maybeRunCallback(eventCode, data, tag)
	}
}

func (l *Linux) maybeRunCallback(eventCode uint64, data []byte, tag uint64) {
	l.m.RLock()
	defer l.m.RUnlock()

	event := unity.NewEventFromCode(eventCode)
	if event == nil {
		log.Printf("Unknown event with code %d (Type:%d, SubType:%d).\n",
			eventCode, eventCode << 32, eventCode & 0xffffffff)
		return
	}

	callback, ok := l.eventCallbackMap[event.Type()]
	if !ok {
		log.Printf("No callback for event %s (tag=%d).\n",
			unity.EventTypeName(event.Type()), tag)
		return
	}

	callback(eventCode, data, tag)
}

func getFd(file *os.File) uintptr {
	rawConn, err := file.SyscallConn()
	if err != nil {
		panic(err)
	}

	var fileFd uintptr
	err = rawConn.Control(func(fd uintptr) {
		fileFd = fd
	})
	if err != nil {
		panic(err)
	}

	return fileFd
}
