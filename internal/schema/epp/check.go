package epp

import "github.com/domainr/epp/internal/schema/domain"

// Check represents an EPP <check> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.1.
type Check struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	DomainCheck *domain.Check
	// TODO: HostCheck *domain.Check
}

func (Check) eppCommand() {}
