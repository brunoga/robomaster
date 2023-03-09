package command

import (
	"encoding/binary"
	"fmt"
)

const (
	getProductVersionResponseSize = 13
)

func init() {
	Register(getProductVersionSet, getProductVersionID, NewGetProductVersionResponse())
}

// GetProductVersionResponse is a response to a GetVersionRequest.
type GetProductVersionResponse struct {
	*baseResponse
}

var _ Response = (*GetProductVersionResponse)(nil)

// NewGetProductVersionResponse creates a new GetVersionResponse.
func NewGetProductVersionResponse() *GetProductVersionResponse {
	return &GetProductVersionResponse{
		baseResponse: newBaseResponse(
			getProductVersionSet,
			getProductVersionID,
			getProductVersionType,
			getProductVersionResponseSize,
		),
	}
}

// New implements the Command interface.
func (g *GetProductVersionResponse) New(data []byte) Command {
	r := NewGetProductVersionResponse()
	r.data = data

	return r
}

// Version returns the version of the SDK that exists in the robot as a string.
func (g *GetProductVersionResponse) Version() string {
	return fmt.Sprintf("%02d.%02d.%04d", g.data[12], g.data[11],
		binary.LittleEndian.Uint16(g.data[9:11]))
}
