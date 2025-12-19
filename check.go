package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nbio/xx"
)

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (*DomainCheckResponse, error) {
	return c.CheckDomainExtensions(domains, nil)
}

// CheckDomainExtensions allows specifying extension data for the following:
//   - "neulevel:unspec": a string of the Key=Value data for the unspec tag
//   - "launch:phase": a string of the launch phase
func (c *Conn) CheckDomainExtensions(domains []string, extData map[string]string) (*DomainCheckResponse, error) {
	x, err := encodeDomainCheck(&c.Greeting, domains, extData)
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

	// The ARI price extension won't return both availability and price data
	// in the same response, so we have to make a separate request for price
	if c.Greeting.SupportsExtension(ExtPrice) {
		x, err = encodePriceCheck(domains)
		if err != nil {
			return nil, err
		}
		err = c.writeRequest(x)
		if err != nil {
			return nil, err
		}
		res2, err := c.readResponse()
		if err != nil {
			return nil, err
		}
		res.DomainCheckResponse.Charges = res2.DomainCheckResponse.Charges
	}

	return &res.DomainCheckResponse, nil
}

func encodeDomainCheck(greeting *Greeting, domains []string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<check>`)
	buf.WriteString(`<domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">`)
	for _, domain := range domains {
		buf.WriteString(`<domain:name>`)
		xml.EscapeText(buf, []byte(domain))
		buf.WriteString(`</domain:name>`)
	}
	buf.WriteString(`</domain:check>`)
	buf.WriteString(`</check>`)

	var feeURN string
	switch {
	case greeting.SupportsExtension(ExtFee10):
		feeURN = ExtFee10
	case greeting.SupportsExtension(ExtFee21):
		feeURN = ExtFee21
	case greeting.SupportsExtension(ExtFee11):
		feeURN = ExtFee11
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
	supportsFeePhase := extData["fee:phase"] != ""
	supportsNeulevel := extData["neulevel:unspec"] != "" && (greeting.SupportsExtension(ExtNeulevel) || greeting.SupportsExtension(ExtNeulevel10))
	supportsNamestore := extData["namestoreExt:subProduct"] != "" && greeting.SupportsExtension(ExtNamestore)

	hasExtension := feeURN != "" || supportsLaunch || supportsNeulevel || supportsNamestore

	if hasExtension {
		buf.WriteString(`<extension>`)
	}

	// https://www.verisign.com/assets/epp-sdk/verisign_epp-extension_namestoreext_v01.html
	if supportsNamestore {
		buf.WriteString(`<namestoreExt:namestoreExt xmlns:namestoreExt="`)
		buf.WriteString(ExtNamestore)
		buf.WriteString(`">`)
		buf.WriteString(`<namestoreExt:subProduct>`)
		buf.WriteString(extData["namestoreExt:subProduct"])
		buf.WriteString(`</namestoreExt:subProduct>`)
		buf.WriteString(`</namestoreExt:namestoreExt>`)
	}

	if supportsLaunch {
		buf.WriteString(`<launch:check xmlns:launch="`)
		buf.WriteString(ExtLaunch)
		buf.WriteString(`" type="avail">`)
		buf.WriteString(`<launch:phase>`)
		buf.WriteString(extData["launch:phase"])
		buf.WriteString(`</launch:phase>`)
		buf.WriteString(`</launch:check>`)
	}

	if supportsNeulevel {
		buf.WriteString(`<neulevel:extension xmlns:neulevel="`)
		buf.WriteString(ExtNeulevel10)
		buf.WriteString(`">`)
		buf.WriteString(`<neulevel:unspec>`)
		buf.WriteString(extData["neulevel:unspec"])
		buf.WriteString(`</neulevel:unspec>`)
		buf.WriteString(`</neulevel:extension>`)
	}

	if len(feeURN) > 0 {
		buf.WriteString(`<fee:check xmlns:fee="`)
		buf.WriteString(feeURN)
		buf.WriteString(`">`)
		feePhase := ""
		if supportsFeePhase {
			feePhase = fmt.Sprintf(" phase=%q", extData["fee:phase"])
		}

		// ExtFee10 uses global commands on this server (and typically in fee-1.0 when not using objects)
		if feeURN == ExtFee10 {
			buf.WriteString(`<fee:command name="create"/>`)
			buf.WriteString(`<fee:command name="renew"/>`)
			buf.WriteString(`<fee:command name="restore"/>`)
			buf.WriteString(`<fee:command name="transfer"/>`)
		} else {
			// Other versions iterate per domain (or we preserve existing behavior for them)
			for _, domain := range domains {
				switch feeURN {
				case ExtFee09: // Version 0.9 changes the XML structure
					buf.WriteString(`<fee:object objURI="urn:ietf:params:xml:ns:domain-1.0">`)
					buf.WriteString(`<fee:objID element="name">`)
					xml.EscapeText(buf, []byte(domain))
					buf.WriteString(`</fee:objID>`)
					buf.WriteString(fmt.Sprintf(`<fee:command%s>create</fee:command>`, feePhase))
					buf.WriteString(`</fee:object>`)
				case ExtFee11: // https://tools.ietf.org/html/draft-brown-epp-fees-07#section-5.1.1
					buf.WriteString(fmt.Sprintf(`<fee:command%s>create</fee:command>`, feePhase))
				case ExtFee21: // Version 0.21 changes the XML structure
					buf.WriteString(`<fee:command name="create"/>`)
				default:
					buf.WriteString(`<fee:domain>`)
					buf.WriteString(`<fee:name>`)
					xml.EscapeText(buf, []byte(domain))
					buf.WriteString(`</fee:name>`)
					buf.WriteString(fmt.Sprintf(`<fee:command%s>create</fee:command>`, feePhase))
					buf.WriteString(`</fee:domain>`)
				}
			}
		}
		buf.WriteString(`</fee:check>`)
	}

	if hasExtension {
		buf.WriteString(`</extension>`)
	}

	buf.WriteString(xmlCommandSuffix)

	return buf.Bytes(), nil
}

func encodePriceCheck(domains []string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
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
	return buf.Bytes(), nil
}

// DomainCheckResponse represents an EPP <response> for a domain check.
type DomainCheckResponse struct {
	Domain  string
	Checks  []DomainCheck
	Charges []DomainCharge
}

// DomainCheck represents an EPP <chkData> and associated extension data.
type DomainCheck struct {
	Domain    string
	Reason    string
	Available bool
}

// DomainCharge represents various EPP charge and fee extension data.
type DomainCharge struct {
	Domain        string
	Category      string
	CategoryName  string
	CreatePrice   string
	RenewPrice    string
	RestorePrice  string
	TransferPrice string
	Currency      string
}

func init() {
	// Default EPP check data
	path := "epp > response > resData > " + ObjDomain + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Checks = append(dcr.Checks, DomainCheck{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		checks := c.Value.(*Response).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Domain = string(c.CharData)
		check.Available = c.AttrBool("", "avail")
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>reason", func(c *xx.Context) error {
		checks := c.Value.(*Response).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Reason = string(c.CharData)
		return nil
	})

	// Scan charge-1.0 extension into Charges
	path = "epp > response > extension > " + ExtCharge + " chkData"
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		c.Value.(*Response).DomainCheckResponse.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleStartElement(path+">cd>set", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>set>category", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = c.Value.(*Response).DomainCheckResponse.Domain
		charge.Category = string(c.CharData)
		charge.CategoryName = c.Attr("", "name")
		return nil
	})

	path = "epp > response > extension > " + ExtFee05 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		if c.AttrBool("", "premium") {
			charge.Category = "premium"
		}
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>fee", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.CategoryName = c.Attr("", "description")
		return nil
	})

	path = "epp > response > extension > " + ExtFee06 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		className := strings.ToLower(string(c.CharData))
		switch {
		case strings.Contains(className, "default"):
		case strings.Contains(className, "normal"):
		case strings.Contains(className, "discount"):
		case strings.Contains(className, "standard"):
		default:
			charge.Category = "premium"
		}
		return nil
	})

	path = "epp > response > extension > " + ExtFee07 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		return nil
	})

	path = "epp > response > extension > " + ExtFee08 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if string(c.CharData) != "standard" {
			charge.Category = "premium"
		}
		return nil
	})

	path = "epp > response > extension > " + ExtFee09 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>objID", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if string(c.CharData) != "standard" {
			charge.Category = "premium"
		}
		return nil
	})

	path = "epp > response > extension > " + ExtFee10 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if string(c.CharData) != "standard" {
			charge.Category = "premium"
		}
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>command>fee", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		if charge.CategoryName == "" {
			charge.CategoryName = c.Attr("", "description") // Fallback
		}

		cmdName := c.Parent.Attr("", "name")
		switch cmdName {
		case "create":
			charge.CreatePrice = string(c.CharData)
		case "renew":
			charge.RenewPrice = string(c.CharData)
		case "restore":
			charge.RestorePrice = string(c.CharData)
		case "transfer":
			charge.TransferPrice = string(c.CharData)
		}
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>currency", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Currency = string(c.CharData)
		return nil
	})

	path = "epp > response > extension > " + ExtFee11 + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>objID", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>class", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
		charge := &charges[len(charges)-1]
		charge.Category = string(c.CharData)
		return nil
	})

	// Scan fee-0.21 phase and subphase into Charges Category and CategoryName, respectively
	// FIXME: stop mangling fee extensions into charges
	path = "epp > response > extension > " + ExtFee21 + " chkData > cd > command > fee"
	scanResponse.MustHandleCharData(path, func(c *xx.Context) error {
		if c.Parent.Attr("", "name") != "create" {
			return nil
		}
		dcr := &c.Value.(*Response).DomainCheckResponse
		check := &dcr.Checks[len(dcr.Checks)-1]
		charge := DomainCharge{
			Domain:       check.Domain,
			Category:     c.Parent.Attr("", "phase"),
			CategoryName: c.Parent.Attr("", "subphase"),
		}
		dcr.Charges = append(dcr.Charges, charge)
		return nil
	})

	// Scan price-1.1 extension into Charges
	path = "epp > response > extension > " + ExtPrice + " chkData"
	scanResponse.MustHandleStartElement(path+">cd", func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse
		dcr.Charges = append(dcr.Charges, DomainCharge{})
		return nil
	})
	scanResponse.MustHandleCharData(path+">cd>name", func(c *xx.Context) error {
		charges := c.Value.(*Response).DomainCheckResponse.Charges
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
		dcr := &c.Value.(*Response).DomainCheckResponse
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

	// Scan frnic-2.0 extension for reserved/forbidden attributes
	path = "epp > response > extension > " + ExtFrnic20 + " ext > resData > chkData > domain > cd > name"
	scanResponse.MustHandleCharData(path, func(c *xx.Context) error {
		dcr := &c.Value.(*Response).DomainCheckResponse

		// Read the reserved attribute and store it as a charge category
		if c.AttrBool("", "reserved") {
			charge := DomainCharge{
				Domain:   string(c.CharData), // domain name from element text
				Category: "reserved",
			}
			dcr.Charges = append(dcr.Charges, charge)
		}

		// Optionally handle forbidden attribute
		if c.AttrBool("", "forbidden") {
			charge := DomainCharge{
				Domain:   string(c.CharData),
				Category: "forbidden",
			}
			dcr.Charges = append(dcr.Charges, charge)
		}

		return nil
	})
}
