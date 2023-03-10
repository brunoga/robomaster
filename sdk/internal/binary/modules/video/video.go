package video

import (
	"fmt"
	"image"
	"net"
	"sync"

	"github.com/brunoga/robomaster/sdk/internal/binary/modules/control"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/command"
	"github.com/brunoga/robomaster/sdk/internal/binary/protocol/message"
	"github.com/brunoga/robomaster/sdk/internal/h264"
	"github.com/brunoga/robomaster/sdk/modules/video"
	"github.com/brunoga/robomaster/sdk/support/logger"
	"github.com/brunoga/robomaster/sdk/types"
)

const (
	streamPort = 40921
)

type Video struct {
	l       *logger.Logger
	control *control.Control

	decoder *h264.Decoder

	resolution video.Resolution

	m        sync.Mutex
	quitChan chan struct{}
	waitChan chan struct{}
	handler  video.Handler
}

var _ video.Video = (*Video)(nil)

func New(control *control.Control, l *logger.Logger) (*Video, error) {
	v := &Video{
		l:          l,
		control:    control,
		resolution: video.ResolutionInvalid,
	}

	decoder, err := h264.NewDecoder(v.frameCallback)
	if err != nil {
		return nil, fmt.Errorf("error creating h264 decoder: %w", err)
	}

	v.decoder = decoder

	return v, nil
}

func (v *Video) Start(resolution video.Resolution, handler video.Handler) error {
	if handler == nil {
		return fmt.Errorf("video handler must not be nil")
	}

	v.m.Lock()
	defer v.m.Unlock()

	if v.handler != nil {
		return fmt.Errorf("video already started")
	}

	v.handler = handler

	v.quitChan = make(chan struct{})
	v.waitChan = make(chan struct{})

	go v.videoLoop(resolution)

	return nil
}

func (v *Video) Stop() error {
	v.m.Lock()
	defer v.m.Unlock()

	if v.handler == nil {
		return fmt.Errorf("video not started")
	}

	close(v.quitChan)

	<-v.waitChan

	v.handler = nil
	v.quitChan = nil
	v.waitChan = nil

	return nil
}

func (v *Video) frameCallback(data []byte, width, height int) {
	frameRGBA := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: width, Y: height},
	})

	frameRGBA.Pix = data

	v.handler(frameRGBA)
}

func (v *Video) videoLoop(resolution video.Resolution) {
	defer close(v.waitChan)

	// Enter SDK streaming mode.
	err := v.sendStreamControlRequest(1, 1, resolution)
	if err != nil {
		v.l.ERROR("Error entering SDK streaming mode: %s", err)
		v.Stop()
		return
	}

	// Start streaming
	err = v.sendStreamControlRequest(2, 1, resolution)
	if err != nil {
		v.l.ERROR("Error starting streaming: %s", err)
		v.Stop()
		return
	}

	streamAddr := fmt.Sprintf("%s:%d", v.control.IP(), streamPort)

	var streamConn net.Conn

	if v.control.ConnectionProtocol() == types.ConnectionProtocolTCP {
		streamConn, err = net.Dial("tcp", streamAddr)
	} else {
		streamConn, err = net.Dial("udp", streamAddr)
	}

	if err != nil {
		return
	}
	defer streamConn.Close()

	err = v.decoder.Open()
	if err != nil {
		return
	}
	defer v.decoder.Close()

	readBuffer := make([]byte, 16384)

L:
	for {
		select {
		case <-v.quitChan:
			break L
		default:
			n, err := streamConn.Read(readBuffer)
			if err != nil {
				break L
			}

			v.decoder.SendData(readBuffer[:n])
		}
	}

	// Stop streaming.
	err = v.sendStreamControlRequest(2, 0, resolution)
	if err != nil {
		v.l.ERROR("Error stopping streaming: %s", err)
		v.Stop()
		return
	}

	// Exit SDK streaming mode.
	err = v.sendStreamControlRequest(1, 0, resolution)
	if err != nil {
		v.l.ERROR("Error exiting SDK streaming mode: %s", err)
		v.Stop()
		return
	}
}

func (v *Video) sendStreamControlRequest(control byte, state byte,
	resolution video.Resolution) error {
	cmd := command.NewStreamCtrlRequest()
	cmd.SetControl(control)
	if v.control.ConnectionMode() == types.ConnectionModeUSB {
		cmd.SetConnectionType(1)
	}
	cmd.SetState(state)
	cmd.SetResolution(byte(resolution))

	resp, err := v.control.SendSync(message.New(v.control.HostByte(), protocol.HostToByte(1, 0), cmd))
	if err != nil {
		return fmt.Errorf("error sending stream control request: %w", err)
	}

	if !resp.Command().(command.Response).Ok() {
		v.l.ERROR("error on stream control request: not ok")
	}

	return nil
}
