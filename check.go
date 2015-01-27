package epp

// DomainCheck represents an EPP <domain:check> command.
// https://tools.ietf.org/html/rfc5730#section-2.9.2.1
type DomainCheckMessage struct {
	MessageNamespace
	Check struct {
		DomainNamespace DomainNamespace `xml:"xmlns:domain,attr"`
		Domains         []string        `xml:"domain:check>domain:name"`
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

// DomainCheckResponse represents the output of the EPP <domain:check> command.
type DomainCheckResponse struct {
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
func (c *Conn) CheckDomain(domains ...string) (dcr *DomainCheckResponse, err error) {
	msg := DomainCheckMessage{TxnID: c.id()}
	msg.Check.Domains = domains
	err = c.WriteMessage(&msg)
	if err != nil {
		return
	}
	rmsg := ResponseMessage{}
	err = c.ReadResponse(&rmsg)
	if err != nil {
		return
	}
	if rmsg.DomainCheck == nil {
		return nil, ErrMalformedResponse
	}
	return rmsg.DomainCheck, nil
}
