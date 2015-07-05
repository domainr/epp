package epp

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
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

func logMarshal(t *testing.T, v interface{}) {
	x, err := xml.Marshal(v)
	st.Expect(t, err, nil)
	t.Logf("<!-- MARSHALED -->\n%s\n", string(x))
}

func TestConnDecodeMessage(t *testing.T) {
	c := newConn(nil)
	c.buf.WriteString(testXMLDomainCheckResponse)
	var msg message
	err := c.decodeMessage(&msg)
	st.Expect(t, err, nil)
	st.Reject(t, msg.Response, nil)
	st.Expect(t, len(msg.Response.Results), 1)
	logMarshal(t, &msg)
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
