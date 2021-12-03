package epp

import "github.com/domainr/epp/internal/schema/std"

// EPP represents an <epp> element as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Hello represents a client <hello> element.
	Hello std.Bool `xml:"hello,selfclosing"`

	// Greeting represents a server <greeting> element.
	Greeting *Greeting `xml:"greeting"`

	// Command represents a client <command> element.
	Command *Command `xml:"command"`

	// Response represents a server <response> element.
	Response *Response `xml:"response"`
}
