package epp

import "encoding/xml"

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() error {
	err := c.writeDataUnit(xmlHello)
	if err != nil {
		return err
	}
	return c.readGreeting()
}

var xmlHello = []byte(xml.Header + startEPP + `<hello/>` + endEPP)

// Greeting is an EPP response that represents server status and capabilities.
// https://tools.ietf.org/html/rfc5730#section-2.4
type Greeting struct {
	ServerName string   `xml:"svID"`
	Versions   []string `xml:"svcMenu>version"`
	Languages  []string `xml:"svcMenu>lang"`
	Objects    []string `xml:"svcMenu>objURI"`
	Extensions []string `xml:"svcMenu>svcExtension>extURI,omitempty"`
}

var (
	defaultObjects = []string{
		"urn:ietf:params:xml:ns:domain-1.0",
		"urn:ietf:params:xml:ns:host-1.0",
		"urn:ietf:params:xml:ns:contact-1.0",
		"http://www.unitedtld.com/epp/finance-1.0",
	}
	defaultExtensions = []string{
		"urn:ietf:params:xml:ns:secDNS-1.1",
		"urn:ietf:params:xml:ns:rgp-1.0",
		"urn:ietf:params:xml:ns:launch-1.0",
		"urn:ietf:params:xml:ns:idn-1.0",
		"http://www.unitedtld.com/epp/charge-1.0",
	}
)

func (c *Conn) readGreeting() error {
	err := c.readDataUnit()
	if err != nil {
		return err
	}
	deleteBufferRange(&c.buf, []byte(`<dcp>`), []byte(`</dcp>`))
	var res response_
	err = IgnoreEOF(scanResponse.Scan(c.decoder, &res))
	if err != nil {
		return err
	}
	c.Greeting = res.Greeting
	return nil
}

func init() {
	scanResponse.MustHandleCharData("epp > greeting > svID", func(c *Context) error {
		res := c.Value.(*response_)
		res.Greeting.ServerName = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData("epp > greeting > svcMenu > version", func(c *Context) error {
		res := c.Value.(*response_)
		res.Greeting.Versions = append(res.Greeting.Versions, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp > greeting > svcMenu > lang", func(c *Context) error {
		res := c.Value.(*response_)
		res.Greeting.Languages = append(res.Greeting.Languages, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp > greeting > svcMenu > objURI", func(c *Context) error {
		res := c.Value.(*response_)
		res.Greeting.Objects = append(res.Greeting.Objects, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData("epp > greeting > svcMenu > svcExtension > extURI", func(c *Context) error {
		res := c.Value.(*response_)
		res.Greeting.Extensions = append(res.Greeting.Extensions, string(c.CharData))
		return nil
	})
}
