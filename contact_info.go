package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

type Addr struct {
	Street []string `xml:"street,omitempty"`
	City   string   `xml:"city"`
	PC     string   `xml:"pc"`
	SP     string   `xml:"sp,omitempty"`
	CC     string   `xml:"cc"`
}

type PostalInfo struct {
	Type string `xml:"type,attr"`
	Name string `xml:"name"`
	Org  string `xml:"org"`
	Addr Addr   `xml:"addr"`
}

type ContactStatus struct {
	S     string `xml:"s,attr"`
	Value string `xml:",chardata"`
}

// ContactInfoResponse represents an EPP <response> for a Contact check.
type ContactInfoResponse struct {
	XMLName      xml.Name        `xml:"epp"`
	Result       Result          `xml:"response>result"`
	ID           string          `xml:"response>resData>infData>id"`
	ROID         string          `xml:"response>resData>infData>roid"`
	Status       []ContactStatus `xml:"response>resData>infData>status"`
	PostalInfo   PostalInfo      `xml:"response>resData>infData>postalInfo"`
	Voice        string          `xml:"response>resData>infData>voice,omitempty"`
	Fax          string          `xml:"response>resData>infData>fax,omitempty"`
	Email        string          `xml:"response>resData>infData>email"`
	CreatedBy    string          `xml:"response>resData>infData>crID"`
	CreatedDate  string          `xml:"response>resData>infData>crDate"`
	UpdatedBy    string          `xml:"response>resData>infData>upID"`
	UpdatedDate  string          `xml:"response>resData>infData>upDate"`
	TransferDate string          `xml:"response>resData>infData>trDate"`
}

func (c *Conn) ContactInfo(roid string) (*ContactInfoResponse, error) {
	err := c.encodeContactInfo(roid)
	if err != nil {
		return nil, err
	}
	return c.processContactInfo(roid)
}

func (c *Conn) encodeContactInfo(roid string) error {
	escaped := &bytes.Buffer{}
	xml.EscapeText(escaped, []byte(roid))

	cmd := fmt.Sprintf(xmlCommandPrefix+`
	<info>
		<contact:info xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
			<contact:id>%s</contact:id>
		</contact:info>
	</info>`+xmlCommandSuffix, escaped)

	c.buf.Reset()

	_, err := c.buf.WriteString(cmd)

	return err
}

func (c *Conn) processContactInfo(roid string) (*ContactInfoResponse, error) {

	if err := c.flushDataUnit(); err != nil {
		return nil, err
	}

	if err := c.readDataUnit(); err != nil {
		return nil, err
	}

	dir := ContactInfoResponse{}
	if err := xml.Unmarshal(c.buf.Bytes(), &dir); err != nil {
		return nil, err
	}

	return &dir, nil
}
