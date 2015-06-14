package epp

import "errors"

var msgHello = message{Hello: &hello{}}

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() (*Greeting, error) {
	err := c.writeMessage(&msgHello)
	if err != nil {
		return nil, err
	}
	return c.ReadGreeting()
}

// Greeting is an EPP response that represents server status and capabilities.
// https://tools.ietf.org/html/rfc5730#section-2.4
type Greeting struct {
	ServerName        string         `xml:"svID"`
	ServerTime        Time           `xml:"svDate"`
	ServiceVersions   []string       `xml:"svcMenu>version"`
	ServiceLanguages  []string       `xml:"svcMenu>lang"`
	ServiceObjects    []string       `xml:"svcMenu>objURI"`
	ServiceExtensions []string       `xml:"svcMenu>svcExtension>extURI,omitempty"`
	DCPAccess         DCPAccess      `xml:"dcp>access"`
	DCPStatements     []DCPStatement `xml:"dcp>statement,omitempty"`
}

// DCPAccess defines the data collection policy for an EPP server.
type DCPAccess struct {
	All              Bool `xml:"all"`
	None             Bool `xml:"none"`
	Null             Bool `xml:"null"`
	Personal         Bool `xml:"personal"`
	PersonalAndOther Bool `xml:"personalAndOther"`
	Other            Bool `xml:"other"`
}

// DCPStatement defines the data collection policy
// for an EPP server. The EPP server may return 1 or more
// DCPStatements in a Greeting message.
// https://tools.ietf.org/html/rfc5730#section-2.4
type DCPStatement struct {
	Purpose struct {
		Admin        Bool `xml:"admin"`
		Contact      Bool `xml:"contact"`
		Provisioning Bool `xml:"prov"`
		Other        Bool `xml:"other"`
	} `xml:"purpose"`
	Recipient struct {
		Other Bool `xml:"other"`
		Ours  *struct {
			Recipient string `xml:"recDesc"`
		} `xml:"ours"`
		Public    Bool `xml:"public"`
		Same      Bool `xml:"same"`
		Unrelated Bool `xml:"unrelated"`
	} `xml:"recipient"`
	Retention struct {
		Business   Bool `xml:"business"`
		Indefinite Bool `xml:"indefinite"`
		Legal      Bool `xml:"legal"`
		None       Bool `xml:"none"`
		Stated     Bool `xml:"stated"`
	} `xml:"retention"`
	Expiry *struct {
		Absolute Bool `xml:"absolute"`
		Relative Bool `xml:"relative"`
	} `xml:"expiry"`
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
