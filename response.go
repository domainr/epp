package epp

type response_ struct {
	Result
	Greeting
	DomainCheckResponse
}

var scanResponse = NewScanner()

func init() {
	scanResponse.MustHandleStartElement("epp", func(c *Context) error {
		*c.Value.(*response_) = response_{}
		return nil
	})
}
