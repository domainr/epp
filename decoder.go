package epp

import (
	"encoding/xml"
	"io"
)

// Decoder implements a resettable XML decoder.
// This is a dirty hack to reduce GC pressure.
type Decoder struct {
	xml.Decoder
	saved xml.Decoder
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
}

// DecodeMessage decodes an EPP XML message into msg,
// returning any EPP protocol-level errors detected in the message.
// It resets the underlying xml.Decoder before attempting to decode
// the input stream.
func (d *Decoder) DecodeMessage(msg *message) error {
	d.Reset()
	err := d.Decoder.Decode(msg)
	if err != nil {
		return err
	}
	return msg.error()
}
