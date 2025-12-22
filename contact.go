package epp

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/nbio/xx"
)

// PostalInfo represents the postal address information for a contact.
type PostalInfo struct {
	Name   string // <contact:name>
	Org    string // <contact:org> (optional)
	Street string // <contact:street> (optional, up to 3)
	City   string // <contact:city>
	SP     string // <contact:sp> (optional)
	PC     string // <contact:pc> (optional)
	CC     string // <contact:cc>
}

// CreateContact requests the creation of a contact.
// https://tools.ietf.org/html/rfc5733#section-3.2.1
func (c *Conn) CreateContact(id string, email string, pi PostalInfo, voice string, auth string, extData map[string]string) (*ContactCreateResponse, error) {
	x, err := encodeContactCreate(&c.Greeting, id, email, pi, voice, auth, extData)
	if err != nil {
		return nil, err
	}
	err = c.writeRequest(x)
	if err != nil {
		return nil, err
	}
	res, err := c.readResponse()
	if err != nil {
		return nil, err
	}
	return &res.ContactCreateResponse, nil
}

func encodeContactCreate(greeting *Greeting, id string, email string, pi PostalInfo, voice string, auth string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<create><contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">`)
	buf.WriteString(`<contact:id>`)
	xml.EscapeText(buf, []byte(id))
	buf.WriteString(`</contact:id>`)

	// Postal Info "int" (international) is standard
	buf.WriteString(`<contact:postalInfo type="int">`)
	buf.WriteString(`<contact:name>`)
	xml.EscapeText(buf, []byte(pi.Name))
	buf.WriteString(`</contact:name>`)

	if pi.Org != "" {
		buf.WriteString(`<contact:org>`)
		xml.EscapeText(buf, []byte(pi.Org))
		buf.WriteString(`</contact:org>`)
	}

	buf.WriteString(`<contact:addr>`)
	if pi.Street != "" {
		buf.WriteString(`<contact:street>`)
		xml.EscapeText(buf, []byte(pi.Street))
		buf.WriteString(`</contact:street>`)
	}
	buf.WriteString(`<contact:city>`)
	xml.EscapeText(buf, []byte(pi.City))
	buf.WriteString(`</contact:city>`)

	if pi.SP != "" {
		buf.WriteString(`<contact:sp>`)
		xml.EscapeText(buf, []byte(pi.SP))
		buf.WriteString(`</contact:sp>`)
	}
	if pi.PC != "" {
		buf.WriteString(`<contact:pc>`)
		xml.EscapeText(buf, []byte(pi.PC))
		buf.WriteString(`</contact:pc>`)
	}

	buf.WriteString(`<contact:cc>`)
	xml.EscapeText(buf, []byte(pi.CC))
	buf.WriteString(`</contact:cc>`)
	buf.WriteString(`</contact:addr>`)
	buf.WriteString(`</contact:postalInfo>`)

	if voice != "" {
		buf.WriteString(`<contact:voice>`)
		xml.EscapeText(buf, []byte(voice))
		buf.WriteString(`</contact:voice>`)
	}

	buf.WriteString(`<contact:email>`)
	xml.EscapeText(buf, []byte(email))
	buf.WriteString(`</contact:email>`)

	if auth != "" {
		buf.WriteString(`<contact:authInfo><contact:pw>`)
		xml.EscapeText(buf, []byte(auth))
		buf.WriteString(`</contact:pw></contact:authInfo>`)
	}

	buf.WriteString(`</contact:create></create>`)
	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// ContactCreateResponse represents an EPP response for a domain create request.
type ContactCreateResponse struct {
	ID     string    // <contact:id>
	CrDate time.Time // <contact:crDate>
}

func init() {
	path := "epp > response > resData > " + ObjContact + " creData"
	scanResponse.MustHandleCharData(path+">id", func(c *xx.Context) error {
		ccr := &c.Value.(*Response).ContactCreateResponse
		ccr.ID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">crDate", func(c *xx.Context) error {
		ccr := &c.Value.(*Response).ContactCreateResponse
		var err error
		ccr.CrDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
}
