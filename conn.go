package epp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync/atomic"
)

// Conn represents a single connection to an EPP server.
// This implementation is not safe for concurrent use.
type Conn struct {
	net.Conn
	buf     bytes.Buffer
	decoder xmlDecoder
	txnID   uint64

	// Greeting holds the last received greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	Greeting *Greeting
}

// NewConn initializes an epp.Conn from a net.Conn and performs the EPP
// handshake. It reads and stores the initial EPP <greeting> message.
// https://tools.ietf.org/html/rfc5730#section-2.4
func NewConn(conn net.Conn) (*Conn, error) {
	c := newConn(conn)
	var err error
	c.Greeting, err = c.ReadGreeting()
	return c, err
}

// newConn initializes an epp.Conn from a net.Conn.
// Used internally for testing.
func newConn(conn net.Conn) *Conn {
	c := Conn{Conn: conn}
	c.decoder = newXMLDecoder(&c.buf)
	return &c
}

// writeMessage serializes msg into XML and writes it to c.
func (c *Conn) writeMessage(msg *message) error {
	data, err := marshal(msg)
	if err != nil {
		return err
	}
	return c.writeDataUnit(data)
}

// writeDataUnit writes a slice of bytes to c.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func (c *Conn) writeDataUnit(p []byte) error {
	logXML("<-- WRITE DATA UNIT -->", p)
	s := uint32(4 + len(xmlHeader) + len(p))
	err := binary.Write(c.Conn, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(xmlHeader)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(p)
	return err
}

// readMessage reads a single EPP response from c and parses the XML into req.
// It returns an error if the EPP response contains an error result.
func (c *Conn) readMessage(msg *message) error {
	err := c.readDataUnit()
	if err != nil {
		return err
	}
	return c.decode(msg)
}

// readDataUnit reads a single EPP message from c into
// c.buf. The bytes in c.buf are valid until the next
// call to readDataUnit.
func (c *Conn) readDataUnit() error {
	c.buf.Reset()
	var s uint32
	err := binary.Read(c.Conn, binary.BigEndian, &s)
	if err != nil {
		return err
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

// decode decodes an EPP XML message from c.buf into msg,
// returning any EPP protocol-level errors detected in the message.
func (c *Conn) decode(msg *message) error {
	c.decoder.reset()
	err := c.decoder.Decode(msg)
	if err != nil {
		return err
	}
	return detectError(msg)
}

// id returns a zero-padded 16-character hex uint64 transaction ID.
func (c *Conn) id() string {
	return fmt.Sprintf("%016x", atomic.AddUint64(&c.txnID, 1))
}

// TxnID returns the current client transaction ID for c.
// This generally corresponds to the number of commands performed.
func (c *Conn) TxnID() uint64 {
	return c.txnID
}
