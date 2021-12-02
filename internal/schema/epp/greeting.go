package epp

import (
	"github.com/domainr/epp/internal/schema/std"
	"github.com/nbio/xml"
)

// Greeting represents an EPP server <greeting> message as defined in RFC 5730.
type Greeting struct {
	ServerName  string       `xml:"svID,omitempty"`
	ServerDate  *std.Time    `xml:"svDate"`
	ServiceMenu *ServiceMenu `xml:"svcMenu"`
	DCP         *DCP         `xml:"dcp"`
}

// ServiceMenu represents an EPP <svcMenu> element as defined in RFC 5730.
type ServiceMenu struct {
	Versions         []string          `xml:"version"`
	Languages        []string          `xml:"lang"`
	Objects          []string          `xml:"objURI"`
	ServiceExtension *ServiceExtension `xml:"svcExtension"`
}

// ServiceExtension represents an EPP <svcExtension> element as defined in RFC 5730.
type ServiceExtension struct {
	Extensions []string `xml:"extURI"`
}

// DCP represents a server data collection policy as defined in RFC 5730.
type DCP struct {
	Access     Access      `xml:"access"`
	Statements []Statement `xml:"statement"`
	Expiry     *Expiry     `xml:"expiry"`
}

// Access represents an EPP server’s scope of data access as defined in RFC 5730.
type Access struct {
	Null             std.Bool `xml:"null"`
	All              std.Bool `xml:"all"`
	None             std.Bool `xml:"none"`
	Other            std.Bool `xml:"other"`
	Personal         std.Bool `xml:"personal"`
	PersonalAndOther std.Bool `xml:"personalAndOther"`
}

var (
	AccessNull             = Access{Null: std.True}
	AccessAll              = Access{All: std.True}
	AccessNone             = Access{None: std.True}
	AccessOther            = Access{Other: std.True}
	AccessPersonal         = Access{Personal: std.True}
	AccessPersonalAndOther = Access{PersonalAndOther: std.True}
)

// Statement describes an EPP server’s data collection purpose, receipient(s), and retention policy.
type Statement struct {
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
}

// Purpose represents an EPP server’s purpose for data collection.
type Purpose struct {
	Admin        std.Bool `xml:"admin"`
	Contact      std.Bool `xml:"contact"`
	Provisioning std.Bool `xml:"provisioning"`
	Other        std.Bool `xml:"other"`
}

var (
	PurposeAdmin        = Purpose{Admin: std.True}
	PurposeContact      = Purpose{Contact: std.True}
	PurposeProvisioning = Purpose{Provisioning: std.True}
	PurposeOther        = Purpose{Other: std.True}
)

// Recipient represents an EPP server’s purpose for data collection.
type Recipient struct {
	Other     std.Bool `xml:"other"`
	Ours      *Ours    `xml:"ours"`
	Public    std.Bool `xml:"public"`
	Same      std.Bool `xml:"same"`
	Unrelated std.Bool `xml:"unrelated"`
}

// Ours represents an EPP server’s description of an <ours> recipient.
type Ours struct {
	Recipient string `xml:"recDesc"`
}

// MarshalXML impements the xml.Marshaler interface.
// Writes a single self-closing <ours/> if v.Recipient is not set.
func (v *Ours) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if v.Recipient == "" {
		return e.EncodeToken(xml.SelfClosingElement(start))
	}
	type T Ours
	return e.EncodeElement((*T)(v), start)
}

// Expiry defines an EPP server’s data retention duration.
type Expiry struct {
	Absolute *std.Time    `xml:"absolute"`
	Relative std.Duration `xml:"relative,omitempty"`
}
