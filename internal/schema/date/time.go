package date

import (
	"time"

	"github.com/nbio/xml"
)

// Time represents EPP date-time values, serialized to XML in RFC 3339 format.
// Because the default encoding.TextMarshaler implementation in time.Time uses
// RFC 3339, we don’t need to create a custom marshaler for this type.
// See https://www.rfc-editor.org/rfc/rfc3339.html.
type Time struct {
	time.Time
}

// Pointer returns a pointer to a Time struct.
func Pointer(tt time.Time) *Time {
	return &Time{tt}
}

// UnmarshalXML implements a custom XML unmarshaler that ignores parsing errors.
// See http://stackoverflow.com/a/25015260.
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if tt, err := time.Parse(time.RFC3339, v); err == nil {
		*t = Time{tt}
	}
	return nil
}
