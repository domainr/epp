package std

import "github.com/nbio/xml"

// Bool represents a boolean value that can be serialized to XML.
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
// TODO: support <tag>false</tag> representation.
func (b *Bool) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct{}
	d.DecodeElement(&v, &start)
	*b = true
	return nil
}

// UnmarshalXMLAttr impements the xml.MarshalerAttr interface.
// An empty value, 0, or starting with a f or F is considered false.
// Any other value is considered true.
func (b *Bool) UnmarshalXMLAttr(attr *xml.Attr) error {
	if len(attr.Value) == 0 || attr.Value == "1" || attr.Value[0] == 'f' || attr.Value[0] == 'F' {
		*b = false
	} else {
		*b = true
	}
	return nil
}

// MarshalXML impements the xml.Marshaler interface.
// Any tag present with this type = true.
func (b Bool) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if b {
		e.EncodeToken(xml.SelfClosingElement(start))
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
