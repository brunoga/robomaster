package datatype

import "encoding/binary"

// DataType is the type of data returned in an event callback.
type DataType int

const (
	String DataType = iota
	Number
)

// String returns the string representation of the DataType.
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

// ParseData parses the data according to the DataType and returns it as an any
// type.
func (d DataType) ParseData(data []byte) any {
	switch d {
	case String:
		return string(data)
	case Number:
		return binary.LittleEndian.Uint64(data)
	default:
		return nil
	}
}

// FromTag returns the DataType and the actual tag value (i.e. without the data
// type) from a tag.
func FromTag(tag uint64) (DataType, uint64) {
	return DataType((tag >> 56) & 0xff), tag & 0x00ffffffffffffff
}
