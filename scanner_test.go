package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"testing"

	"github.com/nbio/st"
)

func TestScanner(t *testing.T) {
	s := NewScanner()
	err := s.Handle("epp", debugScanFunc)
	st.Expect(t, err, nil)
	err = s.Handle("epp>response>result", debugScanFunc)
	st.Expect(t, err, nil)

	x := []byte(`<?xml version="1.0" encoding="utf-8"?>
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
</epp>`)

	d := xml.NewDecoder(bytes.NewBuffer(x))
	var res response_
	err = s.Scan(d, &res)
	st.Expect(t, err, io.EOF)
}

func debugScanFunc(ctx *Context) error {
	if ctx.CharData != nil {
		fmt.Printf("xml.CharData: %s\n", string(ctx.CharData))
	} else {
		fmt.Printf("xml.StartElement: <%s xmlns=\"%s\">\n", ctx.StartElement.Name.Local, ctx.StartElement.Name.Space)
	}
	return nil
}
