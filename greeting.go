package epp

import (
	"encoding/xml"
	"io"
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
	g.Languages = g.Languages[:0]
	g.Versions = g.Versions[:0]
	g.Objects = g.Objects[:0]
	g.Extensions = g.Extensions[:0]
	return d.DecodeWith(func(t xml.Token) error {
		switch node := t.(type) {
		case xml.EndElement:
			// Escape early (skip remaining XML)
			if node.Name.Local == "svcMenu" &&
				g.ServerName != "" {
				return io.EOF
			}

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
