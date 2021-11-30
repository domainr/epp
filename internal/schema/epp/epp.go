package epp

import (
	"github.com/domainr/epp/internal/schema/date"
	"github.com/domainr/epp/internal/schema/domain"
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

// Hello represents an EPP client <hello> message as defined in RFC 5730.
type Hello struct{}

// Greeting represents an EPP server <greeting> message as defined in RFC 5730.
type Greeting struct {
	ServerName  string       `xml:"svID,omitempty"`
	ServerDate  *date.Time   `xml:"svDate"`
	ServiceMenu *ServiceMenu `xml:"svcMenu"`
}

// ServiceMenu represents an EPP <svcMenu> element as defined in RFC 5730.
type ServiceMenu struct {
	Versions         []string          `xml:"version"`
	Languages        []string          `xml:"lang"`
	Objects          []string          `xml:"objURI"`
	ServiceExtension *ServiceExtension `xml:"svcExtension"`
}

// ServiceExtension represents an EPP <svcExtension> element as defined in RFC 5730.
type ServiceExtension struct {
	Extensions []string `xml:"extURI"`
}

// Command represents an EPP client <command> message as defined in RFC 5730.
type Command struct {
	Check               *Check `xml:"check"`
	ClientTransactionID string `xml:"clTRID,omitempty"`
}

// Check represents an EPP <check> command as defined in RFC 5730.
type Check struct {
	DomainCheck *domain.Check
}
