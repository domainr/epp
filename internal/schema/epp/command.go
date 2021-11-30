package epp

import (
	"github.com/domainr/epp/internal/schema/domain"
)

// Command represents an EPP client <command> message as defined in RFC 5730.
type Command struct {
	Check               *Check `xml:"check"`
	ClientTransactionID string `xml:"clTRID,omitempty"`
}

// Check represents an EPP <check> command as defined in RFC 5730.
type Check struct {
	DomainCheck *domain.Check
}
