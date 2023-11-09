package event

import (
	"encoding/binary"
	"fmt"
)

// DataType is the type of data returned in an event callback.
type DataType int

const (
	StringDataType DataType = iota
	NumberDataType
)

// String returns the string representation of the DataType.
func (d DataType) String() string {
	switch d {
	case StringDataType:
		return "String"
	case NumberDataType:
		return "Number"
	default:
		return "Unknown"
	}
}

// ParseData parses the data according to the DataType and returns it as an any
// type.
func (d DataType) ParseData(data []byte) any {
	switch d {
	case StringDataType:
		return string(data)
	case NumberDataType:
		if len(data) != 8 {
			panic(fmt.Sprintf("Invalid number data length. Expected 8, got %d.",
				len(data)))
		}
		return binary.LittleEndian.Uint64(data)
	default:
		return nil
	}
}

// DataTypeFromTag returns the DataType and the actual tag value (i.e. without
// the data type) from a tag.
func DataTypeFromTag(tag uint64) (DataType, uint64) {
	return DataType((tag >> 56) & 0xff), tag & 0x00ffffffffffffff
}
