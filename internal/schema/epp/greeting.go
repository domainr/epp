package epp

import (
	"reflect"
	"strings"

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
	Purpose   Purpose   `xml:"purpose"`
	Recipient Recipient `xml:"recipient"`
}

// Purpose represents an EPP server’s purpose for data collection.
type Purpose int

const (
	PurposeAdmin = 1 << iota
	PurposeContact
	PurposeProvisioning
	PurposeOther
)

type purposeTemplate struct {
	Admin        *struct{} `xml:"admin"`
	Contact      *struct{} `xml:"contact"`
	Provisioning *struct{} `xml:"prov"`
	Other        *struct{} `xml:"other"`
}

// MarshalXML implements xml.Marshaler.
func (p Purpose) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var b strings.Builder
	typ := reflect.TypeOf(purposeTemplate{})
	for i := 0; i < typ.NumField(); i++ {
		if p&(1<<i) == 0 {
			continue
		}
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		tag := f.Tag.Get("xml")
		if tag == "" || tag == "-" {
			continue
		}
		tokens := strings.Split(tag, ",")
		if tokens[0] == "" {
			continue
		}
		b.WriteByte('<')
		b.WriteString(tokens[0])
		b.WriteString("/>")
	}
	// var s string
	// if p&PurposeAdmin != 0 {
	// 	s += "<admin/>"
	// }
	// if p&PurposeContact != 0 {
	// 	s += "<contact/>"
	// }
	// if p&PurposeProvisioning != 0 {
	// 	s += "<prov/>"
	// }
	// if p&PurposeOther != 0 {
	// 	s += "<other/>"
	// }
	return e.EncodeElement(&raw.XML{Value: b.String()}, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (p *Purpose) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v purposeTemplate
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	*p = 0
	val := reflect.ValueOf(v)
	for i := 0; i < val.NumField(); i++ {
		if !val.Field(i).IsNil() {
			*p |= 1 << i
		}
	}
	// if v.Admin != nil {
	// 	*p |= PurposeAdmin
	// }
	// if v.Contact != nil {
	// 	*p |= PurposeContact
	// }
	// if v.Provisioning != nil {
	// 	*p |= PurposeProvisioning
	// }
	// if v.Other != nil {
	// 	*p |= PurposeOther
	// }
	return nil
}

// Recipient represents an EPP server’s purpose for data collection.
type Recipient int

const (
	RecipientOther = 1 << iota
	RecipientOurs
	RecipientPublic
	RecipientSame
	RecipientUnrelated
)

// MarshalXML implements xml.Marshaler.
func (r Recipient) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var s string
	if r&RecipientOther != 0 {
		s += "<other/>"
	}
	if r&RecipientOurs != 0 {
		s += "<ours/>"
	}
	if r&RecipientPublic != 0 {
		s += "<public/>"
	}
	if r&RecipientSame != 0 {
		s += "<same/>"
	}
	if r&RecipientUnrelated != 0 {
		s += "<unrelated/>"
	}
	return e.EncodeElement(&raw.XML{Value: s}, start)
}

// UnmarshalXML implements xml.Unmarshaler.
func (r *Recipient) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v struct {
		Other     *struct{} `xml:"other"`
		Ours      *struct{} `xml:"ours"`
		Public    *struct{} `xml:"public"`
		Same      *struct{} `xml:"same"`
		Unrelated *struct{} `xml:"unrelated"`
	}
	err := d.DecodeElement(&v, &start)
	if err != nil {
		return err
	}
	*r = 0
	if v.Other != nil {
		*r |= RecipientOther
	}
	if v.Ours != nil {
		*r |= RecipientOurs
	}
	if v.Public != nil {
		*r |= RecipientPublic
	}
	if v.Same != nil {
		*r |= RecipientSame
	}
	if v.Unrelated != nil {
		*r |= RecipientUnrelated
	}
	return nil
}
