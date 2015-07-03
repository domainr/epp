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
	err := IgnoreEOF(d.decodeResult(&r))
	st.Expect(t, err, nil)
	st.Expect(t, r.Code, 1000)
	st.Expect(t, r.Message, "Command completed successfully")
	st.Expect(t, r.IsError(), false)
	st.Expect(t, r.IsFatal(), false)

	// Result code >= 2000 is an error.
	buf.Reset()
	buf.WriteString(`<result code="2001"><msg>Command syntax error</msg></result>`)
	d.Reset()
	err = d.decodeResult(&r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2001)
	st.Expect(t, r.Message, "Command syntax error")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), false)

	// Result code > 2500 is a fatal error.
	buf.Reset()
	buf.WriteString(`<result code="2501"><msg>Authentication error; server closing connection</msg></result>`)
	d.Reset()
	err = d.decodeResult(&r)
	st.Expect(t, err, &r)
	st.Expect(t, r.Code, 2501)
	st.Expect(t, r.Message, "Authentication error; server closing connection")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), true)

	// Decoding should stop after </result>.
	buf.Reset()
	buf.WriteString(`<result code="1000"><msg>OK</msg></result><foo></foo>`)
	d.Reset()
	err = d.decodeResult(&r)
	token, err := d.Token()
	st.Expect(t, err, nil)
	se := token.(xml.StartElement)
	st.Expect(t, se.Name.Local, "foo")
}

func TestScanResult(t *testing.T) {
	var res response_

	d := db(`<epp><response><result code="1000"><msg>Command completed successfully</msg></result></response></epp>`)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, res.Result.Code, 1000)
	st.Expect(t, res.Result.Message, "Command completed successfully")
	st.Expect(t, res.Result.IsError(), false)
	st.Expect(t, res.Result.IsFatal(), false)
}

func db(s string) *xml.Decoder {
	return xml.NewDecoder(bytes.NewBufferString(s))
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
		buf.WriteString(`<epp><response><result code="1000"><msg>Command completed successfully</msg></result></response></epp>`)
		d.Reset()
		b.StartTimer()
		d.decodeResult(&r)
	}
}

func BenchmarkScanResult(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		d := db(`<epp><response><result code="1000"><msg>Command completed successfully</msg></result></response></epp>`)
		b.StartTimer()
		var res response_
		scanResponse.Scan(d, &res)
	}
}
