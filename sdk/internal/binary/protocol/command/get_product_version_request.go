package command

const (
	getProductVersionRequestSize = 9
)

func init() {
	Register(getProductVersionSet, getProductVersionID, NewGetProductVersionRequest())
}

// GetProductVersionRequest is a command to get the product version.
type GetProductVersionRequest struct {
	*baseRequest
}

var _ Request = (*GetProductVersionRequest)(nil)

// NewGetProductVersionRequest creates a new GetVersionRequest.
func NewGetProductVersionRequest() *GetProductVersionRequest {
	r := &GetProductVersionRequest{
		baseRequest: newBaseRequest(
			getProductVersionSet,
			getProductVersionID,
			getProductVersionType,
			getProductVersionRequestSize,
		),
	}

	r.data[0] = 4 // Default "file type", whatever that means.
	r.data[5] = 0xff
	r.data[6] = 0xff
	r.data[7] = 0xff
	r.data[8] = 0xff

	return r
}

// New implements the Command interface.
func (g *GetProductVersionRequest) New(data []byte) Command {
	r := NewGetProductVersionRequest()
	r.data = data

	return r
}

func (g *GetProductVersionRequest) SetFileType(fileType byte) {
	g.data[0] = fileType
}

func (g *GetProductVersionRequest) FileType() byte {
	return g.data[0]
}
