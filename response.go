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
		g := c.Value.(*response_).greeting
		g.ServerName = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>version", func(c *Context) error {
		g := c.Value.(*response_).greeting
		g.Versions = append(g.Versions, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>lang", func(c *Context) error {
		g := c.Value.(*response_).greeting
		g.Languages = append(g.Languages, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>objURI", func(c *Context) error {
		g := c.Value.(*response_).greeting
		g.Objects = append(g.Objects, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp>greeting>svcMenu>svcExtension>extURI", func(c *Context) error {
		g := c.Value.(*response_).greeting
		g.Extensions = append(g.Extensions, string(c.CharData))
		return nil
	})
}
