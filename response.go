package epp

import "github.com/nbio/xx"

// Response represents an EPP response.
type Response struct {
	Result
	Greeting
	DomainCheckResponse
	DomainInfoResponse
}

var scanResponse = xx.NewScanner()

func init() {
	scanResponse.MustHandleStartElement("epp", func(c *xx.Context) error {
		*c.Value.(*Response) = Response{}
		return nil
	})
}
