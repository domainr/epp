package protocol

import (
	"github.com/domainr/epp/internal/encoding/xml"
)

type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Command *Command `xml:"command,omitempty"`
}

// func (epp *EPP) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
// 	// start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: ns.EPP})
// 	start.Name.Space = ns.EPP
// 	start.Name.Local = "epp"
// 	type proxy EPP
// 	return e.EncodeElement((*proxy)(epp), start)
// }

type Command struct {
	Check *Check `xml:"check,omitempty"`
}

type Check struct {
	DomainCheck *DomainCheck `xml:"urn:ietf:params:xml:ns:domain-1.0 check,omitempty"`
}

type DomainCheck struct {
	DomainNames DomainNames `xml:"urn:ietf:params:xml:ns:domain-1.0 name,omitempty"`
}

func (dc *DomainCheck) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Space: "xmlns", Local: "domain"}, Value: start.Name.Space})
	type proxy DomainCheck
	return e.EncodeElement((*proxy)(dc), start)
}

type DomainNames []string

func (dn DomainNames) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name.Local = "domain:" + start.Name.Local
	start.Name.Space = ""
	type proxy DomainNames
	return e.EncodeElement((proxy)(dn), start)
}

func encodePrefixed(e *xml.Encoder, v interface{}, start xml.StartElement, prefix string) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:" + prefix}, Value: start.Name.Space})
	start.Name.Local = prefix + ":" + start.Name.Local
	start.Name.Space = ""
	return e.EncodeElement(v, start)
}
