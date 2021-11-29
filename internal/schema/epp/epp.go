package epp

import "github.com/domainr/epp/internal/schema/domain"

// EPP represents a single <epp> message as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Hello   *Hello
	Command *Command
}

// Hello represents a single EPP <hello> message as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type Hello struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`
}

// Command represents a single EPP <command> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type Command struct {
	XMLName             struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`
	Check               *Check
	ClientTransactionID string `xml:"urn:ietf:params:xml:ns:epp-1.0 clTRID,omitempty"`
}

// Check represents a single EPP <check> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type Check struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	DomainCheck *domain.Check
}
