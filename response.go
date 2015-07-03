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
	scanResponse.MustHandle("epp>greeting", func(c *Context) error {
		if c.CharData == nil {
			c.Value.(*response_).greeting = &Greeting{}
		}
		return nil
	})
}
