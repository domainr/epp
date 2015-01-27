package epp

import (
	"crypto/tls"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
	"strconv"
	"sync/atomic"

	"github.com/wsxiaoys/terminal/color"
)

// Conn represents a single connection to an EPP server.
// This implementation is not safe for concurrent use.
type Conn struct {
	// Greeting holds the last recieved greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	Greeting *Greeting

	net.Conn
	txnID uint64
}

// Dial connects to an EPP server via TCP.
// Returns an error if unable to connect, including certificate mismatch errors.
func Dial(addr string) (*Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return newConn(conn)
}

// DialTLS connects to an EPP server via TLS.
// Returns an error if unable to connect, including certificate mismatch errors.
func DialTLS(addr string, cfg *tls.Config) (*Conn, error) {
	conn, err := tls.Dial("tcp", addr, cfg)
	if err != nil {
		return nil, err
	}
	return newConn(conn)
}

func newConn(conn net.Conn) (*Conn, error) {
	c := &Conn{Conn: conn}
	err := c.readGreeting()
	if err != nil {
		c.Close()
		return nil, err
	}
	return c, nil
}

// WriteMessage serializes msg into XML and writes it to c.
func (c *Conn) WriteMessage(msg interface{}) error {
	data, err := xml.Marshal(msg)
	if err != nil {
		return err
	}
	color.Printf("@{|}<!-- REQUEST -->\n%s\n\n", string(data))
	return c.WriteDataUnit(data)
}

// WriteDataUnit writes a slice of bytes to c.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func (c *Conn) WriteDataUnit(p []byte) error {
	s := uint32(4 + len(xml.Header) + len(p))
	err := binary.Write(c.Conn, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = c.Conn.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	_, err = c.Conn.Write(p)
	return err
}

// ReadResponse reads a single EPP message from c and parses the XML into msg.
// It returns an error if the EPP message contains an error result.
func (c *Conn) ReadResponse(rmsg *Response) error {
	data, err := c.ReadDataUnit()
	if err != nil {
		return err
	}
	color.Printf("@{c}<!-- RESPONSE -->\n%s\n\n", string(data))
	err = xml.Unmarshal(data, rmsg)
	if err != nil {
		return err
	}
	// color.Fprintf(os.Stderr, "@{y}%s\n", spew.Sprintf("%+v", msg))
	if len(rmsg.Results) != 0 {
		r := rmsg.Results[0]
		if r.IsError() {
			return r
		}
	}
	return nil
}

// ReadDataUnit reads a single EPP message from c.
// It returns the bytes read and/or an error.
// FIXME: allocate a single buffer per Conn to reduce GC pressure?
func (c *Conn) ReadDataUnit() (data []byte, err error) {
	var s uint32
	err = binary.Read(c.Conn, binary.BigEndian, &s)
	if err != nil {
		return
	}
	data = make([]byte, s)
	n, err := c.Conn.Read(data)
	if err != nil {
		return
	}
	if 4+n != int(s) {
		return data, io.ErrNoProgress
	}
	return data, nil
}

func (c *Conn) id() string {
	return strconv.FormatUint(atomic.AddUint64(&c.txnID, 1), 16)
}
