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

	// The ARI price extension won't return both availability and price data
	// in the same response, so we have to make a separate request for price
	if c.Greeting.SupportsExtension(ExtPrice) {
		err := encodePriceCheck(&c.buf, domains)
		err = c.flushDataUnit()
		if err != nil {
			return nil, err
		}
		var res2 response_
		err = c.readResponse(&res2)
		if err != nil {
			return nil, err
		}
		res.DomainCheckResponse.Charges = res2.DomainCheckResponse.Charges
	}

	return &res.DomainCheckResponse, nil
}

func encodeDomainCheck(buf *bytes.Buffer, domains []string, extFee bool) error {
	buf.Reset()
	buf.WriteString(xmlCommandPrefix)
	buf.WriteString(`<check>`)
	buf.WriteString(`<domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	for _, domain := range domains {
		buf.WriteString(`<domain:name>`)
		xml.EscapeText(buf, []byte(domain))
		buf.WriteString(`</domain:name>`)
	}
	buf.WriteString(`</domain:check>`)
	buf.WriteString(`</check>`)

	if extFee {
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

	buf.WriteString(xmlCommandSuffix)
	return nil
}

func encodePriceCheck(buf *bytes.Buffer, domains []string) error {
	buf.Reset()
	buf.WriteString(xmlCommandPrefix)
	buf.WriteString(`<check>`)
	buf.WriteString(`<domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	for _, domain := range domains {
		buf.WriteString(`<domain:name>`)
		xml.EscapeText(buf, []byte(domain))
		buf.WriteString(`</domain:name>`)
	}
	buf.WriteString(`</domain:check>`)
	buf.WriteString(`</check>`)

	// Extensions
	buf.WriteString(`<extension>`)
	// ARI price extension
	buf.WriteString(`<price:check xmlns:price="urn:ar:params:xml:ns:price-1.1"></price:check>`)
	buf.WriteString(`</extension>`)

	buf.WriteString(xmlCommandSuffix)
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
	// Default EPP check data
	path := "epp > response > resData > " + ObjDomain + " chkData"
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

	// Scan charge-1.0 extension into Charges
	path = "epp > response > extension > " + ExtCharge + " chkData"
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

	// Scan fee-0.5 extension into Charges
	path = "epp > response > extension > " + ExtFee + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Charges = append(dcd.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		if c.AttrBool("", "premium") {
			charge.Category = "premium"
		}
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>fee", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.CategoryName = c.Attr("", "description")
		return nil
	})

	// Scan price-1.1 extension into Charges
	path = "epp > response > extension > " + ExtPrice + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Charges = append(dcd.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		if c.AttrBool("", "premium") {
			charge.Category = "premium"
		}
		return nil
	})
}
