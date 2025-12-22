package epp

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/nbio/xx"
)

// TransferDomain requests a transfer operation for a domain.
// https://tools.ietf.org/html/rfc5731#section-3.2.4
func (c *Conn) TransferDomain(op string, domain string, period int, unit string, auth string, extData map[string]string) (*DomainTransferResponse, error) {
	x, err := encodeDomainTransfer(&c.Greeting, op, domain, period, unit, auth, extData)
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
	return &res.DomainTransferResponse, nil
}

func encodeDomainTransfer(greeting *Greeting, op string, domain string, period int, unit string, auth string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<transfer op="`)
	buf.WriteString(op)
	buf.WriteString(`">`)
	buf.WriteString(`<domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	buf.WriteString(`<domain:name>`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name>`)

	if period > 0 {
		buf.WriteString(`<domain:period unit="`)
		buf.WriteString(unit)
		buf.WriteString(`">`)
		buf.WriteString(xmlInt(period))
		buf.WriteString(`</domain:period>`)
	}

	if auth != "" {
		buf.WriteString(`<domain:authInfo><domain:pw>`)
		xml.EscapeText(buf, []byte(auth))
		buf.WriteString(`</domain:pw></domain:authInfo>`)
	}

	buf.WriteString(`</domain:transfer></transfer>`)

	// Extensions (e.g. fee)
	if fee, ok := extData["fee:fee"]; ok {
		buf.WriteString(`<extension>`)
		buf.WriteString(`<fee:transfer xmlns:fee="`)
		buf.WriteString(ExtFee10)
		buf.WriteString(`">`)
		if currency, ok := extData["fee:currency"]; ok {
			buf.WriteString(`<fee:currency>`)
			buf.WriteString(currency)
			buf.WriteString(`</fee:currency>`)
		}
		buf.WriteString(`<fee:fee>`)
		buf.WriteString(fee)
		buf.WriteString(`</fee:fee>`)
		buf.WriteString(`</fee:transfer>`)
		buf.WriteString(`</extension>`)
	}

	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// DomainTransferResponse represents an EPP response for a domain transfer request.
type DomainTransferResponse struct {
	Domain string    // <domain:name>
	Status string    // <domain:trStatus>
	REID   string    // <domain:reID>
	REDate time.Time // <domain:reDate>
	ACID   string    // <domain:acID>
	ACDate time.Time // <domain:acDate>
	ExDate time.Time // <domain:exDate>
}

func init() {
	path := "epp > response > resData > " + ObjDomain + " trnData"
	scanResponse.MustHandleCharData(path+">name", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">trStatus", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.Status = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">reID", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.REID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">reDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		var err error
		dtr.REDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">acID", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.ACID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">acDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		var err error
		dtr.ACDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">exDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		var err error
		dtr.ExDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
}
