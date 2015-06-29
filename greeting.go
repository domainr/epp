package epp

import "errors"

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() (*Greeting, error) {
	err := c.writeDataUnit(xmlHello)
	if err != nil {
		return nil, err
	}
	return c.ReadGreeting()
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

// ReadGreeting reads a <greeting> message from the EPP server.
// Performed automatically during a Handshake or Hello command.
func (c *Conn) ReadGreeting() (*Greeting, error) {
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
