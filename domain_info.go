package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type DomainStatus struct {
	S     string `xml:"s,attr"`
	Value string `xml:",chardata"`
}
type DomainContact struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",chardata"`
}

type AuthInfo struct {
	Date     string `xml:"date"`
	ClientID string `xml:"clID"`
	Code     string `xml:"code"`
}

// DomainInfoResponse represents an EPP <response> for a domain check.
type DomainInfoResponse struct {
	XMLName      xml.Name        `xml:"epp"`
	Result       Result          `xml:"response>result"`
	Name         string          `xml:"response>resData>infData>name"`
	ROID         string          `xml:"response>resData>infData>roid"`
	Status       []DomainStatus  `xml:"response>resData>infData>status"`
	CLID         string          `xml:"response>resData>infData>clID"`
	Registrant   DomainContact   `xml:"response>resData>infData>registrant"`
	Contacts     []DomainContact `xml:"response>resData>infData>contact,omitempty"`
	NS           []string        `xml:"response>resData>infData>ns,omitempty"`
	Hosts        []string        `xml:"response>resData>infData>host,omitempty"`
	CreatedBy    string          `xml:"response>resData>infData>crID"`
	CreatedDate  string          `xml:"response>resData>infData>crDate"`
	Expiration   string          `xml:"response>resData>infData>exDate"`
	UpdatedBy    string          `xml:"response>resData>infData>upID"`
	UpdatedDate  string          `xml:"response>resData>infData>upDate"`
	TransferDate string          `xml:"response>resData>infData>trDate"`
	AuID         string          `xml:"response>resData>infData>auID,omitempty"`
	AuthInfo     AuthInfo        `xml:"response>resData>infData>authInfo"`
}

func (c *Conn) DomainInfo(domain string) (*DomainInfoResponse, error) {
	err := c.encodeDomainInfo(domain)
	if err != nil {
		return nil, err
	}
	return c.processDomainInfo(domain)
}

func (c *Conn) encodeDomainInfo(domain string) error {
	escaped := &bytes.Buffer{}
	xml.EscapeText(escaped, []byte(domain))

	cmd := fmt.Sprintf(xmlCommandPrefix+`
	<info>
		<domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
			<domain:name>%s</domain:name>
		</domain:info>
	</info>`+xmlCommandSuffix, escaped)

	c.buf.Reset()

	_, err := c.buf.WriteString(cmd)

	return err
}

func (c *Conn) processDomainInfo(domain string) (*DomainInfoResponse, error) {

	if err := c.flushDataUnit(); err != nil {
		return nil, err
	}

	if err := c.readDataUnit(); err != nil {
		return nil, err
	}

	dir := DomainInfoResponse{}
	if err := xml.Unmarshal(c.buf.Bytes(), &dir); err != nil {
		return nil, err
	}

	return &dir, nil
}
