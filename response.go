package epp

import "github.com/nbio/xx"

type response_ struct {
	Result
	Greeting
	DomainCheckResponse
	DomainInfoResponse
}

var scanResponse = xx.NewScanner()

func init() {
	scanResponse.MustHandleStartElement("epp", func(c *xx.Context) error {
		*c.Value.(*response_) = response_{}
		return nil
	})
}
