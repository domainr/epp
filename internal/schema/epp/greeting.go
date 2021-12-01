package epp

import (
	"github.com/domainr/epp/internal/schema/date"
	"github.com/domainr/epp/internal/schema/option"
	"github.com/domainr/epp/internal/schema/raw"
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
type Access string

const (
	AccessNull             Access = "null"
	AccessNone             Access = "none"
	AccessPersonal         Access = "personal"
	AccessOther            Access = "other"
	AccessPersonalAndOther Access = "personalAndOther"
	AccessAll              Access = "all"
)

// MarshalXML implements xml.Marshaler.
func (a Access) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if a == "" {
		return nil
	}
	return e.EncodeElement(&raw.XML{Value: "<" + string(a) + "/>"}, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (a *Access) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Null             *struct{} `xml:"null"`
		None             *struct{} `xml:"none"`
		Personal         *struct{} `xml:"personal"`
		Other            *struct{} `xml:"other"`
		PersonalAndOther *struct{} `xml:"personalAndOther"`
		All              *struct{} `xml:"all"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	switch {
	case v.Null != nil:
		*a = "null"
	case v.None != nil:
		*a = "none"
	case v.Personal != nil:
		*a = "personal"
	case v.Other != nil:
		*a = "other"
	case v.PersonalAndOther != nil:
		*a = "personalAndOther"
	case v.All != nil:
		*a = "all"
	}
	return nil
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
