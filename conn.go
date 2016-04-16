package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
)

// IgnoreEOF returns err unless err == io.EOF,
// in which case it returns nil.
func IgnoreEOF(err error) error {
	if err == io.EOF {
		return nil
	}
	return err
}

// Conn represents a single connection to an EPP server.
// This implementation is not safe for concurrent use.
type Conn struct {
	net.Conn
	buf     bytes.Buffer
	decoder *xml.Decoder
	saved   xml.Decoder

	// Greeting holds the last received greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	Greeting
}

// NewConn initializes an epp.Conn from a net.Conn and performs the EPP
// handshake. It reads and stores the initial EPP <greeting> message.
// https://tools.ietf.org/html/rfc5730#section-2.4
func NewConn(conn net.Conn) (*Conn, error) {
	c := newConn(conn)
	g, err := c.readGreeting()
	if err == nil {
		c.Greeting = g
	}
	return c, err
}

// Close sends an EPP <logout> command and closes the connection c.
func (c *Conn) Close() error {
	c.Logout()
	return c.Conn.Close()
}

// newConn initializes an epp.Conn from a net.Conn.
// Used internally for testing.
func newConn(conn net.Conn) *Conn {
	c := Conn{Conn: conn}
	c.decoder = xml.NewDecoder(&c.buf)
	c.saved = *c.decoder
	return &c
}

// reset resets the underlying xml.Decoder and bytes.Buffer,
// restoring the original state of the underlying
// xml.Decoder (pos 1, line 1, stack, etc.) using a hack.
func (c *Conn) reset() {
	c.buf.Reset()
	*c.decoder = c.saved // Heh.
}

// flushDataUnit writes bytes from c.buf to c using writeDataUnit.
func (c *Conn) flushDataUnit() error {
	return c.writeDataUnit(c.buf.Bytes())
}

// writeDataUnit writes a slice of bytes to c.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func (c *Conn) writeDataUnit(x []byte) error {
	logXML("<-- WRITE DATA UNIT -->", x)
	s := uint32(4 + len(x))
	err := binary.Write(c.Conn, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(x)
	return err
}

// readResponse reads a single EPP response from c and parses the XML into req.
// It returns an error if the EPP response contains an error Result.
func (c *Conn) readResponse(res *response_) error {
	err := c.readDataUnit()
	if err != nil {
		return err
	}
	err = IgnoreEOF(scanResponse.Scan(c.decoder, res))
	if err != nil {
		return err
	}
	if res.Result.IsError() {
		return res.Result
	}
	return nil
}

// readDataUnit reads a single EPP message from c into
// c.buf. The bytes in c.buf are valid until the next
// call to readDataUnit.
func (c *Conn) readDataUnit() error {
	c.reset()
	var s int32
	err := binary.Read(c.Conn, binary.BigEndian, &s)
	if err != nil {
		return err
	}
	s -= 4
	if s < 0 {
		return io.ErrUnexpectedEOF
	}
	lr := io.LimitedReader{R: c.Conn, N: int64(s)}
	n, err := c.buf.ReadFrom(&lr)
	if err != nil {
		return err
	}
	if n != int64(s) || lr.N != 0 {
		return io.ErrUnexpectedEOF
	}
	logXML("<-- READ DATA UNIT -->", c.buf.Bytes())
	return nil
}

func deleteRange(s, pfx, sfx []byte) []byte {
	start := bytes.Index(s, pfx)
	if start < 0 {
		return s
	}
	end := bytes.Index(s[start+len(pfx):], sfx)
	if end < 0 {
		return s
	}
	end += start + len(pfx) + len(sfx)
	size := len(s) - (end - start)
	copy(s[start:size], s[end:])
	return s[:size]
}

func deleteBufferRange(buf *bytes.Buffer, pfx, sfx []byte) {
	v := deleteRange(buf.Bytes(), pfx, sfx)
	buf.Truncate(len(v))
}
