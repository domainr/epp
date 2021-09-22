package epp

import (
	"encoding/xml"
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
	ln, err := net.Listen("tcp", ":0")
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
		sc := newConn(conn)
		// Respond with greeting
		err = sc.writeDataUnit([]byte(testXMLGreeting))
		st.Assert(t, err, nil)
		var res Response
		// Read logout message
		err = sc.readResponse(&res)
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

func TestConnDecoderReuse(t *testing.T) {
	c := newConn(nil)
	v := struct {
		XMLName struct{} `xml:"hello"`
		Foo     string   `xml:"foo"`
	}{}

	c.reset()
	c.buf.WriteString(`<hello><foo>foo</foo></hello>`)
	st.Expect(t, c.decoder.InputOffset(), int64(0))
	c.decoder.Decode(&v)
	st.Expect(t, v.Foo, "foo")
	st.Expect(t, c.decoder.InputOffset(), int64(29))

	c.reset()
	c.buf.WriteString(`<hello><foo>bar</foo></hello>`)
	st.Expect(t, c.decoder.InputOffset(), int64(0))
	tok, _ := c.decoder.Token()
	se := tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "hello")
	tok, _ = c.decoder.Token()
	se = tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "foo")
	st.Expect(t, c.decoder.InputOffset(), int64(12))

	c.reset()
	c.buf.WriteString(`<hello><foo>blam&lt;</foo></hello>`)
	st.Expect(t, c.decoder.InputOffset(), int64(0))
	c.decoder.Decode(&v)
	st.Expect(t, v.Foo, "blam<")
	st.Expect(t, c.decoder.InputOffset(), int64(34))
}

func TestDeleteRange(t *testing.T) {
	v := deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`<baz`), []byte(`</baz>`))
	st.Expect(t, string(v), `<foo><bar></bar></foo>`)

	v = deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`</bar>`), []byte(`o>`))
	st.Expect(t, string(v), `<foo><bar><baz></baz>`)
}
