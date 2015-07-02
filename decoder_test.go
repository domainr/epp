package epp

import (
	"bytes"
	"encoding/xml"
	"io"
	"testing"

	"github.com/nbio/st"
)

func logMarshal(t *testing.T, msg *message) {
	x, err := xml.Marshal(&msg)
	st.Expect(t, err, nil)
	t.Logf("<!-- MARSHALED -->\n%s\n", string(x))
}

func TestDecoderReuse(t *testing.T) {
	buf := bytes.Buffer{}
	d := NewDecoder(&buf)

	v := struct {
		XMLName struct{} `xml:"hello"`
		Foo     string   `xml:"foo"`
	}{}

	buf.Reset()
	buf.Write([]byte(`<hello><foo>foo</foo></hello>`))
	d.Reset()
	st.Expect(t, d.decoder.InputOffset(), int64(0))
	d.decoder.Decode(&v)
	st.Expect(t, v.Foo, "foo")
	st.Expect(t, d.decoder.InputOffset(), int64(29))

	buf.Reset()
	buf.Write([]byte(`<hello><foo>bar</foo></hello>`))
	d.Reset()
	st.Expect(t, d.decoder.InputOffset(), int64(0))
	tok, _ := d.Token()
	se := tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "hello")
	tok, _ = d.Token()
	se = tok.(xml.StartElement)
	st.Expect(t, se.Name.Local, "foo")
	st.Expect(t, d.decoder.InputOffset(), int64(12))

	buf.Reset()
	buf.Write([]byte(`<hello><foo>blam&lt;</foo></hello>`))
	d.Reset()
	st.Expect(t, d.decoder.InputOffset(), int64(0))
	d.decoder.Decode(&v)
	st.Expect(t, v.Foo, "blam<")
	st.Expect(t, d.decoder.InputOffset(), int64(34))
}

func TestUnmarshalCheckDomainResponse(t *testing.T) {
	x := []byte(`<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">good.memorial</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<charge:chkData xmlns:charge="http://www.unitedtld.com/epp/charge-1.0">
				<charge:cd>
					<charge:name>good.memorial</charge:name>
					<charge:set>
						<charge:category name="BBB+">premium</charge:category>
						<charge:type>price</charge:type>
						<charge:amount command="create">100.00</charge:amount>
						<charge:amount command="renew">100.00</charge:amount>
						<charge:amount command="transfer">100.00</charge:amount>
						<charge:amount command="update" name="restore">50.00</charge:amount>
					</charge:set>
				</charge:cd>
			</charge:chkData>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`)

	d := NewDecoder(bytes.NewBuffer(x))
	var msg message
	err := d.DecodeMessage(&msg)
	st.Expect(t, err, nil)
	st.Reject(t, msg.Response, nil)
	st.Expect(t, msg.Response.ResponseData.DomainCheckData.Results[0].Domain.Domain, "good.memorial")
	st.Expect(t, msg.Response.ResponseData.DomainCheckData.Results[0].Domain.IsAvailable, true)
	logMarshal(t, &msg)
}

func TestDecoderAtPath(t *testing.T) {
	x := []byte(`<foo><bar><baz></baz></bar></foo>`)
	d := NewDecoder(bytes.NewBuffer(x))
	d.Token()
	st.Expect(t, d.AtPath("blam"), false)
	st.Expect(t, d.AtPath("foo"), true)
	st.Expect(t, d.AtPath("foo", "bar"), false)
	d.Token()
	st.Expect(t, len(d.Stack), 2)
	st.Expect(t, d.Element(-1).Name.Local, "bar")
	st.Expect(t, d.AtPath("foo", "bar"), true)
	st.Expect(t, d.AtPath("foo", "bar", "baz"), false)
	d.Token()
	st.Expect(t, len(d.Stack), 3)
	st.Expect(t, d.Element(-1).Name.Local, "baz")
	st.Expect(t, d.AtPath("foo", "bar"), false)
	st.Expect(t, d.AtPath("foo", "bar", "baz"), true)
	d.Token()
	st.Expect(t, len(d.Stack), 2)
	st.Expect(t, d.Element(-1).Name.Local, "bar")
	st.Expect(t, d.AtPath("foo", "bar"), true)
	st.Expect(t, d.AtPath("foo", "bar", "baz"), false)
	d.Token()
	st.Expect(t, len(d.Stack), 1)
	st.Expect(t, d.Element(-1).Name.Local, "foo")
	st.Expect(t, d.AtPath("foo"), true)
	st.Expect(t, d.AtPath("foo", "bar"), false)
	d.Token()
	st.Expect(t, len(d.Stack), 0)
	st.Expect(t, d.Element(-1).Name.Local, "")
	st.Expect(t, d.AtPath(), true)
	_, err := d.Token()
	st.Expect(t, err, io.EOF)
}
