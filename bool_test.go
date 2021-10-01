package epp

import (
	"testing"

	"github.com/nbio/xml"

	"github.com/nbio/st"
)

func TestBool(t *testing.T) {
	x := []byte(`<example><fred/><susan/></example>`)
	var y struct {
		XMLName struct{} `xml:"example"`
		Fred    Bool     `xml:"fred"`
		Jane    Bool     `xml:"jane"`
		Susan   Bool     `xml:"susan"`
	}

	err := xml.Unmarshal(x, &y)
	st.Expect(t, err, nil)
	st.Expect(t, y.Fred, True)
	st.Expect(t, y.Jane, False)
	st.Expect(t, y.Susan, True)
	z, err := xml.Marshal(&y)
	st.Expect(t, err, nil)
	st.Expect(t, string(z), `<example><fred></fred><susan></susan></example>`)
}

func TestBoolAttr(t *testing.T) {
	x := []byte(`<example fred="1" jane="FALSE"></example>`)
	var y struct {
		XMLName struct{} `xml:"example"`
		Fred    Bool     `xml:"fred,attr"`
		Jane    Bool     `xml:"jane,attr"`
		Susan   Bool     `xml:"susan,attr"`
	}

	err := xml.Unmarshal(x, &y)
	st.Expect(t, err, nil)
	st.Expect(t, y.Fred, True)
	st.Expect(t, y.Jane, False)
	st.Expect(t, y.Susan, False)
	z, err := xml.Marshal(&y)
	st.Expect(t, err, nil)
	st.Expect(t, string(z), `<example fred="1" jane="0" susan="0"></example>`)
}
