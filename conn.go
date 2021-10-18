package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
	"time"
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

	// LoginResult holds the last received login response message's Result
	// from the server, in which some servers might include diagnostics such
	// as connection count limits.
	LoginResult Result

	// Timeout defines the timeout for network operations.
	Timeout time.Duration
}

// NewConn initializes an epp.Conn from a net.Conn and performs the EPP
// handshake. It reads and stores the initial EPP <greeting> message.
// https://tools.ietf.org/html/rfc5730#section-2.4
func NewConn(conn net.Conn) (*Conn, error) {
	return NewTimeoutConn(conn, 0)
}

// NewTimeoutConn initializes an epp.Conn like NewConn, limiting the duration of network
// operations on conn using Set(Read|Write)Deadline.
func NewTimeoutConn(conn net.Conn, timeout time.Duration) (*Conn, error) {
	c := newConn(conn)
	c.Timeout = timeout
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

// writeDataUnit writes a slice of bytes to c.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func (c *Conn) writeDataUnit(x []byte) error {
	logXML("<-- WRITE DATA UNIT -->", x)
	s := uint32(4 + len(x))
	if c.Timeout > 0 {
		c.Conn.SetWriteDeadline(time.Now().Add(c.Timeout))
	}
	err := binary.Write(c.Conn, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(x)
	return err
}

// readResponse reads a single EPP response from c and parses the XML into req.
// It returns an error if the EPP response contains an error Result.
func (c *Conn) readResponse(res *Response) error {
	err := c.readDataUnit()
	if err != nil {
		return err
	}
	err = IgnoreEOF(scanResponse.Scan(c.decoder, res))
	if err != nil {
		return err
	}
	if res.Result.IsError() {
		return &res.Result
	}
	return nil
}

// readDataUnit reads a single EPP message from c into
// c.buf. The bytes in c.buf are valid until the next
// call to readDataUnit.
func (c *Conn) readDataUnit() error {
	c.reset()
	var s int32
	if c.Timeout > 0 {
		c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
	}
	err := binary.Read(c.Conn, binary.BigEndian, &s)
	if err != nil {
		return err
	}
	s -= 4 // https://tools.ietf.org/html/rfc5734#section-4
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
