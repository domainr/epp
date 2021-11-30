package epp

import (
	"github.com/domainr/epp/internal/schema/raw"
	"github.com/nbio/xml"
)

// EPP represents an <epp> message envelope as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Hello represents a client <hello> message.
	Hello *Hello `xml:"hello"`

	// Greeting represents a server <greeting> message.
	Greeting *Greeting `xml:"greeting"`

	// Command represents a client <command> message.
	Command *Command `xml:"command"`
}

// MarshalXML implements xml.Marshaler, so a <hello/> message will be generated with a self-closing tag.
func (epp *EPP) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Space = "urn:ietf:params:xml:ns:epp-1.0"
	start.Name.Local = "epp"
	if epp.Hello != nil {
		return e.EncodeElement(&raw.XML{Value: "<hello/>"}, start)
	}
	type T EPP
	return e.EncodeElement((*T)(epp), start)
}
