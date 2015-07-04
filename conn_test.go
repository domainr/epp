package epp

import (
	"bytes"
	"crypto/tls"
	"net"
	"testing"

	"github.com/nbio/st"
)

const (
	// https://wiki.hexonet.net/wiki/Domain_API
	addr     = "api.1api.net:1700"
	user     = "test.user"
	password = "test.passw0rd"
)

func testDial(t *testing.T) net.Conn {
	if testing.Short() {
		t.Skip("network-dependent")
	}
	conn, err := tls.Dial("tcp", addr, nil)
	st.Assert(t, err, nil)
	return conn
}

func testLogin(t *testing.T) *Conn {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
	return c
}

func TestNewConn(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Expect(t, err, nil)
	st.Reject(t, c, nil)
	st.Reject(t, c.Greeting.ServerName, "")
}

func TestDeleteRange(t *testing.T) {
	v := deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`<baz`), []byte(`</baz>`))
	st.Expect(t, string(v), `<foo><bar></bar></foo>`)

	v = deleteRange([]byte(`<foo><bar><baz></baz></bar></foo>`), []byte(`</bar>`), []byte(`o>`))
	st.Expect(t, string(v), `<foo><bar><baz></baz>`)
}

func TestDeleteBufferRange(t *testing.T) {
	buf := bytes.NewBufferString(`<foo><bar><baz></baz></bar></foo>`)
	deleteBufferRange(buf, []byte(`<baz`), []byte(`</baz>`))
	st.Expect(t, buf.String(), `<foo><bar></bar></foo>`)

	buf = bytes.NewBufferString(`<foo><bar><baz></baz></bar></foo>`)
	deleteBufferRange(buf, []byte(`</bar>`), []byte(`o>`))
	st.Expect(t, buf.String(), `<foo><bar><baz></baz>`)
}
