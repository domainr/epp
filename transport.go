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

// IOTransport implements Transport using an io.Reader and an io.Writer.
type IOTransport struct {
	R io.Reader
	W io.Writer
}

var _ Transport = &IOTransport{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (t *IOTransport) ReadDataUnit() ([]byte, error) {
	return readDataUnit(t.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *IOTransport) WriteDataUnit(b []byte) error {
	return writeDataUnit(t.W, b)
}

// NetTransport implements Transport using a net.Conn.
type NetTransport struct {
	Conn net.Conn
}

var _ Transport = &NetTransport{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (t *NetTransport) ReadDataUnit() ([]byte, error) {
	return readDataUnit(t.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *NetTransport) WriteDataUnit(b []byte) error {
	return writeDataUnit(t.Conn, b)
}

// readDataUnit reads a single EPP data unit from r, returning the payload or an error.
// An EPP data unit is prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func readDataUnit(r io.Reader) ([]byte, error) {
	var n uint32
	err := binary.Read(r, binary.BigEndian, &n)
	if err != nil {
		return nil, err
	}
	if n < 4 {
		return nil, io.ErrUnexpectedEOF
	}
	// An EPP data unit size includes the 4 byte header.
	// See https://tools.ietf.org/html/rfc5734#section-4.
	b := make([]byte, n-4)
	_, err = io.ReadAtLeast(r, b, int(n))
	return b, err
}

// writeDataUnit writes a single EPP data unit to w.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func writeDataUnit(w io.Writer, b []byte) error {
	s := uint32(4 + len(b))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}
