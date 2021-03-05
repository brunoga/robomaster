package protocol

import "fmt"

const (
	getVersionRespSize = 30
)

type GetVersion struct {
	*Data

	aa int
	bb int
	cc int
	dd int
}

func NewGetVersion() *GetVersion {
	return &GetVersion{
		NewData("GetVersion", 0x00, 0x01),
		0,
		1,
		0,
		0,
	}
}

func (g *GetVersion) AA() int {
	return g.aa
}

func (g *GetVersion) BB() int {
	return g.bb
}

func (g *GetVersion) CC() int {
	return g.cc
}

func (g *GetVersion) DD() int {
	return g.dd
}

func (g *GetVersion) UnpackResp(buf []byte, offset int) error {
	if len(buf) < getVersionRespSize {
		return fmt.Errorf("buffer is too small")
	}

	if err := g.Data.UnpackResp(buf, offset); err != nil {
		return err
	}

	g.aa = int(buf[0])
	g.bb = int(buf[1])
	g.cc = int(buf[2])
	g.dd = int(buf[3])

	return nil
}
