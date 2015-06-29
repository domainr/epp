package epp

import (
	"encoding/xml"
	"errors"
	"io"
)

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() (*Greeting, error) {
	err := c.writeDataUnit(xmlHello)
	if err != nil {
		return nil, err
	}
	return c.readGreeting()
}

// Greeting is an EPP response that represents server status and capabilities.
// https://tools.ietf.org/html/rfc5730#section-2.4
type Greeting struct {
	ServerName        string   `xml:"svID"`
	ServerTime        Time     `xml:"svDate"`
	ServiceVersions   []string `xml:"svcMenu>version"`
	ServiceLanguages  []string `xml:"svcMenu>lang"`
	ServiceObjects    []string `xml:"svcMenu>objURI"`
	ServiceExtensions []string `xml:"svcMenu>svcExtension>extURI,omitempty"`
}

// A ErrGreetingNotFound is reported when a <greeting> message is expected but not found.
var ErrGreetingNotFound = errors.New("missing <greeting> message")

// readGreeting reads a <greeting> message from the EPP server.
// Performed automatically during a Handshake or Hello command.
func (c *Conn) readGreeting() (*Greeting, error) {
	var msg message
	err := c.readMessage(&msg)
	if err != nil {
		return nil, err
	}
	if msg.Greeting == nil {
		return nil, ErrGreetingNotFound
	}
	return msg.Greeting, nil
}

func (d *Decoder) decodeGreeting(g *Greeting) error {
	d.Reset()
	g.ServiceLanguages = g.ServiceLanguages[:0]
	g.ServiceVersions = g.ServiceVersions[:0]
	g.ServiceObjects = g.ServiceObjects[:0]
	g.ServiceExtensions = g.ServiceExtensions[:0]
	for {
		t, err := d.Token()
		if err != nil && err != io.EOF {
			return err
		}
		if t == nil {
			break
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
			var t Time
			if node.Name.Local == "svcMenu" &&
				g.ServerName != "" &&
				g.ServerTime != t {
				return nil
			}

		// Extract character data
		case xml.CharData:
			if len(d.Stack) == 0 {
				continue
			}
			switch d.Stack[len(d.Stack)-1].Name.Local {
			case "svID":
				g.ServerName = string(node)
			case "svDate":
				g.ServerTime.UnmarshalText(node)
			case "version":
				g.ServiceVersions = append(g.ServiceVersions, string(node))
			case "lang":
				g.ServiceLanguages = append(g.ServiceLanguages, string(node))
			case "objURI":
				g.ServiceObjects = append(g.ServiceObjects, string(node))
			case "extURI":
				g.ServiceExtensions = append(g.ServiceExtensions, string(node))
			}
		}
	}
	return nil
}
