package epp

// EPP represents an <epp> message envelope as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Hello represents a client <hello> message.
	Hello *Hello `xml:"hello,selfclosing"`

	// Greeting represents a server <greeting> message.
	Greeting *Greeting `xml:"greeting"`

	// Command represents a client <command> message.
	Command *Command `xml:"command"`
}
