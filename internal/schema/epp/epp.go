package epp

import (
	"github.com/domainr/epp/internal/schema/date"
	"github.com/domainr/epp/internal/schema/domain"
)

// EPP represents a single <epp> message as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Hello represents a client <hello> message.
	Hello *Hello

	// Greeting represents a server <greeting> message.
	Greeting *Greeting

	// Command represents a client <command> message.
	Command *Command
}

// Hello represents an EPP client <hello> message as defined in RFC 5730.
type Hello struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 hello"`
}

// Greeting represents an EPP server <greeting> message as defined in RFC 5730.
type Greeting struct {
	XMLName     struct{}   `xml:"urn:ietf:params:xml:ns:epp-1.0 greeting"`
	ServerName  string     `xml:"urn:ietf:params:xml:ns:epp-1.0 svID,omitempty"`
	ServerDate  *date.Time `xml:"urn:ietf:params:xml:ns:epp-1.0 svDate,omitempty"`
	ServiceMenu *ServiceMenu
}

// ServiceMenu represents an EPP <svcMenu> element as defined in RFC 5730.
type ServiceMenu struct {
	XMLName    struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 svcMenu"`
	Versions   []string `xml:"urn:ietf:params:xml:ns:epp-1.0 version,omitempty"`
	Languages  []string `xml:"urn:ietf:params:xml:ns:epp-1.0 lang,omitempty"`
	Objects    []string `xml:"urn:ietf:params:xml:ns:epp-1.0 objURI,omitempty"`
	Extensions []string `xml:"urn:ietf:params:xml:ns:epp-1.0 svcExtension>extURI,omitempty"`
}

// Command represents an EPP client <command> message as defined in RFC 5730.
type Command struct {
	XMLName             struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`
	Check               *Check
	ClientTransactionID string `xml:"urn:ietf:params:xml:ns:epp-1.0 clTRID,omitempty"`
}

// Check represents an EPP <check> command as defined in RFC 5730.
type Check struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	DomainCheck *domain.Check
}