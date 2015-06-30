package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
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
	decoder Decoder
	encoder *xml.Encoder
	txnID   uint64

	// Greeting holds the last received greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	Greeting
}

// NewConn initializes an epp.Conn from a net.Conn and performs the EPP
// handshake. It reads and stores the initial EPP <greeting> message.
// https://tools.ietf.org/html/rfc5730#section-2.4
func NewConn(conn net.Conn) (*Conn, error) {
	c := newConn(conn)
	err := c.readGreeting()
	return c, err
}

// newConn initializes an epp.Conn from a net.Conn.
// Used internally for testing.
func newConn(conn net.Conn) *Conn {
	c := Conn{Conn: conn}
	c.decoder = NewDecoder(&c.buf)
	c.encoder = xml.NewEncoder(&c.buf)
	return &c
}

// writeMessage serializes msg into XML and writes it to c.
// It reuses (and therefore clobbers) the internal buffer
// shared with the parsing side.
func (c *Conn) writeMessage(msg *message) error {
	c.buf.Reset()
	c.buf.Write(xmlHeader)
	err := c.encoder.Encode(msg)
	if err != nil {
		return err
	}
	return c.writeDataUnit(c.buf.Bytes())
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

// readMessage reads a single EPP response from c and parses the XML into req.
// It returns an error if the EPP response contains an error result.
func (c *Conn) readMessage(msg *message) error {
	err := c.readDataUnit()
	if err != nil {
		return err
	}
	return c.decoder.DecodeMessage(msg)
}

// readDataUnit reads a single EPP message from c into
// c.buf. The bytes in c.buf are valid until the next
// call to readDataUnit.
func (c *Conn) readDataUnit() error {
	c.buf.Reset()
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

// id returns a zero-padded 16-character hex uint64 transaction ID.
func (c *Conn) id() string {
	return fmt.Sprintf("%016x", atomic.AddUint64(&c.txnID, 1))
}

// encodeID writes the XML for the transaction ID to c.buf.
func (c *Conn) encodeID() error {
	_, err := fmt.Fprintf(&c.buf, "<clTRID>%016x</clTRID>", atomic.AddUint64(&c.txnID, 1))
	return err
}

// TxnID returns the current client transaction ID for c.
// This generally corresponds to the number of commands performed.
func (c *Conn) TxnID() uint64 {
	return c.txnID
}
