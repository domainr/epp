package epp

import "encoding/xml"

type Bool bool

var (
	True  = Bool(true)
	False = Bool(false)
)

// UnmarshalXML impements the xml.Marshaler interface.
// Any tag present with this type = true.
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct{}
	d.DecodeElement(&v, &start)
	*b = true
	return nil
}

// UnmarshalXML impements the xml.Unmarshaler interface.
// Any tag present with this type = true.
func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b {
		e.EncodeToken(start)
		e.EncodeToken(xml.EndElement{Name: start.Name})
	}
	return nil
}

var (
	openAngle  = xml.CharData("<")
	closeAngle = xml.CharData(">")
)
