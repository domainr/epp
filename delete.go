package epp

import (
	"bytes"
	"encoding/xml"
)

// DeleteDomain requests the deletion of a domain.
// https://tools.ietf.org/html/rfc5731#section-3.2.1
func (c *Conn) DeleteDomain(domain string, extData map[string]string) error {
	x, err := encodeDomainDelete(&c.Greeting, domain, extData)
	if err != nil {
		return err
	}
	err = c.writeRequest(x)
	if err != nil {
		return err
	}
	_, err = c.readResponse()
	return err
}

func encodeDomainDelete(greeting *Greeting, domain string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<delete><domain:delete xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name></domain:delete></delete>`)
	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// DeleteContact requests the deletion of a contact.
// https://tools.ietf.org/html/rfc5733#section-3.2.2
func (c *Conn) DeleteContact(id string, extData map[string]string) error {
	x, err := encodeContactDelete(&c.Greeting, id, extData)
	if err != nil {
		return err
	}
	err = c.writeRequest(x)
	if err != nil {
		return err
	}
	_, err = c.readResponse()
	return err
}

func encodeContactDelete(greeting *Greeting, id string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<delete><contact:delete xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>`)
	xml.EscapeText(buf, []byte(id))
	buf.WriteString(`</contact:id></contact:delete></delete>`)
	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}
