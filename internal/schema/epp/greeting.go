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
	Null             *struct{} `xml:"null,selfclosing"`
	All              *struct{} `xml:"all,selfclosing"`
	None             *struct{} `xml:"none,selfclosing"`
	Other            *struct{} `xml:"other,selfclosing"`
	Personal         *struct{} `xml:"personal,selfclosing"`
	PersonalAndOther *struct{} `xml:"personalAndOther,selfclosing"`
}

var (
	AccessNull             = Access{Null: &struct{}{}}
	AccessAll              = Access{All: &struct{}{}}
	AccessNone             = Access{None: &struct{}{}}
	AccessOther            = Access{Other: &struct{}{}}
	AccessPersonal         = Access{Personal: &struct{}{}}
	AccessPersonalAndOther = Access{PersonalAndOther: &struct{}{}}
)

// Statement describes an EPP server’s data collection purpose, receipient(s), and retention policy.
type Statement struct {
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
}

// Purpose represents an EPP server’s purpose for data collection.
type Purpose struct {
	Admin        *struct{} `xml:"admin,selfclosing"`
	Contact      *struct{} `xml:"contact,selfclosing"`
	Provisioning *struct{} `xml:"provisioning,selfclosing"`
	Other        *struct{} `xml:"other,selfclosing"`
}

var (
	PurposeAdmin        = Purpose{Admin: &struct{}{}}
	PurposeContact      = Purpose{Contact: &struct{}{}}
	PurposeProvisioning = Purpose{Provisioning: &struct{}{}}
	PurposeOther        = Purpose{Other: &struct{}{}}
)

// Recipient represents an EPP server’s purpose for data collection.
type Recipient struct {
	Other     *struct{} `xml:"other,selfclosing"`
	Ours      *Ours     `xml:"ours"`
	Public    *struct{} `xml:"public,selfclosing"`
	Same      *struct{} `xml:"same,selfclosing"`
	Unrelated *struct{} `xml:"unrelated,selfclosing"`
}

// MarshalXML implements xml.Marshaler.
func (v *Recipient) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	// If v.Ours.Recipient contains a value, emit a non-self-closing <ours> tag.
	if v.Ours != nil && v.Ours.Recipient != "" {
		type T Recipient
		return e.EncodeElement((*T)(v), start)
	}

	// Otherwise, emit a self-closing <ours/> tag.
	// This hack takes advantage of the fact that Go will let you typecast
	// between two structs that differ only by struct tags.
	// See https://go-review.googlesource.com/c/go/+/24190/.
	type T struct {
		Other     *struct{} `xml:"other,selfclosing"`
		Ours      *Ours     `xml:"ours,selfclosing"`
		Public    *struct{} `xml:"public,selfclosing"`
		Same      *struct{} `xml:"same,selfclosing"`
		Unrelated *struct{} `xml:"unrelated,selfclosing"`
	}
	return e.EncodeElement((*T)(v), start)
}

// Ours represents an EPP server’s description of an <ours> recipient.
type Ours struct {
	Recipient string `xml:"recDesc"`
}

// Expiry defines an EPP server’s data retention duration.
type Expiry struct {
	Absolute *std.Time    `xml:"absolute"`
	Relative std.Duration `xml:"relative"`
}
