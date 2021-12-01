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
type Recipient struct {
	Other     *struct{} `xml:"other,selfclosing"`
	Ours      *Ours     `xml:"ours"`
	Public    *struct{} `xml:"public,selfclosing"`
	Same      *struct{} `xml:"same,selfclosing"`
	Unrelated *struct{} `xml:"unrelated,selfclosing"`
}

// MarshalXML implements xml.Marshaler.
func (v *Recipient) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if v.Ours != nil && v.Ours.Recipient != "" {
		type T Recipient
		return e.EncodeElement((*T)(v), start)
	}
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
