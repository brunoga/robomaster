package datatype

// DataType is the type of data returned in an event callback.
type DataType int

const (
	String DataType = iota
	Number
)

func (d DataType) String() string {
	switch d {
	case String:
		return "String"
	case Number:
		return "Number"
	default:
		return "Unknown"
	}
}

// FromTag returns the DataType and the actual tag value (i.e. without the data
// type) from a tag.
func FromTag(tag uint64) (DataType, uint64) {
	return DataType((tag >> 56) & 0xff), tag & 0x00ffffffffffffff
}
