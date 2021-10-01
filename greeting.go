package epp

import (
	"github.com/nbio/xml"

	"github.com/nbio/xx"
)

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() error {
	err := c.writeRequest(xmlHello)
	if err != nil {
		return err
	}
	_, err = c.readGreeting()
	return err
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

// SupportsObject returns true if the EPP server supports
// the object specified by uri.
func (g *Greeting) SupportsObject(uri string) bool {
	if g == nil {
		return false
	}
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
	if g == nil {
		return false
	}
	for _, v := range g.Extensions {
		if v == uri {
			return true
		}
	}
	return false
}

func (c *Conn) readGreeting() (Greeting, error) {
	res, err := c.readResponse()
	if err != nil {
		return Greeting{}, err
	}
	return res.Greeting, nil
}

func init() {
	path := "epp>greeting"
	scanResponse.MustHandleCharData(path+">svID", func(c *xx.Context) error {
		res := c.Value.(*Response)
		res.Greeting.ServerName = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>version", func(c *xx.Context) error {
		res := c.Value.(*Response)
		res.Greeting.Versions = append(res.Greeting.Versions, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>lang", func(c *xx.Context) error {
		res := c.Value.(*Response)
		res.Greeting.Languages = append(res.Greeting.Languages, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>objURI", func(c *xx.Context) error {
		res := c.Value.(*Response)
		res.Greeting.Objects = append(res.Greeting.Objects, string(c.CharData))
		return nil
	})
	scanResponse.MustHandleCharData(path+">svcMenu>svcExtension>extURI", func(c *xx.Context) error {
		res := c.Value.(*Response)
		res.Greeting.Extensions = append(res.Greeting.Extensions, string(c.CharData))
		return nil
	})
}
