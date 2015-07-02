package epp

import (
	"encoding/xml"
	"fmt"
	"io"
)

// Decoder implements a resettable XML decoder.
// This is a dirty hack to reduce GC pressure.
type Decoder struct {
	decoder xml.Decoder
	saved   xml.Decoder
	Stack   []xml.StartElement
}

// NewDecoder returns an initialized decoder.
// The initial state of the xml.Decoder is copied to saved.
func NewDecoder(r io.Reader) Decoder {
	d := xml.NewDecoder(r)
	return Decoder{decoder: *d, saved: *d}
}

// Element returns StartElement indexed by i.
// Indexes < 0 are offset from len(d.stack).
// Returns a zero-value StartElement if i is out of bounds, or
// the decoder is not inside an XML tag.
func (d *Decoder) Element(i int) xml.StartElement {
	if i < 0 {
		i += len(d.Stack)
	}
	if i < 0 || i >= len(d.Stack) {
		return xml.StartElement{}
	}
	return d.Stack[i]
}

// AtPath determines if the current local element
// path ends with path.
func (d *Decoder) AtPath(path ...string) bool {
	ls := len(d.Stack)
	lp := len(path)
	if ls < lp {
		return false
	}
	for i := 1; i <= lp; i++ {
		if d.Stack[ls-i].Name.Local != path[lp-i] {
			return false
		}
	}
	return true
}

// Reset restores the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.).
func (d *Decoder) Reset() {
	d.decoder = d.saved
	d.Stack = d.Stack[:0]
}

// DecodeMessage decodes an EPP XML message into msg,
// returning any EPP protocol-level errors detected in the message.
// It resets the underlying xml.Decoder before attempting to decode
// the input stream.
func (d *Decoder) DecodeMessage(msg *message) error {
	d.Reset()
	err := d.decoder.Decode(msg)
	if err != nil {
		return err
	}
	return msg.error()
}

// Token returns an xml.Token from its internal xml.Decoder or an error.
// It maintains a stack of xml.StartElements.
func (d *Decoder) Token() (xml.Token, error) {
	t, err := d.decoder.Token()
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

// DecodeWith is an experimental function to wrap the underlying
// for loop + switch pattern when using xml.Decoder.
// If a call to f returns an error, it exits early, returning that error.
// It does not reset the decoder.
func (d *Decoder) DecodeWith(f func(xml.Token) error) error {
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		err = f(t)
		if err != nil {
			return err
		}
	}
	return nil
}
