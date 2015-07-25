package epp

import (
	"encoding/xml"

	"github.com/nbio/xx"
)

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

// SupportsExtension returns true if the EPP server supports
// the object specified by uri.
func (g *Greeting) SupportsObject(uri string) bool {
	for _, v := range g.Objects {
		if v == uri {
			return true
		}
	}
	return false
}

// SupportsExtension returns true if the EPP server supports
// the extension specified by uri.
func (g *Greeting) SupportsExtension(uri string) bool {
	for _, v := range g.Extensions {
		if v == uri {
			return true
		}
	}
	return false
}

const (
	ObjDomain  = "urn:ietf:params:xml:ns:domain-1.0"
	ObjHost    = "urn:ietf:params:xml:ns:host-1.0"
	ObjContact = "urn:ietf:params:xml:ns:contact-1.0"
	ObjFinance = "http://www.unitedtld.com/epp/finance-1.0"
	ExtSecDNS  = "urn:ietf:params:xml:ns:secDNS-1.1"
	ExtRGP     = "urn:ietf:params:xml:ns:rgp-1.0"
	ExtLaunch  = "urn:ietf:params:xml:ns:launch-1.0"
	ExtIDN     = "urn:ietf:params:xml:ns:idn-1.0"
	ExtCharge  = "http://www.unitedtld.com/epp/charge-1.0"
	ExtFee     = "urn:ietf:params:xml:ns:fee-0.5"
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
	path := "epp>greeting"
	scanResponse.MustHandleCharData(path+">svID", func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Greeting.ServerName = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>version", func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Greeting.Versions = append(res.Greeting.Versions, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>lang", func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Greeting.Languages = append(res.Greeting.Languages, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>objURI", func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Greeting.Objects = append(res.Greeting.Objects, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>svcExtension>extURI", func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Greeting.Extensions = append(res.Greeting.Extensions, string(c.CharData))
		return nil
	})
}
