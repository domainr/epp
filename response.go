package epp

type response_ struct {
	result   *Result
	greeting *Greeting
	domains  []domainCheckData_
	// charges  []chargeCheckData
}

var scanResponse *Scanner

func init() {
	scanResponse = NewScanner()
	scanResponse.MustHandleStartElement("epp>greeting", func(c *Context) error {
		c.Value.(*response_).greeting = &Greeting{}
		return nil
	})
}
