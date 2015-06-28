package epp

import (
	"encoding/xml"
	"fmt"
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
	stack []xml.StartElement
}

// newDecoder returns an initialized decoder.
// The initial state of the xml.Decoder is copied to saved.
func newDecoder(r io.Reader) decoder {
	d := xml.NewDecoder(r)
	return decoder{Decoder: *d, saved: *d}
}

// reset restores the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.).
func (d *decoder) reset() {
	d.Decoder = d.saved
	d.stack = d.stack[:0]
}

// decode decodes an EPP XML message into msg,
// returning any EPP protocol-level errors detected in the message.
// It resets the underlying xml.Decoder before attempting to decode
// the input stream.
func (d *decoder) decode(msg *message) error {
	d.reset()
	err := d.Decode(msg)
	if err != nil {
		return err
	}
	return msg.error()
}

// Token returns an xml.Token from its internal xml.Decoder or an error.
// It maintains a stack of xml.StartElements.
func (d *decoder) Token() (xml.Token, error) {
	t, err := d.Decoder.Token()
	switch node := t.(type) {
	case xml.StartElement:
		d.stack = append(d.stack, node)
	case xml.EndElement:
		n := len(d.stack)
		if n == 0 || node.Name != d.stack[n-1].Name {
			return t, fmt.Errorf("unbalanced end tag: %s", node.Name)
		}
		d.stack = d.stack[:n-1]
	}
	return t, err
}
