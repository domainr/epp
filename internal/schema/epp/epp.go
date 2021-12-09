package epp

import (
	"github.com/nbio/xml"
)

// EPP represents an <epp> element as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
//
// Body is a pointer to a valid EPP child element:
// <hello>, <greeting>, <command>, and <response>.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Body    body
}

func (e *EPP) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Hello    *Hello    `xml:"hello"`
		Greeting *Greeting `xml:"greeting"`
		Command  *Command  `xml:"command"`
		Response *Response `xml:"response"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	switch {
	case v.Hello != nil:
		e.Body = v.Hello
	case v.Greeting != nil:
		e.Body = v.Greeting
	case v.Command != nil:
		e.Body = v.Command
	case v.Response != nil:
		e.Body = v.Response
	}
	return nil
}

type body interface {
	eppBody()
}
