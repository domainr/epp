package epp

import (
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
