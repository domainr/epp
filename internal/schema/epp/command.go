package epp

import (
	"github.com/nbio/xml"
)

// Command represents an EPP client <command> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.5.
type Command struct {
	XMLName             struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`
	Command             command
	ClientTransactionID string `xml:"clTRID,omitempty"`
}

func (Command) eppBody() {}

// UnmarshalXML implements the xml.Unmarshaler interface.
// It maps known EPP commands to their corresponding Go type.
func (c *Command) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Command
	var v struct {
		Check    *Check    `xml:"check"`
		Create   *Create   `xml:"create"`
		Delete   *Delete   `xml:"delete"`
		Info     *Info     `xml:"info"`
		Login    *Login    `xml:"login"`
		Logout   *Logout   `xml:"logout"`
		Poll     *Poll     `xml:"poll"`
		Renew    *Renew    `xml:"renew"`
		Transfer *Transfer `xml:"transfer"`
		Update   *Update   `xml:"update"`
		*T
	}
	v.T = (*T)(c)
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	switch {
	case v.Check != nil:
		c.Command = v.Check
	case v.Create != nil:
		c.Command = v.Create
	case v.Delete != nil:
		c.Command = v.Delete
	case v.Info != nil:
		c.Command = v.Info
	case v.Login != nil:
		c.Command = v.Login
	case v.Logout != nil:
		c.Command = v.Logout
	case v.Poll != nil:
		c.Command = v.Poll
	case v.Renew != nil:
		c.Command = v.Renew
	case v.Transfer != nil:
		c.Command = v.Transfer
	case v.Update != nil:
		c.Command = v.Update
	}
	return nil
}

// command is a child element of EPP <command>.
// Concrete command types implement this interface.
type command interface {
	eppCommand()
}
