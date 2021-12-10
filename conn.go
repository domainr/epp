package epp

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
)

// Conn is a generic connection that can read and write EPP data units.
// Multiple goroutines may invoke methods on a Conn simultaneously.
type Conn interface {
	// ReadDataUnit reads a single EPP data unit, returning the payload bytes or an error.
	ReadDataUnit() ([]byte, error)

	// WriteDataUnit writes a single EPP data unit, returning any error.
	WriteDataUnit([]byte) error

	// Close closes the connection.
	Close() error

	// LocalAddr returns the local network address, if any.
	LocalAddr() net.Addr

	// RemoteAddr returns the remote network address, if any.
	RemoteAddr() net.Addr
}

// Pipe implements Conn using an io.Reader and an io.Writer.
type Pipe struct {
	// R is from by ReadDataUnit.
	R io.Reader

	// W is written to by WriteDataUnit.
	W io.Writer

	r sync.Mutex
	w sync.Mutex
}

var _ Conn = &Pipe{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload bytes or an error.
func (t *Pipe) ReadDataUnit() ([]byte, error) {
	t.r.Lock()
	defer t.r.Unlock()
	return ReadDataUnit(t.R)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *Pipe) WriteDataUnit(data []byte) error {
	t.w.Lock()
	defer t.w.Unlock()
	return WriteDataUnit(t.W, data)
}

// Close attempts to close both the underlying reader and writer.
// It will return the first error encountered.
func (t *Pipe) Close() error {
	var rerr, werr error
	if c, ok := t.R.(io.Closer); ok {
		rerr = c.Close()
	}
	if r, ok := t.W.(io.Reader); ok && r == t.R {
		return rerr
	}
	if c, ok := t.W.(io.Closer); ok {
		werr = c.Close()
	}
	if rerr != nil {
		return rerr
	}
	return werr
}

// LocalAddr attempts to return the local address of p.
// If p.R implements LocalAddr, it will be called.
// Otherwise, LocalAddr will return nil.
func (p *Pipe) LocalAddr() net.Addr {
	if a, ok := p.R.(interface{ LocalAddr() net.Addr }); ok {
		return a.LocalAddr()
	}
	return nil
}

// RemoteAddr attempts to return the remote address of p.
// If p.W implements RemoteAddr, it will be called.
// Otherwise, RemoteAddr will return nil.
func (p *Pipe) RemoteAddr() net.Addr {
	if a, ok := p.W.(interface{ RemoteAddr() net.Addr }); ok {
		return a.RemoteAddr()
	}
	return nil
}

// NetConn implements Conn using a net.Conn.
type NetConn struct {
	net.Conn
	r sync.Mutex
	w sync.Mutex
}

var _ Conn = &NetConn{}

// ReadDataUnit reads a single EPP data unit from t, returning the payload or an error.
func (t *NetConn) ReadDataUnit() ([]byte, error) {
	t.r.Lock()
	defer t.r.Unlock()
	return ReadDataUnit(t.Conn)
}

// WriteDataUnit writes a single EPP data unit to t or returns an error.
func (t *NetConn) WriteDataUnit(p []byte) error {
	t.w.Lock()
	defer t.w.Unlock()
	return WriteDataUnit(t.Conn, p)
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
	p := make([]byte, n)
	_, err = io.ReadAtLeast(r, p, int(n))
	return p, err
}

// WriteDataUnit writes a single EPP data unit to w.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// See http://www.ietf.org/rfc/rfc4934.txt for more information.
func WriteDataUnit(w io.Writer, p []byte) error {
	s := uint32(4 + len(p))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(p)
	return err
}
