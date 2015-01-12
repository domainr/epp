package epp

import (
	"encoding/xml"
	"time"
)

// Time represents EPP date-time values.
type Time struct {
	time.Time
}

// UnmarshalXML implements a custom XML unmarshaler.
// http://stackoverflow.com/a/25015260
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	parse, err := time.Parse(time.RFC3339, v)
	if err != nil {
		return nil
	}
	*t = Time{parse}
	return nil
}
