package epp

// domainCheckRequest represents an EPP <domain:check> command.
// https://tools.ietf.org/html/rfc5730#section-2.9.2.1
type domainCheckRequest struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Check   struct {
		XMLNamespace DomainNamespace `xml:"xmlns:domain,attr"`
		Domains      []string        `xml:"domain:check>domain:name"`
	} `xml:"command>check"`
	TxnID string `xml:"command>clTRID"`
}

// <epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
//   <command>
//     <check>
//     <domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
//       <domain:name>dmnr-test-1234.com</domain:name>
//     </domain:check>
//     </check>
//     <clTRID>ABC-12345</clTRID>
//   </command>
// </epp>

// The DomainNamespace type exists solely to emit an XML attribute.
type DomainNamespace struct{}

// MarshalText returns a byte slice for the xmlns:xsi attribute.
func (n DomainNamespace) MarshalText() (text []byte, err error) {
	return []byte("urn:ietf:params:xml:ns:domain-1.0"), nil
}

// DomainCheck represents the output of the EPP <domain:check> command.
type DomainCheck struct {
	Results []struct {
		Domain struct {
			Domain      string `xml:",chardata"`
			IsAvailable bool   `xml:"avail,attr"`
		} `xml:"name"`
		Reason string `xml:"reason"`
	} `xml:"cd"`
}

// <resData>
//  <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd">
//   <domain:cd>
//    <domain:name avail="0">ace.pizza</domain:name>
//    <domain:reason>Premium Domain Name</domain:reason>
//   </domain:cd>
//  </domain:chkData>
// </resData>

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (dc *DomainCheck, err error) {
	req := domainCheckRequest{TxnID: c.id()}
	req.Check.Domains = domains
	err = c.WriteRequest(&req)
	if err != nil {
		return
	}
	res := Response{}
	err = c.ReadResponse(&res)
	if err != nil {
		return
	}
	if res.DomainCheck == nil {
		return nil, ErrMalformedResponse
	}
	return res.DomainCheck, nil
}
