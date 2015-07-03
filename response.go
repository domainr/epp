package epp

type response_ struct {
	result   Result
	greeting Greeting
	domains  []domainCheckData_
	// charges  []chargeCheckData
}

var scanResponse = NewScanner()

func init() {
	scanResponse.MustHandleStartElement("epp>greeting", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting = Greeting{}
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svID", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting.ServerName = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>version", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting.Versions = append(res.greeting.Versions, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>lang", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting.Languages = append(res.greeting.Languages, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>objURI", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting.Objects = append(res.greeting.Objects, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>svcExtension>extURI", func(c *Context) error {
		res := c.Value.(*response_)
		res.greeting.Extensions = append(res.greeting.Extensions, string(c.CharData))
		return nil
	})
}
