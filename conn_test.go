package epp

import (
	"bytes"
	"net"
	"sync"
	"testing"

	"github.com/nbio/st"
)

type localServer struct {
	lnmu sync.RWMutex
	net.Listener
	done chan bool // signal that indicates server stopped
}

func (ls *localServer) buildup(handler func(*localServer, net.Listener)) error {
	go func() {
		handler(ls, ls.Listener)
		close(ls.done)
	}()
	return nil
}

func (ls *localServer) teardown() {
	ls.lnmu.Lock()
	defer ls.lnmu.Unlock()
	if ls.Listener != nil {
		ls.Listener.Close()
		<-ls.done
		ls.Listener = nil
	}
}

func newLocalServer() (*localServer, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, err
	}
	return &localServer{Listener: ln, done: make(chan bool)}, nil
}

func TestNewConn(t *testing.T) {
	ls, err := newLocalServer()
	st.Assert(t, err, nil)
	defer ls.teardown()
	ls.buildup(func(ls *localServer, ln net.Listener) {
		conn, err := ls.Accept()
		st.Assert(t, err, nil)
		// Respond with greeting
		err = writeDataUnit(conn, []byte(testXMLGreeting))
		st.Assert(t, err, nil)
		// Read logout message
		_, err = readDataUnitHeader(conn)
		st.Assert(t, err, nil)
		// Close connection
		err = conn.Close()
		st.Assert(t, err, nil)
	})
	nc, err := net.Dial(ls.Listener.Addr().Network(), ls.Listener.Addr().String())
	st.Assert(t, err, nil)
	c, err := NewConn(nc)
	st.Assert(t, err, nil)
	st.Reject(t, c, nil)
	st.Reject(t, c.Greeting.ServerName, "")
	err = c.Close()
	st.Expect(t, err, nil)
}

func TestDeleteRange(t *testing.T) {
	v := deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`<baz`), []byte(`</baz>`))
	st.Expect(t, string(v), `<foo><bar></bar></foo>`)

	v = deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`</bar>`), []byte(`o>`))
	st.Expect(t, string(v), `<foo><bar><baz></baz>`)
}

func deleteBufferRange(buf *bytes.Buffer, pfx, sfx []byte) {
	v := deleteRange(buf.Bytes(), pfx, sfx)
	buf.Truncate(len(v))
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
