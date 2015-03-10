package epp

import "errors"

// Hello represents a client <hello> (request for <greeting>).
// https://tools.ietf.org/html/rfc5730#section-2.3
type Hello struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Hello   struct{} `xml:"hello"`
}

// <epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello/></epp>
var hello = Hello{}

// Hello sends a <hello> command to request a <greeting> from the EPP server.
func (c *Conn) Hello() (err error) {
	err = c.WriteMessage(&hello)
	if err != nil {
		return
	}
	return c.readGreeting()
}

// Greeting is an EPP response that represents server status and capabilities.
// https://tools.ietf.org/html/rfc5730#section-2.4
type Greeting struct {
	ServerName  string `xml:"svID"`
	ServerTime  Time   `xml:"svDate"`
	ServiceMenu struct {
		Versions   []string `xml:"version"`
		Languages  []string `xml:"lang"`
		Objects    []string `xml:"objURI"`
		Extensions []string `xml:"svcExtension>extURI"`
	} `xml:"svcMenu"`
	DCP struct {
		Access struct {
			All              *struct{} `xml:"all"`
			None             *struct{} `xml:"none"`
			Null             *struct{} `xml:"null"`
			Personal         *struct{} `xml:"personal"`
			PersonalAndOther *struct{} `xml:"personalAndOther"`
			Other            *struct{} `xml:"other"`
		} `xml:"access"`
		Statement []struct {
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
		} `xml:"statement"`
	} `xml:"dcp"`
}

// A ErrMissingGreeting is reported when a <greeting> message is expected but not found.
var ErrMissingGreeting = errors.New("expected <greeting> message in EPP message, but none found")

// readGreeting reads a <greeting> message from the EPP server.
// It stores the last-read <greeting> message on the connection,
func (c *Conn) readGreeting() (err error) {
	var rmsg Response
	err = c.ReadResponse(&rmsg)
	if err != nil {
		return
	}
	if rmsg.Greeting == nil {
		return ErrMissingGreeting
	}
	c.Greeting = rmsg.Greeting
	return
}
