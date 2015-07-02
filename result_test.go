package epp

import (
	"bytes"
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
	err := decodeResult(&d, &r)
	st.Expect(t, err, nil)
	st.Expect(t, r.Code, 1000)
	st.Expect(t, r.Message, "Command completed successfully")
}
