package epp

import (
	"bytes"
	"encoding/xml"
	"strings"

	"github.com/nbio/xx"
)

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (*DomainCheckResponse, error) {
	err := c.encodeDomainCheck(domains, nil)
	if err != nil {
		return nil, err
	}
	return c.processDomainCheck(domains)
}

// CheckDomainExtensions allows specifying extension data for the following:
//  - "neulevel:unspec": a string of the Key=Value data for the unspec tag
//  - "launch:phase": a string of the launch phase
func (c *Conn) CheckDomainExtensions(domains []string, extData map[string]string) (*DomainCheckResponse, error) {
	err := c.encodeDomainCheck(domains, extData)
	if err != nil {
		return nil, err
	}
	return c.processDomainCheck(domains)
}

func (c *Conn) encodeDomainCheck(domains []string, extData map[string]string) error {
	c.buf.Reset()
	c.buf.WriteString(xmlCommandPrefix)
	c.buf.WriteString(`<check>`)
	c.buf.WriteString(`<domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	for _, domain := range domains {
		c.buf.WriteString(`<domain:name>`)
		xml.EscapeText(&c.buf, []byte(domain))
		c.buf.WriteString(`</domain:name>`)
	}
	c.buf.WriteString(`</domain:check>`)
	c.buf.WriteString(`</check>`)

	greeting := c.Greeting

	var feeURN string
	switch {
	// Versions 0.8-0.9 require the returned class to be "standard" for
	// non-premium domains
	case greeting.SupportsExtension(ExtFee08):
		feeURN = ExtFee08
	case greeting.SupportsExtension(ExtFee09):
		feeURN = ExtFee09
	// Version 0.5 has an attribute premium="1" for premium domains
	case greeting.SupportsExtension(ExtFee05):
		feeURN = ExtFee05
	// Version 0.6 and 0.7 don't have a standard way of detecting premiums,
	// so instead there must be matching done on class names
	case greeting.SupportsExtension(ExtFee06):
		feeURN = ExtFee06
	case greeting.SupportsExtension(ExtFee07):
		feeURN = ExtFee07
	}

	supportsLaunch := extData["launch:phase"] != "" && greeting.SupportsExtension(ExtLaunch)
	supportsNeulevel := extData["neulevel:unspec"] != "" && (greeting.SupportsExtension(ExtNeulevel) || greeting.SupportsExtension(ExtNeulevel10))

	hasExtension := feeURN != "" || supportsLaunch || supportsNeulevel

	if hasExtension {
		c.buf.WriteString(`<extension>`)
	}

	if supportsLaunch {
		c.buf.WriteString(`<launch:check xmlns:launch="` + ExtLaunch + `" type="avail">`)
		c.buf.WriteString(`<launch:phase>` + extData["launch:phase"] + `</launch:phase>`)
		c.buf.WriteString(`</launch:check>`)
	}

	if supportsNeulevel {
		c.buf.WriteString(`<neulevel:extension xmlns:neulevel="` + ExtNeulevel10 + `">`)
		c.buf.WriteString(`<neulevel:unspec>` + extData["neulevel:unspec"] + `</neulevel:unspec>`)
		c.buf.WriteString(`</neulevel:extension>`)
	}

	if len(feeURN) > 0 {
		c.buf.WriteString(`<fee:check xmlns:fee="` + feeURN + `">`)
		for _, domain := range domains {
			if feeURN == ExtFee09 {
				// Version 0.9 changes the XML structure
				c.buf.WriteString(`<fee:object objURI="urn:ietf:params:xml:ns:domain-1.0">`)
				c.buf.WriteString(`<fee:objID element="name">`)
				xml.EscapeText(&c.buf, []byte(domain))
				c.buf.WriteString(`</fee:objID>`)
				c.buf.WriteString(`<fee:command>create</fee:command>`)
				c.buf.WriteString(`</fee:object>`)
			} else {
				c.buf.WriteString(`<fee:domain>`)
				c.buf.WriteString(`<fee:name>`)
				xml.EscapeText(&c.buf, []byte(domain))
				c.buf.WriteString(`</fee:name>`)
				c.buf.WriteString(`<fee:command>create</fee:command>`)
				c.buf.WriteString(`</fee:domain>`)
			}
		}
		c.buf.WriteString(`</fee:check>`)
	}

	if hasExtension {
		c.buf.WriteString(`</extension>`)
	}

	c.buf.WriteString(xmlCommandSuffix)
	return nil
}

func (c *Conn) processDomainCheck(domains []string) (*DomainCheckResponse, error) {
	err := c.flushDataUnit()
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

	path = "epp > response > extension > " + ExtFee05 + " chkData"
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

	path = "epp > response > extension > " + ExtFee06 + " chkData"
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
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		className := strings.ToLower(string(c.CharData))
		isDefault := strings.Index(className, "default") != -1
		isNormal := strings.Index(className, "normal") != -1
		isDiscount := strings.Index(className, "discount") != -1
		if isDefault == false && isNormal == false && isDiscount == false {
			charge.Category = "premium"
		}
		return nil
	})

	path = "epp > response > extension > " + ExtFee07 + " chkData"
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
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		return nil
	})

	path = "epp > response > extension > " + ExtFee08 + " chkData"
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
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if string(c.CharData) != "standard" {
			charge.Category = "premium"
		}
		return nil
	})

	path = "epp > response > extension > " + ExtFee09 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Charges = append(dcd.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>objID", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*response_).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if string(c.CharData) != "standard" {
			charge.Category = "premium"
		}
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

	// Scan neulevel-1.0 extension
	path = "epp > response > extension > " + ExtNeulevel10 + " extension > unspec"
	scanResponse.MustHandleCharData(path, func(c *xx.Context) error {
		dcr := &c.Value.(*response_).DomainCheckResponse
		if len(dcr.Checks) == 0 {
			return nil
		}

		check := &dcr.Checks[len(dcr.Checks)-1]
		charge := DomainCharge{Domain: check.Domain}
		data := string(c.CharData)
		pairs := strings.Split(data, " ")
		for _, pair := range pairs {
			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 && parts[0] == "TierName" {
				charge.Category = parts[1]
				break
			}
		}
		dcr.Charges = append(dcr.Charges, charge)

		return nil
	})
}
