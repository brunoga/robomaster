package rgb

import (
	"image"
	"image/color"
)

type Image struct {
	Pix    []uint8
	Stride int
	Rect   image.Rectangle
}

func NewImage(r image.Rectangle) *Image {
	w, h := r.Dx(), r.Dy()
	return &Image{
		Pix:    make([]uint8, 3*w*h),
		Stride: 3 * w,
		Rect:   r,
	}
}

func (p *Image) ColorModel() color.Model {
	return ColorModel
}

func (p *Image) Bounds() image.Rectangle {
	return p.Rect
}

func (p *Image) At(x, y int) color.Color {
	return p.RGBAAt(x, y)
}

func (p *Image) RGBAAt(x, y int) color.RGBA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return color.RGBA{}
	}
	i := (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*3
	return color.RGBA{p.Pix[i+0], p.Pix[i+1], p.Pix[i+2], 0xFF}
}

var ColorModel = color.ModelFunc(rgbModel)

func rgbModel(c color.Color) color.Color {
	if _, ok := c.(RGB); ok {
		return c
	}
	r, g, b, _ := c.RGBA()
	return RGB{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8)}
}

type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(0xFFFF)
	return
}
