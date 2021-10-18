package epp

import (
	"bytes"
	"encoding/xml"
	"time"

	"github.com/nbio/xx"
)

// DomainInfo retrieves info for a domain.
// https://tools.ietf.org/html/rfc5731#section-3.1.2
func (c *Conn) DomainInfo(domain string, extData map[string]string) (*DomainInfoResponse, error) {
	x, err := encodeDomainInfo(&c.Greeting, domain, extData)
	if err != nil {
		return nil, err
	}
	err = c.writeDataUnit(x)
	if err != nil {
		return nil, err
	}
	var res Response
	err = c.readResponse(&res)
	if err != nil {
		return nil, err
	}
	return &res.DomainInfoResponse, nil
}

func encodeDomainInfo(greeting *Greeting, domain string, extData map[string]string) ([]byte, error) {
	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<info><domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name hosts="none">`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name></domain:info></info>`)

	supportsNamestore := extData["namestoreExt:subProduct"] != "" && greeting.SupportsExtension(ExtNamestore)
	hasExtension := supportsNamestore

	if hasExtension {
		buf.WriteString(`<extension>`)
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
		buf.WriteString(`</extension>`)
	}

	buf.WriteString(xmlCommandSuffix)

	return buf.Bytes(), nil
}

// DomainInfoResponse represents an EPP response for a domain info request.
// https://tools.ietf.org/html/rfc5731#section-3.1.2
type DomainInfoResponse struct {
	Domain string    // <domain:name>
	ID     string    // <domain:roid>
	ClID   string    // <domain:clID>
	UpID   string    // <domain:upID>
	CrDate time.Time // <domain:crDate>
	ExDate time.Time // <domain:exDate>
	UpDate time.Time // <domain:upDate>
	TrDate time.Time // <domain:trDate>
	Status []string  // <domain:status>
}

func init() {
	// Default EPP check data
	path := "epp > response > resData > " + ObjDomain + " infData"
	scanResponse.MustHandleCharData(path+">name", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		dir.Domain = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">roid", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		dir.ID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">clID", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		dir.ClID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">upID", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		dir.UpID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">crDate", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		var err error
		dir.CrDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">exDate", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		var err error
		dir.ExDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">upDate", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		var err error
		dir.UpDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">trDate", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		var err error
		dir.TrDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleStartElement(path+">status", func(c *xx.Context) error {
		dir := &c.Value.(*Response).DomainInfoResponse
		dir.Status = append(dir.Status, c.Attr("", "s"))
		return nil
	})
}

//lint:ignore U1000 keeping around for reference
func encodeVerisignDomainInfo(buf *bytes.Buffer, domain string) error {
	buf.Reset()
	buf.WriteString(xmlCommandPrefix)
	buf.WriteString(`<info><domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name hosts="none">`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name></domain:info></info>`)
	buf.WriteString(`<extension>`)
	buf.WriteString(`<namestoreExt:namestoreExt xmlns:namestoreExt="http://www.verisign-grs.com/epp/namestoreExt-1.1">`)
	buf.WriteString(`<namestoreExt:subProduct>`)
	buf.WriteString(`com`)
	buf.WriteString(`</namestoreExt:subProduct>`)
	buf.WriteString(`</namestoreExt:namestoreExt>`)
	buf.WriteString(`</extension>`)
	buf.WriteString(xmlCommandSuffix)
	return nil
}

//lint:ignore U1000 keeping around for reference
func encodeVerisignContactInfo(buf *bytes.Buffer, contact string) error {
	buf.Reset()
	buf.WriteString(xmlCommandPrefix)
	buf.WriteString(`<info><contact:info xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>`)
	xml.EscapeText(buf, []byte(contact))
	buf.WriteString(`</contact:id></contact:info></info>`)
	buf.WriteString(`<extension>`)
	buf.WriteString(`<namestoreExt:namestoreExt xmlns:namestoreExt="http://www.verisign-grs.com/epp/namestoreExt-1.1">`)
	buf.WriteString(`<namestoreExt:subProduct>`)
	buf.WriteString(`com`)
	buf.WriteString(`</namestoreExt:subProduct>`)
	buf.WriteString(`</namestoreExt:namestoreExt>`)
	buf.WriteString(`</extension>`)
	buf.WriteString(xmlCommandSuffix)
	return nil
}

/*
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <info>
      <domain:info
       xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name hosts="all">example.com</domain:name>
      </domain:info>
    </info>
		<extension>
			<namestoreExt:namestoreExt xmlns:namestoreExt="http://www.verisign-grs.com/epp/namestoreExt-1.1">
			   <namestoreExt:subProduct>TLD</namestoreExt:subProduct>
			</namestoreExt:namestoreExt>
    </extension>
    <clTRID>ABC-12345</clTRID>
  </command>
</epp>


<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <command>
    <info>
      <contact:info
       xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
        <contact:id>sh8013</contact:id>
        <contact:authInfo>
          <contact:pw>2fooBAR</contact:pw>
        </contact:authInfo>
      </contact:info>
    </info>
    <clTRID>ABC-12345</clTRID>
  </command>
</epp>
*/
