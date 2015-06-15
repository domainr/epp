package epp

import (
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
	// Greeting holds the last received greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	Greeting *Greeting

	net.Conn
	txnID uint64
}

// NewConn initializes an epp.Conn from a net.Conn and performs the EPP
// handshake. It reads and stores the initial EPP <greeting> message.
// https://tools.ietf.org/html/rfc5730#section-2.4
func NewConn(conn net.Conn) (*Conn, error) {
	c := &Conn{Conn: conn}
	var err error
	c.Greeting, err = c.ReadGreeting()
	return c, err
}

// writeMessage serializes msg into XML and writes it to c.
func (c *Conn) writeMessage(msg *message) error {
	data, err := marshal(msg)
	if err != nil {
		return err
	}
	logRequest(data)
	return c.writeDataUnit(data)
}

// writeDataUnit writes a slice of bytes to c.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func (c *Conn) writeDataUnit(p []byte) error {
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

// xmlHeader is a byte-slice representation of the
// standard XML header. Declared as a global to relieve GC pressure.
var xmlHeader = []byte(xml.Header)

// readMessage reads a single EPP response from c and parses the XML into req.
// It returns an error if the EPP response contains an error result.
func (c *Conn) readMessage(msg *message) error {
	data, err := c.readDataUnit()
	if err != nil {
		return err
	}
	logResponse(data)
	return unmarshal(data, msg)
}

// readDataUnit reads a single EPP message from c.
// It returns the bytes read and/or an error.
// FIXME: allocate a single buffer per Conn to reduce GC pressure?
func (c *Conn) readDataUnit() (data []byte, err error) {
	var s uint32
	err = binary.Read(c.Conn, binary.BigEndian, &s)
	if err != nil {
		return
	}
	data = make([]byte, s-4)
	n, err := c.Conn.Read(data)
	if err != nil {
		return
	}
	if n != len(data) {
		return data, io.ErrNoProgress
	}
	return data, nil
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
