package epp

import (
	"bytes"
	"encoding/xml"
)

// RestoreDomain requests the restoration of a domain (usually via RGP extension).
// This is actually an <update> command with an RGP extension <restore> op.
func (c *Conn) RestoreDomain(domain string, extData map[string]string) (*DomainUpdateResponse, error) {
	x, err := encodeDomainRestore(&c.Greeting, domain, extData)
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
	return &res.DomainUpdateResponse, nil
}

func encodeDomainRestore(greeting *Greeting, domain string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<update><domain:update xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	buf.WriteString(`<domain:name>`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name>`)
	buf.WriteString(`</domain:update></update>`)

	// RGP Extension for restore
	// https://tools.ietf.org/html/rfc3915
	buf.WriteString(`<extension><rgp:update xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0">`)
	buf.WriteString(`<rgp:restore op="request"/>`)
	buf.WriteString(`</rgp:update></extension>`)

	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// DomainUpdateResponse might be generic, but for now defining it here if not exists.
// Logic check: does DomainUpdateResponse exist? No.
type DomainUpdateResponse struct {
	// Usually empty resData for update?
	// RFC 5731 says <domain:upData> is optional and currently not defined to return anything useful other than basic success.
}

func init() {
	// No specific scan logic needed for update response data yet as it is usually empty
}
