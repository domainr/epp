package epp

import (
	"encoding/binary"
	"io"
	"net"
)

// Transport is an interface that can read and write EPP data units.
type Transport interface {
	ReadDataUnit() ([]byte, error)
	WriteDataUnit([]byte) error
}

// IO implements Transport using an io.Reader and an io.Writer.
type IO struct {
	R io.Reader
	W io.Writer
}

var _ Transport = &IO{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (t *IO) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(t.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *IO) WriteDataUnit(b []byte) error {
	return WriteDataUnit(t.W, b)
}

// Conn implements Transport using a net.Conn.
type Conn struct {
	net.Conn
}

var _ Transport = &Conn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (t *Conn) ReadDataUnit() ([]byte, error) {
	return ReadDataUnit(t.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *Conn) WriteDataUnit(b []byte) error {
	return WriteDataUnit(t.Conn, b)
}

// ReadDataUnit reads a single EPP data unit from r, returning the payload or an error.
// An EPP data unit is prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func ReadDataUnit(r io.Reader) ([]byte, error) {
	var n uint32
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}
	// An EPP data unit size includes the 4 byte header.
	// See https://tools.ietf.org/html/rfc5734#section-4.
	if n < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	n -= 4
	b := make([]byte, n)
	_, err = io.ReadAtLeast(r, b, int(n))
	return b, err
}

// WriteDataUnit writes a single EPP data unit to w.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func WriteDataUnit(w io.Writer, b []byte) error {
	s := uint32(4 + len(b))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
