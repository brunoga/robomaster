package qrcode

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"image"
	"net"
	"strings"

	"github.com/brunoga/robomaster/unitybridge/support"
	"github.com/skip2/go-qrcode"
)

// QRCode handles generating and parsing the data for the Robomaster connection
// qrcode that is used to associate the robot with an app ID and also to tell
// it about the network and password to use.
type QRCode struct {
	appID       uint64
	countryCode string
	ssID        string
	password    string
	bssID       *net.HardwareAddr
}

// New creates a new QRCode instance with the given parameters.
func New(appID uint64,
	countryCode, ssID, password, bssID string) (*QRCode, error) {
	trimmedCountryCode := strings.TrimSpace(countryCode)
	if len(trimmedCountryCode) != 2 {
		return nil, fmt.Errorf("country code must be 2 characters")
	}

	trimmedSsID := strings.TrimSpace(ssID)
	if len(trimmedSsID) == 0 {
		return nil, fmt.Errorf("SSID must be non-empty")
	}

	trimmedPassword := strings.TrimSpace(password)
	if len(trimmedPassword) == 0 {
		return nil, fmt.Errorf("password must be non-empty")
	}

	var resultBssID *net.HardwareAddr
	if len(strings.TrimSpace(bssID)) != 0 {
		parsedBssID, err := net.ParseMAC(bssID)
		if err != nil {
			return nil, err
		}
		if len(parsedBssID) != 6 {
			return nil, fmt.Errorf(
				"BSSID must have exactly 6 octets")
		}

		resultBssID = &parsedBssID
	}

	return &QRCode{
		appID,
		trimmedCountryCode,
		trimmedSsID,
		trimmedPassword,
		resultBssID,
	}, nil
}

// NewFromMessage parses the given message and returns a QRCode instance based
// on it. This message is what you got if you use a normal QRCode reader to
// read the Robomaster app generated QRCode.
func NewFromMessage(message string) (*QRCode, error) {
	q := &QRCode{}
	err := q.decodeMessage(message)
	if err != nil {
		return nil, err
	}

	return q, nil
}

// AppID returns the app ID for this QRCode.
func (q *QRCode) AppID() uint64 {
	return q.appID
}

// CountryCode returns the country code for this QRCode.
func (q *QRCode) CountryCode() string {
	return q.countryCode
}

// SSID returns the SSID for this QRCode.
func (q *QRCode) SsID() string {
	return q.ssID
}

// Password returns the password for this QRCode.
func (q *QRCode) Password() string {
	return q.password
}

// BssID returns the BSSID for this QRCode.
func (q *QRCode) BssID() *net.HardwareAddr {
	return q.bssID
}

// String returns a string representation of this QRCode.
func (q *QRCode) String() string {
	return fmt.Sprintf("App Id : %d, Country Code : %q, SSID : %q, "+
		"Password : %q, BSSID : %q", q.appID, q.countryCode, q.ssID,
		q.password, q.bssID)
}

// Message returns the message for this QRCode. This is the result of encoding
// the QRCode instance into a string that can be used to generate an actual
// qr-code image that can be used by a Robomaster robot.
func (q *QRCode) Message() string {
	return q.encodeMessage()
}

// Image returns an size X size image.Image that represents the QRCode instance.
// This image is readable by a Robomaster robot.
func (q *QRCode) Image(size int) (image.Image, error) {
	qrc, err := qrcode.New(q.encodeMessage(), qrcode.Medium)
	if err != nil {
		return nil, err
	}

	return qrc.Image(size), nil
}

// Text returns a string representation of the QRCode instance that can be
// printed to the console. This should be readable by a Robomaster robot.
func (q *QRCode) Text() (string, error) {
	qrc, err := qrcode.New(q.encodeMessage(), qrcode.Medium)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, line := range qrc.Bitmap() {
		for _, black := range line {
			if black {
				sb.WriteString("  ") // White space
			} else {
				sb.WriteString("██") // Black square
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (q *QRCode) encodeMessage() string {
	var b bytes.Buffer

	bytesSsid := []byte(q.ssID)
	bytesPassword := []byte(q.password)

	hasBssId := uint16(0)
	if q.bssID != nil {
		hasBssId = uint16(1)
	}

	metadata := (hasBssId << 11) | (uint16(len(bytesPassword)) << 6) |
		uint16(len(bytesSsid))

	data1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(data1, metadata)

	b.Write(data1)

	data2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data2, q.appID)

	b.Write(data2)

	b.Write([]byte(q.countryCode))

	b.Write([]byte(q.ssID))
	b.Write([]byte(q.password))

	if q.bssID != nil {
		bssIdString := strings.ReplaceAll(q.bssID.String(), ":", "")
		b.Write([]byte(bssIdString))
	}

	data := b.Bytes()
	support.SimpleEncryptDecrypt(data)

	return base64.StdEncoding.EncodeToString(data)
}

func (q *QRCode) decodeMessage(message string) error {
	data, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return err
	}

	support.SimpleEncryptDecrypt(data)

	metadata := binary.LittleEndian.Uint16(data)

	hasBssId := metadata >> 11
	lenPassword := (metadata & 0b0000001111111111) >> 6
	lenSsId := (metadata & 0b0000000001111111)

	q.appID = binary.LittleEndian.Uint64(data[2:])

	q.countryCode = string(data[10:12])

	q.ssID = string(data[12 : 12+lenSsId])

	q.password = string(data[12+lenSsId : 12+lenSsId+lenPassword])

	if hasBssId != 0 {
		parsedBssId, err := net.ParseMAC(
			string(data[12+lenSsId+lenPassword : 12+lenSsId+lenPassword+12]))
		if err != nil {
			return err
		}

		q.bssID = &parsedBssId
	} else {
		q.bssID = nil
	}

	return nil
}
