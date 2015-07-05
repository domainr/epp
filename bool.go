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

// UnmarshalXML impements the xml.Unmarshaler interface.
// Any tag present with this type = true.
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct{}
	d.DecodeElement(&v, &start)
	*b = true
	return nil
}

// UnmarshalXMLAttr impements the xml.MarshalerAttr interface.
// A value of 1 or starting with a t or T is considered true.
func (b *Bool) UnmarshalXMLAttr(attr *xml.Attr) error {
	if attr.Value == "1" || attr.Value[0] == 'T' || attr.Value[0] == 't' {
		*b = true
	} else {
		*b = false
	}
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

// MarshalXMLAttr implements the xml.MarshalerAttr interface.
// Attributes will be serialized with a value of "0" or "1".
func (b Bool) MarshalXMLAttr(name xml.Name) (attr xml.Attr, err error) {
	attr.Name = name
	if b {
		attr.Value = "1"
	} else {
		attr.Value = "0"
	}
	return
}
