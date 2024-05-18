package mobile

import (
	"image/draw"

	"github.com/brunoga/robomaster/module/camera"
	"github.com/brunoga/robomaster/unitybridge/support/qrcode"
)

// QRCode allows generating QRCodes that can be read by a Robomaster robot.
type QRCode struct {
	q *qrcode.QRCode
}

// NewQRCode creates a new QRCode instance. The bssID parameter is optional
// and may be an empty string.
func NewQRCode(appID int64, countryCode, ssID, password,
	bssID string) (*QRCode, error) {
	q, err := qrcode.New(uint64(appID), countryCode, ssID, password, bssID)
	if err != nil {
		return nil, err
	}

	return &QRCode{
		q: q,
	}, nil
}

// Image returns the QRCode image (with size as its width and height) as a byte
// slice with RGB24 image data.
func (q *QRCode) Image(size int) ([]byte, error) {
	img, err := q.q.Image(size)
	if err != nil {
		return nil, err
	}

	rgb := camera.NewRGB(img.Bounds())

	draw.Draw(rgb, img.Bounds(), img, img.Bounds().Min, draw.Src)

	return rgb.Pix, nil
}
