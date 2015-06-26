package epp

import (
	"encoding/xml"
	"io"
)

// xmlHeader is a byte-slice representation of the standard XML header.
// Declared as a global to relieve GC pressure.
var xmlHeader = []byte(xml.Header)

// decoder implements a resettable XML decoder.
// This is a dirty hack to reduce GC pressure.
type decoder struct {
	xml.Decoder
	saved xml.Decoder
}

// newDecoder returns an initialized decoder.
// The initial state of the xml.Decoder is copied to saved.
func newDecoder(r io.Reader) decoder {
	d := xml.NewDecoder(r)
	return decoder{*d, *d}
}

// reset restores the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.).
func (d *decoder) reset() {
	d.Decoder = d.saved
}

// decode decodes an EPP XML message from c.buf into msg,
// returning any EPP protocol-level errors detected in the message.
func (d *decoder) decode(msg *message) error {
	d.reset()
	err := d.Decode(msg)
	if err != nil {
		return err
	}
	return msg.error()
}
