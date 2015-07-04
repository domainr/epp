package epp

import (
	"encoding/xml"
	"fmt"
	"io"
	"testing"

	"github.com/nbio/st"
)

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

	d := db(x)
	var res response_
	err := s.Scan(d, &res)
	st.Expect(t, err, io.EOF)
}

func TestScannerInvalidXML(t *testing.T) {
	s := NewScanner()
	d := db(`<foo><bar/><baz/`)
	err := s.Scan(d, nil)
	_, ok := err.(*xml.SyntaxError)
	st.Expect(t, ok, true)
}

func debugStartElement(ctx *Context) error {
	fmt.Printf("xml.StartElement: <%s xmlns=\"%s\">\n", ctx.StartElement.Name.Local, ctx.StartElement.Name.Space)
	return nil
}

func debugCharData(ctx *Context) error {
	fmt.Printf("xml.CharData: %s\n", string(ctx.CharData))
	return nil
}
