package epp

import (
	"testing"

	"github.com/nbio/st"
)

func TestScanCheckDomainResponseWithPremiumAttributeWithFee10(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">zero.work</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:epp:fee-1.0">
				<fee:currency>USD</fee:currency>
				<fee:cd avail="1">
					<fee:name premium="true">zero.work</fee:name>
					<fee:objID>zero.work</fee:objID>
					<fee:class>premium</fee:class>
					<fee:command name="create">
						<fee:period unit="y">1</fee:period>
						<fee:fee description="Registration Fee" refundable="1" grace-period="P5D">500.000</fee:fee>
					</fee:command>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<svTRID>14470834306141</svTRID>
		</trID>
	</response>
</epp>`

	var res Response
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "zero.work")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "zero.work")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "Registration Fee")
}
