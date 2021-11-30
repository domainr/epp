package epp

import (
	"github.com/domainr/epp/internal/schema/date"
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

/*
<dcp>
	<access>
		<all/>
	</access>
	<statement>
		<purpose>
			<admin/>
			<other/>
			<prov/>
		</purpose>
		<recipient>
			<ours/>
			<public/>
			<unrelated/>
		</recipient>
		<retention>
			<indefinite/>
		</retention>
	</statement>
</dcp>
*/

// Statement describes an EPP server’s data collection purpose, receipient(s), and retention policy.
type Statement struct {
	Purpose Purpose `xml:"purpose"`
}

type Purpose int

const (
	PurposeAdmin = 1 << iota
	PurposeContact
	PurposeProvisioning
	PurposeOther
)

// MarshalXML implements xml.Marshaler.
func (p Purpose) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string
	if p&PurposeAdmin != 0 {
		s = "<admin/>"
	}
	if p&PurposeContact != 0 {
		s += "<contact/>"
	}
	if p&PurposeProvisioning != 0 {
		s += "<prov/>"
	}
	if p&PurposeOther != 0 {
		s += "<other/>"
	}
	return e.EncodeElement(&raw.XML{Value: s}, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (p *Purpose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Admin        *struct{} `xml:"admin"`
		Contact      *struct{} `xml:"contact"`
		Provisioning *struct{} `xml:"prov"`
		Other        *struct{} `xml:"other"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	*p = 0
	if v.Admin != nil {
		*p |= PurposeAdmin
	}
	if v.Contact != nil {
		*p |= PurposeContact
	}
	if v.Provisioning != nil {
		*p |= PurposeProvisioning
	}
	if v.Other != nil {
		*p |= PurposeOther
	}
	return nil
}
