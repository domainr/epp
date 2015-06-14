package epp

import "encoding/xml"

// Bool represents a bool that can be serialized to XML.
// True: <tag>
// False: (no tag)
type Bool bool

var (
	// True is a Bool of value true.
	True = Bool(true)

	// False is a Bool of value false.
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

// MarshalXML impements the xml.Marshaler interface.
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
