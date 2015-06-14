package epp

import "errors"

var msgHello = Message{Hello: &hello{}}

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() (*Greeting, error) {
	err := c.WriteMessage(&msgHello)
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
	All              *struct{} `xml:"all"`
	None             *struct{} `xml:"none"`
	Null             *struct{} `xml:"null"`
	Personal         *struct{} `xml:"personal"`
	PersonalAndOther *struct{} `xml:"personalAndOther"`
	Other            *struct{} `xml:"other"`
}

// DCPStatement defines the data collection policy
// for an EPP server. The EPP server may return 1 or more
// DCPStatements in a Greeting message.
// https://tools.ietf.org/html/rfc5730#section-2.4
type DCPStatement struct {
	Purpose struct {
		Admin        *struct{} `xml:"admin"`
		Contact      *struct{} `xml:"contact"`
		Provisioning *struct{} `xml:"prov"`
		Other        *struct{} `xml:"other"`
	} `xml:"purpose"`
	Recipient struct {
		Other *struct{} `xml:"other"`
		Ours  *struct {
			Recipient string `xml:"recDesc"`
		} `xml:"ours"`
		Public    *struct{} `xml:"public"`
		Same      *struct{} `xml:"same"`
		Unrelated *struct{} `xml:"unrelated"`
	} `xml:"recipient"`
	Retention struct {
		Business   *struct{} `xml:"business"`
		Indefinite *struct{} `xml:"indefinite"`
		Legal      *struct{} `xml:"legal"`
		None       *struct{} `xml:"none"`
		Stated     *struct{} `xml:"stated"`
	} `xml:"retention"`
	Expiry *struct {
		Absolute *struct{} `xml:"absolute"`
		Relative *struct{} `xml:"relative"`
	} `xml:"expiry"`
}

// A ErrGreetingNotFound is reported when a <greeting> message is expected but not found.
var ErrGreetingNotFound = errors.New("missing <greeting> message")

// ReadGreeting reads a <greeting> message from the EPP server.
// Performed automatically during a Handshake or Hello command.
func (c *Conn) ReadGreeting() (*Greeting, error) {
	var msg Message
	err := c.ReadMessage(&msg)
	if err != nil {
		return nil, err
	}
	if msg.Greeting == nil {
		return nil, ErrGreetingNotFound
	}
	return msg.Greeting, nil
}
