package epp

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/nbio/xx"
)

// RenewDomain requests the renewal of a domain.
// https://tools.ietf.org/html/rfc5731#section-3.2.2
func (c *Conn) RenewDomain(domain string, curExpDate time.Time, period int, unit string, extData map[string]string) (*DomainRenewResponse, error) {
	x, err := encodeDomainRenew(&c.Greeting, domain, curExpDate, period, unit, extData)
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
	return &res.DomainRenewResponse, nil
}

func encodeDomainRenew(greeting *Greeting, domain string, curExpDate time.Time, period int, unit string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<renew><domain:renew xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	buf.WriteString(`<domain:name>`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name>`)

	buf.WriteString(`<domain:curExpDate>`)
	buf.WriteString(curExpDate.Format("2006-01-02")) // Date only required? RFC says date or dateTime. usually date for existing. simple date format.
	// Actually RFC 5731 says "date" type which is YYYY-MM-DD.
	buf.WriteString(`</domain:curExpDate>`)

	if period > 0 {
		buf.WriteString(`<domain:period unit="`)
		buf.WriteString(unit)
		buf.WriteString(`">`)
		buf.WriteString(xmlInt(period))
		buf.WriteString(`</domain:period>`)
	}

	buf.WriteString(`</domain:renew></renew>`)

	if fee, ok := extData["fee:fee"]; ok {
		buf.WriteString(`<extension>`)
		buf.WriteString(`<fee:renew xmlns:fee="`)
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
		buf.WriteString(`</fee:renew>`)
		buf.WriteString(`</extension>`)
	}

	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// DomainRenewResponse represents an EPP response for a domain renew request.
type DomainRenewResponse struct {
	Domain string    // <domain:name>
	ExDate time.Time // <domain:exDate>
}

func init() {
	path := "epp > response > resData > " + ObjDomain + " renData"
	scanResponse.MustHandleCharData(path+">name", func(c *xx.Context) error {
		drr := &c.Value.(*Response).DomainRenewResponse
		drr.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">exDate", func(c *xx.Context) error {
		drr := &c.Value.(*Response).DomainRenewResponse
		var err error
		drr.ExDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
}
