package epp

import (
	"encoding/xml"
	"fmt"
	"io"
)

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
			TxnID: c.id(),
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

func decodeCheckDomainResponse(d *decoder) (*domainCheckData, error) {
	d.reset()
	data := &domainCheckData{}
	for {
		t, err := d.Token()
		if err != nil && err != io.EOF {
			return nil, err
		}
		if t == nil {
			break
		}
		switch node := t.(type) {
		case xml.StartElement:
			fmt.Printf("StartElement: %s %s\n", node.Name.Space, node.Name.Local)
		case xml.EndElement:
			fmt.Printf("EndElement: %s %s\n", node.Name.Space, node.Name.Local)
		case xml.CharData:
			fmt.Printf("CharData: %s\n", string(node))
		}
		fmt.Printf("Stack: %+v\n", d.stack)
	}
	return data, nil
}
