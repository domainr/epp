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
	c.decoder.Reset()
	return decodeGreeting(&c.decoder, &c.Greeting)
}

func decodeGreeting(d *Decoder, g *Greeting) error {
	g.ServerName = ""
	g.Languages = nil
	g.Versions = nil
	g.Objects = nil
	g.Extensions = nil
	return d.DecodeWith(func(t xml.Token) error {
		switch node := t.(type) {
		// FIXME: re-optimize this with a special error type?
		// case xml.EndElement:
		// 	// Escape early (skip remaining XML)
		// 	if node.Name.Local == "svcMenu" &&
		// 		g.ServerName != "" {
		// 		return io.EOF
		// 	}

		case xml.CharData:
			switch d.Element(-1).Name.Local {
			case "svID":
				g.ServerName = string(node)
			case "version":
				g.Versions = append(g.Versions, string(node))
			case "lang":
				g.Languages = append(g.Languages, string(node))
			case "objURI":
				g.Objects = append(g.Objects, string(node))
			case "extURI":
				g.Extensions = append(g.Extensions, string(node))
			}
		}
		return nil
	})
}
