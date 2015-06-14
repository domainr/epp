package epp

// domainCheckRequest represents an EPP <domain:check> command.
// https://tools.ietf.org/html/rfc5730#section-2.9.2.1
type domainCheckRequest struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Check   struct {
		DomainNS domainNS `xml:"xmlns:domain,attr"`
		Domains  []string `xml:"domain:check>domain:name"`
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
func (c *Conn) CheckDomain(domains ...string) (dc *DomainCheckData, err error) {
	req := domainCheckRequest{TxnID: c.id()}
	req.Check.Domains = domains
	err = c.WriteMessage(&req)
	if err != nil {
		return
	}
	msg := Message{}
	err = c.ReadMessage(&msg)
	if err != nil {
		return
	}
	res := msg.Response
	if res == nil || res.ResponseData == nil || res.ResponseData.DomainCheckData == nil {
		return nil, ErrResponseMalformed
	}
	return res.ResponseData.DomainCheckData, nil
}
