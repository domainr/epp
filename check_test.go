package epp

import (
	"testing"

	"github.com/nbio/st"
)

func TestConnCheck(t *testing.T) {
	c := testLogin(t)

	dc, err := c.CheckDomain("google.com")
	st.Expect(t, err, nil)
	st.Reject(t, dc, nil)
	st.Expect(t, len(dc.Results), 1)
	st.Expect(t, dc.Results[0].Domain.Domain, "google.com")
	st.Expect(t, dc.Results[0].Domain.IsAvailable, false)

	dc, err = c.CheckDomain("dmnr-test-x759824vim-i2.com")
	st.Expect(t, err, nil)
	st.Reject(t, dc, nil)
	st.Expect(t, len(dc.Results), 1)
	st.Expect(t, dc.Results[0].Domain.Domain, "dmnr-test-x759824vim-i2.com")
	st.Expect(t, dc.Results[0].Domain.IsAvailable, true)

	dc, err = c.CheckDomain("--dmnr-test--.com")
	st.Reject(t, err, nil)
	st.Expect(t, dc, (*DomainCheck)(nil))
}

func TestScanCheckDomainResponse(t *testing.T) {
	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(testXMLDomainCheckResponse)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "good.memorial")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "premium name")
}

var testXMLDomainCheckResponse = `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">good.memorial</domain:name>
					<domain:reason>premium name</domain:reason>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<charge:chkData xmlns:charge="http://www.unitedtld.com/epp/charge-1.0">
				<charge:cd>
					<charge:name>good.memorial</charge:name>
					<charge:set>
						<charge:category name="BBB+">premium</charge:category>
						<charge:type>price</charge:type>
						<charge:amount command="create">100.00</charge:amount>
						<charge:amount command="renew">100.00</charge:amount>
						<charge:amount command="transfer">100.00</charge:amount>
						<charge:amount command="update" name="restore">50.00</charge:amount>
					</charge:set>
				</charge:cd>
			</charge:chkData>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`
