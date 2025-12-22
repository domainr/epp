package epp

import "github.com/nbio/xx"

// Response represents an EPP response.
// Response represents an EPP response.
type Response struct {
	Result
	Greeting
	DomainCheckResponse
	DomainInfoResponse
	// Additions for new commands
	DomainCreateResponse
	DomainRenewResponse
	DomainTransferResponse
	DomainUpdateResponse
	ContactCreateResponse
	ContactInfoResponse
}

var scanResponse = xx.NewScanner()

func init() {
	scanResponse.MustHandleStartElement("epp", func(c *xx.Context) error {
		*c.Value.(*Response) = Response{}
		return nil
	})
}
