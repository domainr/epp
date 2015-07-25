package epp

import (
	"bytes"
	"encoding/xml"

	"github.com/nbio/xx"
)

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (*DomainCheckResponse, error) {
	err := encodeDomainCheck(&c.buf, domains, c.Greeting.SupportsExtension(ExtFee))
	if err != nil {
		return nil, err
	}
	err = c.flushDataUnit()
	if err != nil {
		return nil, err
	}
	var res response_
	err = c.readResponse(&res)
	if err != nil {
		return nil, err
	}
	return &res.DomainCheckResponse, nil
}

func encodeDomainCheck(buf *bytes.Buffer, domains []string, extFee bool) error {
	buf.Reset()
	buf.Write(xmlCommandPrefix)
	buf.WriteString(`<check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	for _, domain := range domains {
		buf.WriteString(`<domain:check><domain:name>`)
		xml.EscapeText(buf, []byte(domain))
		buf.WriteString(`</domain:name></domain:check>`)
	}
	buf.WriteString(`</check>`)

	if extFee && false {
		// Extensions
		buf.WriteString(`<extension>`)

		// CentralNic fee extension
		buf.WriteString(`<fee:check xmlns:fee="urn:ietf:params:xml:ns:fee-0.5">`)
		for _, domain := range domains {
			buf.WriteString(`<fee:domain>`)
			buf.WriteString(`<fee:name>`)
			xml.EscapeText(buf, []byte(domain))
			buf.WriteString(`</fee:name>`)
			buf.WriteString(`<fee:command>create</fee:command>`)
			buf.WriteString(`</fee:domain>`)
		}
		buf.WriteString(`</fee:check>`)

		buf.WriteString(`</extension>`)
	}

	buf.Write(xmlCommandSuffix)
	return nil
}

type DomainCheckResponse struct {
	Checks  []DomainCheck
	Charges []DomainCharge
}

type DomainCheck struct {
	Domain    string
	Reason    string
	Available bool
}

type DomainCharge struct {
	Domain       string
	Category     string
	CategoryName string
}

func init() {
	path := "epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Checks = append(dcd.Checks, DomainCheck{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Domain = string(c.CharData)
		check.Available = c.AttrBool("", "avail")
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>reason", func(c *xx.Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Reason = string(c.CharData)
		return nil
	})
}

func init() {
	path := "epp > response > extension > http://www.unitedtld.com/epp/charge-1.0 chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Charges = append(dcd.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>set>category", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		charge.CategoryName = c.Attr("", "name")
		return nil
	})
}
