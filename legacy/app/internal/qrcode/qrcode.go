package qrcode

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"net"
	"strings"

	"git.bug-br.org.br/bga/robomasters1/app/internal/support"
)

type QRCode struct {
	appId       uint64
	countryCode string
	ssId        string
	password    string
	bssId       *net.HardwareAddr
}

func NewQRCode(appId uint64,
	countryCode, ssId, password, bssId string) (*QRCode, error) {
	trimmedCountryCode := strings.TrimSpace(countryCode)
	if len(trimmedCountryCode) != 2 {
		return nil, fmt.Errorf("country code must be 2 characters")
	}

	trimmedSsId := strings.TrimSpace(ssId)
	if len(trimmedSsId) == 0 {
		return nil, fmt.Errorf("SSID must be non-empty")
	}

	trimmedPassword := strings.TrimSpace(password)
	if len(trimmedPassword) == 0 {
		return nil, fmt.Errorf("password must be non-empty")
	}

	var resultBssId *net.HardwareAddr = nil
	if len(strings.TrimSpace(bssId)) != 0 {
		parsedBssId, err := net.ParseMAC(bssId)
		if err != nil {
			return nil, err
		}
		if len(parsedBssId) != 6 {
			return nil, fmt.Errorf(
				"BSSID must have exactly 6 octets")
		}

		resultBssId = &parsedBssId
	}

	return &QRCode{
		appId,
		trimmedCountryCode,
		trimmedSsId,
		trimmedPassword,
		resultBssId,
	}, nil
}

func ParseQRCodeMessage(message string) (*QRCode, error) {
	q := &QRCode{}
	err := q.decodeQRCodeMessage(message)
	if err != nil {
		return nil, err
	}

	return q, nil
}

func (q *QRCode) AppId() uint64 {
	return q.appId
}

func (q *QRCode) CountryCode() string {
	return q.countryCode
}

func (q *QRCode) SsId() string {
	return q.ssId
}

func (q *QRCode) Password() string {
	return q.password
}

func (q *QRCode) BssId() *net.HardwareAddr {
	return q.bssId
}

func (q *QRCode) EncodedMessage() string {
	return q.encodeQRCodeMessage()
}

func (q *QRCode) String() string {
	return fmt.Sprintf("App Id : %d, Country Code : %q, SSID : %q, "+
		"Password : %q, BSSID : %q", q.appId, q.countryCode, q.ssId,
		q.password, q.bssId)
}

func (q *QRCode) encodeQRCodeMessage() string {
	var b bytes.Buffer

	bytesSsid := []byte(q.ssId)
	bytesPassword := []byte(q.password)

	hasBssId := uint16(0)
	if q.bssId != nil {
		hasBssId = uint16(1)
	}

	metadata := (hasBssId << 11) | (uint16(len(bytesPassword)) << 6) |
		uint16(len(bytesSsid))

	data1 := make([]byte, 2)
	binary.LittleEndian.PutUint16(data1, metadata)

	b.Write(data1)

	data2 := make([]byte, 8)
	binary.LittleEndian.PutUint64(data2, q.appId)

	b.Write(data2)

	b.Write([]byte(q.countryCode))

	b.Write([]byte(q.ssId))
	b.Write([]byte(q.password))

	if q.bssId != nil {
		bssIdString := strings.ReplaceAll(q.bssId.String(), ":", "")
		b.Write([]byte(bssIdString))
	}

	data := b.Bytes()
	support.InPlaceEncodeDecode(data)

	return base64.StdEncoding.EncodeToString(data)
}

func (q *QRCode) decodeQRCodeMessage(message string) error {
	data, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		return err
	}

	support.InPlaceEncodeDecode(data)

	metadata := binary.LittleEndian.Uint16(data)

	hasBssId := metadata >> 11
	lenPassword := (metadata & 0b0000001111111111) >> 6
	lenSsId := (metadata & 0b0000000001111111)

	q.appId = binary.LittleEndian.Uint64(data[2:])

	q.countryCode = string(data[10:12])

	q.ssId = string(data[12 : 12+lenSsId])

	q.password = string(data[12+lenSsId : 12+lenSsId+lenPassword])

	if hasBssId != 0 {
		parsedBssId, err := net.ParseMAC(
			string(data[12+lenSsId+lenPassword : 12+lenSsId+lenPassword+12]))
		if err != nil {
			return err
		}

		q.bssId = &parsedBssId
	} else {
		q.bssId = nil
	}

	return nil
}
