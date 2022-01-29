package internal

import (
	"encoding/xml"
)

type Dji struct {
	XMLName xml.Name `xml:"dji"`

	Attribute Attribute `xml:"attribute"`
	Code      Code      `xml:"code"`
}
