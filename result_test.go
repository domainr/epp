package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestDecodeResult(t *testing.T) {
	var r Result
	var buf bytes.Buffer
	d := NewDecoder(&buf)

	buf.Reset()
	buf.WriteString(`<result code="1000"><msg>Command completed successfully</msg></result>`)
	d.Reset()
	err := IgnoreEOF(decodeResult(&d, &r))
	st.Expect(t, err, nil)
	st.Expect(t, r.Code, 1000)
	st.Expect(t, r.Message, "Command completed successfully")
	st.Expect(t, r.IsError(), false)
	st.Expect(t, r.IsFatal(), false)

	// Result code >= 2000 is an error.
	buf.Reset()
	buf.WriteString(`<result code="2001"><msg>Command syntax error</msg></result>`)
	d.Reset()
	err = decodeResult(&d, &r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2001)
	st.Expect(t, r.Message, "Command syntax error")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), false)

	// Result code > 2500 is a fatal error.
	buf.Reset()
	buf.WriteString(`<result code="2501"><msg>Authentication error; server closing connection</msg></result>`)
	d.Reset()
	err = decodeResult(&d, &r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2501)
	st.Expect(t, r.Message, "Authentication error; server closing connection")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), true)

	// Decoding should stop after </result>.
	buf.Reset()
	buf.WriteString(`<result code="1000"><msg>OK</msg></result><foo></foo>`)
	d.Reset()
	err = decodeResult(&d, &r)
	token, err := d.Token()
	st.Expect(t, err, nil)
	se := token.(xml.StartElement)
	st.Expect(t, se.Name.Local, "foo")
}

func BenchmarkDecodeResult(b *testing.B) {
	b.StopTimer()
	var r Result
	var buf bytes.Buffer
	d := NewDecoder(&buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf.Reset()
		buf.WriteString(`<result code="1000"><msg>Command completed successfully</msg></result>`)
		d.Reset()
		b.StartTimer()
		decodeResult(&d, &r)
	}
}
