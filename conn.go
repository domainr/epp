package epp

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"io"
	"net"
	"sync"
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
// Reads and writes are serialized, so it is safe for concurrent use.
type Conn struct {
	// Conn is the underlying net.Conn (usually a TLS connection).
	net.Conn

	// Timeout defines the timeout for network operations.
	// It must be set at initialization. Changing it after
	// a connection is already opened will have no effect.
	Timeout time.Duration

	// m protects Greeting and LoginResult.
	m sync.Mutex

	// Greeting holds the last received greeting message from the server,
	// indicating server name, status, data policy and capabilities.
	//
	// Deprecated: This field is written to upon opening a new EPP connection and should not be modified.
	Greeting

	// LoginResult holds the last received login response message's Result
	// from the server, in which some servers might include diagnostics such
	// as connection count limits.
	//
	// Deprecated: this field is written to by the Login method but otherwise is not used by this package.
	LoginResult Result

	requests  chan []byte
	responses chan *Response

	done     chan struct{}
	readErr  error
	writeErr error
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
	select {
	case <-c.done:
		return net.ErrClosed
	default:
	}
	c.Logout()
	close(c.done)
	return c.Conn.Close()
}

// newConn initializes an epp.Conn from a net.Conn.
// Used internally for testing.
func newConn(conn net.Conn) *Conn {
	c := Conn{
		Conn:      conn,
		requests:  make(chan []byte),
		responses: make(chan *Response),
		done:      make(chan struct{}),
	}
	go c.writeLoop()
	go c.readLoop()
	return &c
}

// writeRequest queues a single EPP request (x) for writing on c.
// It returns net.ErrClosed if the underlying connection is closed.
// writeRequest can be called from multiple goroutines.
func (c *Conn) writeRequest(x []byte) error {
	select {
	case c.requests <- x:
		return nil
	case <-c.done:
		return net.ErrClosed
	}
}

// readResponse dequeues and returns a EPP response from c.
// It returns an error if the EPP response contains an error Result.
// readResponse can be called from multiple goroutines.
func (c *Conn) readResponse() (*Response, error) {
	select {
	case res := <-c.responses:
		if res == nil {
			return res, c.readErr
		}
		if res.Result.IsError() {
			return nil, &res.Result
		}
		return res, nil
	case <-c.done:
		return nil, net.ErrClosed
	}
}

func (c *Conn) writeLoop() {
	defer c.Close()
	for {
		// TODO(ydnar): figure out how to handle timeouts for continous write loop.
		select {
		case x := <-c.requests:
			err := writeDataUnit(c.Conn, x)
			if err != nil {
				c.writeErr = err
				return
			}
		case <-c.done:
			return
		}
	}
}

func (c *Conn) readLoop() {
	defer close(c.responses)
	buf := &bytes.Buffer{}
	decoder := xml.NewDecoder(buf)
	saved := *decoder
	for {
		// TODO(ydnar): figure out how to handle timeouts for continous read loop.
		// if c.Timeout > 0 {
		// 	c.Conn.SetReadDeadline(time.Now().Add(c.Timeout))
		// }
		err := readDataUnit(buf, c.Conn)
		if err != nil {
			c.readErr = err
			return
		}
		res := &Response{}
		*decoder = saved
		err = IgnoreEOF(scanResponse.Scan(decoder, res))
		if err != nil {
			c.readErr = err
			return
		}
		c.responses <- res
	}
}

// writeDataUnit writes x to w.
// Bytes written are prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func writeDataUnit(w io.Writer, x []byte) error {
	logXML("<-- WRITE DATA UNIT -->", x)
	s := uint32(4 + len(x))
	err := binary.Write(w, binary.BigEndian, s)
	if err != nil {
		return err
	}
	_, err = w.Write(x)
	return err
}

// readDataUnit reads a single EPP data unit from r, returning the payload bytes or an error.
// An EPP data units is prefixed with 32-bit header specifying the total size
// of the data unit (message + 4 byte header), in network (big-endian) order.
// http://www.ietf.org/rfc/rfc4934.txt
func readDataUnit(buf *bytes.Buffer, r io.Reader) error {
	var s int32
	err := binary.Read(r, binary.BigEndian, &s)
	if err != nil {
		return err
	}
	s -= 4 // https://tools.ietf.org/html/rfc5734#section-4
	if s < 0 {
		return io.ErrUnexpectedEOF
	}
	buf.Reset()
	buf.Grow(int(s))
	_, err = io.CopyN(buf, r, int64(s))
	logXML("<-- READ DATA UNIT -->", buf.Bytes())
	return err
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
