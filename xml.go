package epp

import (
	"encoding/xml"
	"io"
)

// xmlHeader is a byte-slice representation of the standard XML header.
// Declared as a global to relieve GC pressure.
var xmlHeader = []byte(xml.Header)

// xmlDecoder implements a resettable XML decoder.
// This is a dirty hack to reduce GC pressure.
type xmlDecoder struct {
	xml.Decoder
	saved xml.Decoder
}

// newXMLDecoder returns an initialized xmlDecoder.
// The initial state of the xml.Decoder is copied to saved.
func newXMLDecoder(r io.Reader) xmlDecoder {
	d := xml.NewDecoder(r)
	return xmlDecoder{*d, *d}
}

// reset restores the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.).
func (d *xmlDecoder) reset() {
	d.Decoder = d.saved
}
