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

// <resData>
//  <domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd">
//   <domain:cd>
//    <domain:name avail="0">ace.pizza</domain:name>
//    <domain:reason>Premium Domain Name</domain:reason>
//   </domain:cd>
//  </domain:chkData>
// </resData>

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (*DomainCheck, error) {
	req := message{
		Command: &command{
			Check: &check{
				DomainCheck: &domainCheck{
					Domains: domains,
				},
			},
		},
	}
	err := c.writeMessage(&req)
	if err != nil {
		return nil, err
	}
	msg := message{}
	err = c.readMessage(&msg)
	if err != nil {
		return nil, err
	}
	res := msg.Response
	if res == nil || res.ResponseData.DomainCheckData == nil {
		return nil, ErrResponseMalformed
	}
	dc := DomainCheck(*res.ResponseData.DomainCheckData)
	return &dc, nil
}

// DomainCheck exists for backwards compatibility.
// FIXME: remove/improve this.
type DomainCheck domainCheckData

type DomainCheckResponse struct {
	Checks  []DomainCheck_
	Charges []DomainCharge
}

type DomainCheck_ struct {
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
	scanResponse.MustHandleStartElement("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData", func(c *Context) error {
		c.Value.(*response_).DomainCheckResponse = DomainCheckResponse{}
		return nil
	})
	scanResponse.MustHandleStartElement("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd", func(c *Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Checks = append(dcd.Checks, DomainCheck_{})
		return nil
	})
	scanResponse.MustHandleCharData("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd > name", func(c *Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Domain = string(c.CharData)
		a := c.Attr("", "avail")
		check.Available = (a == "1" || a == "true")
		return nil
	})
	scanResponse.MustHandleCharData("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd > reason", func(c *Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Reason = string(c.CharData)
		return nil
	})
}
