package epp

import (
	"bytes"
	"testing"

	"github.com/nbio/st"
)

func TestScanResult(t *testing.T) {
	var res response_
	r := &res.Result

	d := decoder(`<epp><response><result code="1000"><msg>Command completed successfully</msg></result></response></epp>`)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, r.Code, 1000)
	st.Expect(t, r.Message, "Command completed successfully")
	st.Expect(t, r.IsError(), false)
	st.Expect(t, r.IsFatal(), false)

	// Result code >= 2000 is an error.
	d = decoder(`<epp><response><result code="2001"><msg>Command syntax error</msg></result></response></epp>`)
	err = IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, r.Code, 2001)
	st.Expect(t, r.Message, "Command syntax error")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), false)

	// Result code > 2500 is a fatal error.
	d = decoder(`<epp><response><result code="2501"><msg>Authentication error; server closing connection</msg></result></response></epp>`)
	err = IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, r.Code, 2501)
	st.Expect(t, r.Message, "Authentication error; server closing connection")
	st.Expect(t, r.IsError(), true)
	st.Expect(t, r.IsFatal(), true)
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
		d := decoder(`<epp><response><result code="1000"><msg>Command completed successfully</msg></result></response></epp>`)
		b.StartTimer()
		var res response_
		scanResponse.Scan(d, &res)
	}
}
