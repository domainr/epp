package epp

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

// DomainCheckOld represents the output of the EPP <domain:check> command.
type DomainCheckOld struct {
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
	req := Message{
		Command: &command{
			Check: &check{
				DomainCheck: &domainCheck{
					Domains: domains,
				},
			},
			TxnID: c.id(),
		},
	}
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
