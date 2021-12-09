package epp

import (
	"github.com/domainr/epp/internal/schema/domain"
	"github.com/nbio/xml"
)

// Check represents an EPP <check> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.1.
type Check struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	Check   check
}

func (Check) eppCommand() {}

// UnmarshalXML implements the xml.Unmarshaler interface.
// It maps known EPP check commands to their corresponding Go type.
func (c *Check) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Check
	var v struct {
		DomainCheck *domain.Check
		// TODO: HostCheck, etc.
		*T
	}
	v.T = (*T)(c)
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	switch {
	case v.DomainCheck != nil:
		c.Check = v.DomainCheck
	}
	return nil
}

type check interface {
	EPPCheck()
}
