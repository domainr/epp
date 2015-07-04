package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"testing"

	"github.com/nbio/st"
)

func decoder(s string) *xml.Decoder {
	return xml.NewDecoder(bytes.NewBufferString(s))
}

func debugStartElement(ctx *Context) error {
	fmt.Printf("xml.StartElement: <%s xmlns=\"%s\">\n", ctx.StartElement.Name.Local, ctx.StartElement.Name.Space)
	return nil
}

func debugCharData(ctx *Context) error {
	fmt.Printf("xml.CharData: %s\n", string(ctx.CharData))
	return nil
}

func TestScanner(t *testing.T) {
	s := NewScanner()
	s.MustHandleStartElement("epp", func(ctx *Context) error { return nil })
	s.MustHandleStartElement("epp>response>result", func(ctx *Context) error { return nil })
	s.MustHandleCharData("epp>response>result>msg", func(ctx *Context) error { return nil })

	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1000">
      <msg>Command completed successfully</msg>
    </result>
    <trID>
      <clTRID>0000000000000001</clTRID>
      <svTRID>3577a51b-5a4b-4025-8556-0a3e392b4097:1</svTRID>
    </trID>
  </response>
</epp>`

	d := decoder(x)
	var res response_
	err := s.Scan(d, &res)
	st.Expect(t, err, io.EOF)
}

func TestScannerInvalidXML(t *testing.T) {
	s := NewScanner()
	d := decoder(`<foo><bar/><baz/`)
	err := s.Scan(d, nil)
	_, ok := err.(*xml.SyntaxError)
	st.Expect(t, ok, true)
}

func TestContextAttr(t *testing.T) {
	ctx := Context{
		StartElement: xml.StartElement{
			Name: xml.Name{"urn:ietf:params:xml:ns:epp-1.0", "example"},
			Attr: []xml.Attr{
				{xml.Name{"urn:ietf:params:xml:ns:epp-1.0", "avail"}, "1"},
				{xml.Name{"", "alpha"}, "TRUE"},
				{xml.Name{"", "beta"}, "1"},
				{xml.Name{"other", "gamma"}, "false"},
				{xml.Name{"other", "delta"}, "FALSE"},
				{xml.Name{"", "omega"}, "hammertime"},
				{xml.Name{"", "number"}, "42"},
			},
		},
	}

	st.Expect(t, ctx.Attr("", "avail"), "1")
	st.Expect(t, ctx.AttrBool("", "avail"), true)
	st.Expect(t, ctx.AttrInt("", "avail"), 1)
	st.Expect(t, ctx.Attr("other", "avail"), "")
	st.Expect(t, ctx.AttrBool("other", "avail"), false)
	st.Expect(t, ctx.AttrInt("other", "avail"), 0)
	st.Expect(t, ctx.AttrBool("", "alpha"), true)
	st.Expect(t, ctx.AttrBool("", "beta"), true)
	st.Expect(t, ctx.AttrBool("", "gamma"), false)
	st.Expect(t, ctx.AttrBool("", "delta"), false)
	st.Expect(t, ctx.Attr("other", "omega"), "")
	st.Expect(t, ctx.AttrBool("other", "omega"), false)
	st.Expect(t, ctx.Attr("", "omega"), "hammertime")
	st.Expect(t, ctx.AttrBool("", "omega"), true)
	st.Expect(t, ctx.AttrInt("", "number"), 42)
}
