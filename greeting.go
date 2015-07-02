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
	return decodeGreeting(&c.decoder, &c.Greeting)
}

func decodeGreeting(d *Decoder, g *Greeting) error {
	d.Reset()
	g.ServerName = ""
	g.Languages = g.Languages[:0]
	g.Versions = g.Versions[:0]
	g.Objects = g.Objects[:0]
	g.Extensions = g.Extensions[:0]
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch node := t.(type) {
		case xml.StartElement:
			// Ignore <dcp> section entirely
			if node.Name.Local == "dcp" {
				err = d.Skip()
				if err != nil {
					return err
				}
			}

		case xml.EndElement:
			// Escape early (skip remaining XML)
			if node.Name.Local == "svcMenu" &&
				g.ServerName != "" {
				return nil
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
	}
	return nil
}
