package epp

// Msg represents EPP request and response messages as defined by RFC 5730.
// https://tools.ietf.org/html/rfc5730#section-2.2
type Msg struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Request elements. Leave nil if unused.
	Hello   *Hello   `xml:"hello"`
	Command *Command `xml:"command"`

	// Response elements. Will be nil if not present in response message.
	Greeting *Greeting `xml:"greeting"`
	Response *Response `xml:"response"`
}
