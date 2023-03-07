package command

import (
	"fmt"
)

const (
	getVersionResponseSize = 6
)

func init() {
	Register(getVersionSet, getVersionID, NewGetVersionResponse())
}

// GetVersionResponse is a response to a GetVersionRequest.
type GetVersionResponse struct {
	*baseResponse
}

var _ Response = (*GetVersionResponse)(nil)

// NewGetVersionResponse creates a new GetVersionResponse.
func NewGetVersionResponse() *GetVersionResponse {
	return &GetVersionResponse{
		baseResponse: newBaseResponse(
			getVersionSet,
			getVersionID,
			getVersionType,
			getVersionResponseSize,
		),
	}
}

// New implements the Command interface.
func (g *GetVersionResponse) New() Command {
	return NewGetVersionResponse()
}

// Version returns the version of the SDK that exists in the robot as a string.
func (g *GetVersionResponse) Version() string {
	return fmt.Sprintf("%02d.%02d.%02d.%02d", g.data[0], g.data[1], g.data[2],
		g.data[3])
}
