package epp

import (
	"encoding/xml"
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
	Checks  []DomainCheck
	Charges []DomainCharge
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
}

func decodeDomainCheckResponse(d *Decoder) ([]DomainCheck_, error) {
	var r Result
	var dcs []DomainCheck_
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return dcs, err
		}
		switch node := t.(type) {
		case xml.StartElement:
			switch {
			case d.AtPath("epp", "response", "result"):
				err := d.decodeResult(&r)
				if err != nil {
					return dcs, err
				}

			case d.AtPath("resData", "chkData", "cd"):
				_, err := decodeDomainCheckData(d)
				if err != nil {
					return dcs, err
				}
			}

		case xml.CharData:
			if string(node) != "" {

			}
		}
		// fmt.Printf("Stack: %+v\n", d.Stack)
	}
	return dcs, nil
}

type DomainCheck_ struct {
	Domain    string
	Reason    string
	Available bool
}

type domainCheckData_ struct {
	domain    string
	available Bool
	reason    string
}

func decodeDomainCheckData(d *Decoder) (domainCheckData_, error) {
	var data domainCheckData_
outer:
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return data, err
		}
		switch node := t.(type) {
		case xml.EndElement:
			if node.Name.Local == "cd" {
				break outer
			}
		case xml.CharData:
			e := d.Element(-1)
			switch e.Name.Local {
			case "name":
				data.domain = string(node)
				for _, a := range e.Attr {
					if a.Name.Local == "avail" {
						data.available.UnmarshalXMLAttr(&a)
					}
				}
			case "reason":
				data.reason = string(node)
			}
		}
	}
	return data, nil
}
