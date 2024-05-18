package camera

import (
	"image"
	"image/color"
)

type RGB struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func NewRGB(r image.Rectangle) *RGB {
	w, h := r.Dx(), r.Dy()
	buf := make([]uint8, 3*w*h)
	return &RGB{Pix: buf, Stride: 3 * w, Rect: r}
}

func NewRGBFromBytes(data []byte, r image.Rectangle) *RGB {
	if len(data) != 3*r.Dx()*r.Dy() {
		panic("unexpected image data")
	}
	return &RGB{Pix: data, Stride: 3 * r.Dx(), Rect: r}
}

func (im *RGB) ColorModel() color.Model {
	return color.RGBAModel
}

func (im *RGB) Bounds() image.Rectangle {
	return im.Rect
}

func (im *RGB) At(x, y int) color.Color {
	if !(image.Point{x, y}.In(im.Rect)) {
		return color.RGBA{}
	}
	i := im.PixOffset(x, y)
	return color.RGBA{R: im.Pix[i+0], G: im.Pix[i+1], B: im.Pix[i+2], A: 255}
}

func (im *RGB) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(im.Rect)) {
		return
	}
	i := im.PixOffset(x, y)
	c1 := color.RGBAModel.Convert(c).(color.RGBA)
	im.Pix[i+0] = c1.R
	im.Pix[i+1] = c1.G
	im.Pix[i+2] = c1.B
}

func (im *RGB) PixOffset(x, y int) int {
	return (y-im.Rect.Min.Y)*im.Stride + (x-im.Rect.Min.X)*3
}
