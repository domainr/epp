package epp

import (
	"github.com/domainr/epp/internal/schema/date"
	"github.com/domainr/epp/internal/schema/option"
	"github.com/nbio/xml"
)

// Greeting represents an EPP server <greeting> message as defined in RFC 5730.
type Greeting struct {
	ServerName  string       `xml:"svID,omitempty"`
	ServerDate  *date.Time   `xml:"svDate"`
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
}

// Access represents an EPP server’s scope of data access as defined in RFC 5730.
type Access uint64

const (
	AccessNull = iota
	AccessNone
	AccessPersonal
	AccessOther
	AccessPersonalAndOther
	AccessAll
)

var accessNames = map[uint64]string{
	AccessNull:             "null",
	AccessNone:             "none",
	AccessPersonal:         "personal",
	AccessOther:            "other",
	AccessPersonalAndOther: "personalAndOther",
	AccessAll:              "all",
}

var accessValues = option.Values(accessNames)

// MarshalXML implements xml.Marshaler.
func (v Access) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return option.EncodeOne(e, start, uint64(v), accessNames)
}

// UnmarshalXML implements xml.Unmarshaler.
func (v *Access) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return option.DecodeOne((*uint64)(v), d, accessValues)
}

// Statement describes an EPP server’s data collection purpose, receipient(s), and retention policy.
type Statement struct {
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
}

// Purpose represents an EPP server’s purpose for data collection.
type Purpose uint64

const (
	PurposeAdmin = 1 << iota
	PurposeContact
	PurposeProvisioning
	PurposeOther
)

var purposeNames = map[uint64]string{
	PurposeAdmin:        "admin",
	PurposeContact:      "contact",
	PurposeProvisioning: "provisioning",
	PurposeOther:        "other",
}

var purposeValues = option.Values(purposeNames)

// MarshalXML implements xml.Marshaler.
func (v Purpose) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return option.Encode(e, start, uint64(v), purposeNames)
}

// UnmarshalXML implements xml.Unmarshaler.
func (v *Purpose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return option.Decode((*uint64)(v), d, purposeValues)
}

// Recipient represents an EPP server’s purpose for data collection.
type Recipient uint64

const (
	RecipientOther = 1 << iota
	RecipientOurs
	RecipientPublic
	RecipientSame
	RecipientUnrelated
)

var recipientNames = map[uint64]string{
	RecipientOther:     "other",
	RecipientOurs:      "ours",
	RecipientPublic:    "public",
	RecipientSame:      "same",
	RecipientUnrelated: "unrelated",
}

var recipientValues = option.Values(recipientNames)

// MarshalXML implements xml.Marshaler.
func (v Recipient) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return option.Encode(e, start, uint64(v), recipientNames)
}

// UnmarshalXML implements xml.Unmarshaler.
func (v *Recipient) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	return option.Decode((*uint64)(v), d, recipientValues)
}
