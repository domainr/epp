package epp

import "github.com/domainr/epp/internal/schema/std"

// Command represents an EPP client <command> as defined in RFC 5730.
type Command struct {
	Login               *Login    `xml:"login"`
	Logout              std.Bool  `xml:"logout"`
	Check               *Check    `xml:"check"`
	Info                *Info     `xml:"info"`
	Poll                *Poll     `xml:"poll"`
	Create              *Create   `xml:"create"`
	Update              *Update   `xml:"update"`
	Delete              *Delete   `xml:"delete"`
	Renew               *Renew    `xml:"renew"`
	Transfer            *Transfer `xml:"transfer"`
	ClientTransactionID string    `xml:"clTRID,omitempty"`
}
