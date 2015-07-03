package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestMarshalOmitEmpty(t *testing.T) {
	v := struct {
		XMLName struct{} `xml:"hello"`
		Foo     string   `xml:"foo"`
		Bar     struct {
			Baz string `xml:"baz"`
		} `xml:"bar,omitempty"`
	}{}

	x, err := xml.Marshal(&v)
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<hello><foo></foo><bar><baz></baz></bar></hello>`)
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
