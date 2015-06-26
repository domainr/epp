package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestReuseXMLDecoder(t *testing.T) {
	buf := &bytes.Buffer{}
	d := newXMLDecoder(buf)

	v := struct {
		XMLName struct{} `xml:"hello"`
		Foo     string   `xml:"foo"`
	}{}

	buf.Reset()
	buf.Write([]byte(`<hello><foo>foo</foo></hello>`))
	d.reset()
	st.Expect(t, d.InputOffset(), int64(0))
	d.Decode(&v)
	st.Expect(t, v.Foo, "foo")
	st.Expect(t, d.InputOffset(), int64(29))

	buf.Reset()
	buf.Write([]byte(`<hello><foo>bar</foo></hello>`))
	d.reset()
	st.Expect(t, d.InputOffset(), int64(0))
	tok, _ := d.Token()
	se := tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "hello")
	tok, _ = d.Token()
	se = tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "foo")
	st.Expect(t, d.InputOffset(), int64(12))

	buf.Reset()
	buf.Write([]byte(`<hello><foo>blam&lt;</foo></hello>`))
	d.reset()
	st.Expect(t, d.InputOffset(), int64(0))
	d.Decode(&v)
	st.Expect(t, v.Foo, "blam<")
	st.Expect(t, d.InputOffset(), int64(34))
}
