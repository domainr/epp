package epp

import (
	"encoding/xml"
	"fmt"
	"io"
)

// xmlHeader is a byte-slice representation of the standard XML header.
// Declared as a global to relieve GC pressure.
var xmlHeader = []byte(xml.Header)

// Decoder implements a resettable XML decoder.
// This is a dirty hack to reduce GC pressure.
type Decoder struct {
	xml.Decoder
	saved xml.Decoder
	Stack []xml.StartElement
}

// NewDecoder returns an initialized decoder.
// The initial state of the xml.Decoder is copied to saved.
func NewDecoder(r io.Reader) Decoder {
	d := xml.NewDecoder(r)
	return Decoder{Decoder: *d, saved: *d}
}

// Reset restores the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.).
func (d *Decoder) Reset() {
	d.Decoder = d.saved
	d.Stack = d.Stack[:0]
}

// DecodeMessage decodes an EPP XML message into msg,
// returning any EPP protocol-level errors detected in the message.
// It resets the underlying xml.Decoder before attempting to decode
// the input stream.
func (d *Decoder) DecodeMessage(msg *message) error {
	d.Reset()
	err := d.Decode(msg)
	if err != nil {
		return err
	}
	return msg.error()
}

// Token returns an xml.Token from its internal xml.Decoder or an error.
// It maintains a stack of xml.StartElements.
func (d *Decoder) Token() (xml.Token, error) {
	t, err := d.Decoder.Token()
	switch node := t.(type) {
	case xml.StartElement:
		d.Stack = append(d.Stack, node)
	case xml.EndElement:
		n := len(d.Stack)
		if n == 0 || node.Name != d.Stack[n-1].Name {
			return t, fmt.Errorf("unbalanced end tag: %s", node.Name)
		}
		d.Stack = d.Stack[:n-1]
	}
	return t, err
}
