package epp

import (
	"encoding/xml"
	"testing"

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
