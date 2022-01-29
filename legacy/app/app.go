package app

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"git.bug-br.org.br/bga/robomasters1/app/internal"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity"
	"git.bug-br.org.br/bga/robomasters1/app/internal/dji/unity/bridge"
	"git.bug-br.org.br/bga/robomasters1/app/internal/pairing"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/skratchdot/open-golang/open"

	internalqrcode "git.bug-br.org.br/bga/robomasters1/app/internal/qrcode"
)

type App struct {
	id  uint64
	qrc *internalqrcode.QRCode
	pl  *pairing.Listener
	cc  *internal.CommandController
}

func New(countryCode, ssId, password, bssId string) (*App, error) {
	appId, err := generateAppId()
	if err != nil {
		return nil, err
	}

	qrc, err := internalqrcode.NewQRCode(appId, countryCode, ssId,
		password, bssId)
	if err != nil {
		return nil, err
	}

	a, err := NewWithAppID(countryCode, ssId, password, bssId, appId)
	if err != nil {
		return nil, err
	}

	a.qrc = qrc

	return a, nil
}

func NewWithAppID(countryCode, ssId, password, bssId string,
	appId uint64) (*App, error) {
	log.Printf("Using app ID %d.\n", appId)

	cc, err := internal.NewCommandController()
	if err != nil {
		return nil, err
	}

	return &App{
		appId,
		nil,
		pairing.NewListener(appId),
		cc,
	}, nil
}

func (a *App) Start(textMode bool) error {
	if a.qrc != nil {
		var err error
		if textMode {
			err = a.showTextQRCode()
		} else {
			err = a.showPNGQRCode()
		}
		if err != nil {
			return fmt.Errorf("error showing QR code: %w", err)
		}
	}

	// Setup Unity Bridge.
	if !bridge.IsSetup() {
		bridge.Setup("Robomaster", true)
	}

	ub := bridge.Instance()

	connectingIP := net.IP{}

	a.cc.StartListening(dji.KeyAirLinkConnection,
		func(result *dji.Result, wg *sync.WaitGroup) {
			if result.Value().(bool) {
				a.pl.SendACK(connectingIP)
			}

			wg.Done()
		})

	// Reset connection to defaults.
	err := ub.SendEvent(unity.NewEventWithSubType(
		unity.EventTypeConnection, 2), "192.168.2.1")
	if err != nil {
		panic(err)
	}
	err = ub.SendEvent(unity.NewEventWithSubType(
		unity.EventTypeConnection, 3), uint64(10607))
	if err != nil {
		panic(err)
	}
	err = ub.SendEvent(unity.NewEvent(unity.EventTypeConnection))
	if err != nil {
		panic(err)
	}

	eventChan, err := a.pl.Start()
	if err != nil {
		return fmt.Errorf("error starting pairing listener: %w", err)
	}

L:
	for {
		select {
		case pairingEvent, ok := <-eventChan:
			if !ok {
				break L
			}

			if pairingEvent.Type() == pairing.EventAdd {
				connectingIP = pairingEvent.IP()
				err = ub.SendEvent(unity.NewEventWithSubType(
					unity.EventTypeConnection, 1))
				if err != nil {
					panic(err)
				}
				err = ub.SendEvent(unity.NewEventWithSubType(
					unity.EventTypeConnection, 2),
					pairingEvent.IP().String())
				if err != nil {
					panic(err)
				}
				err = ub.SendEvent(unity.NewEventWithSubType(
					unity.EventTypeConnection, 3),
					uint64(10607))
				if err != nil {
					panic(err)
				}
				err = ub.SendEvent(unity.NewEvent(
					unity.EventTypeConnection))
				if err != nil {
					panic(err)
				}

				break L
			}
		}
	}

	return nil
}

func (a *App) CommandController() *internal.CommandController {
	return a.cc
}

func (a *App) showTextQRCode() error {
	qrc, err := qrcode.New(a.qrc.EncodedMessage(), qrcode.Medium)
	if err != nil {
		return err
	}

	fmt.Println(qrc.ToString(false))

	return nil
}

func (a *App) showPNGQRCode() error {
	pngData, err := qrcode.Encode(a.qrc.EncodedMessage(), qrcode.Medium,
		256)
	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return err
	}

	fileName := f.Name()

	_, err = f.Write(pngData)
	if err != nil {
		f.Close()
		return err
	}

	f.Close()

	return open.Run(fileName)
}

func generateAppId() (uint64, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return 0, err
	}

	// Create an app ID out of the first 8 bytes of the UUID.
	return binary.LittleEndian.Uint64(id[0:9]), nil
}
