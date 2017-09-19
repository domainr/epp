package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func decoder(s string) *xml.Decoder {
	return xml.NewDecoder(bytes.NewBufferString(s))
}

func TestConnCheck(t *testing.T) {
	c := testLogin(t)

	dcr, err := c.CheckDomain("google.com")
	st.Expect(t, err, nil)
	st.Reject(t, dcr, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "google.com")
	st.Expect(t, dcr.Checks[0].Available, false)

	dcr, err = c.CheckDomain("dmnr-test-x759824vim-i2.com")
	st.Expect(t, err, nil)
	st.Reject(t, dcr, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "dmnr-test-x759824vim-i2.com")
	st.Expect(t, dcr.Checks[0].Available, true)

	dcr, err = c.CheckDomain("--dmnr-test--.com")
	st.Reject(t, err, nil)
	st.Expect(t, dcr, (*DomainCheckResponse)(nil))
}

func TestEncodeDomainCheck(t *testing.T) {
	con := &Conn{}
	err := con.encodeDomainCheck([]string{"hello.com", "foo.domains", "xn--ninja.net"}, nil)
	st.Expect(t, err, nil)
	st.Expect(t, con.buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>hello.com</domain:name><domain:name>foo.domains</domain:name><domain:name>xn--ninja.net</domain:name></domain:check></check></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(con.buf.Bytes(), &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainCheckLaunchPhase(t *testing.T) {
	con := &Conn{}
	con.Greeting.Extensions = []string{ExtLaunch}
	err := con.encodeDomainCheck([]string{"hello.com", "foo.domains", "xn--ninja.net"}, map[string]string{"launch:phase": "claims"})
	st.Expect(t, err, nil)
	st.Expect(t, con.buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>hello.com</domain:name><domain:name>foo.domains</domain:name><domain:name>xn--ninja.net</domain:name></domain:check></check><extension><launch:check xmlns:launch="urn:ietf:params:xml:ns:launch-1.0" type="avail"><launch:phase>claims</launch:phase></launch:check></extension></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(con.buf.Bytes(), &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainCheckNeulevelUnspec(t *testing.T) {
	con := &Conn{}
	con.Greeting.Extensions = []string{ExtNeulevel}
	err := con.encodeDomainCheck([]string{"hello.com", "foo.domains", "xn--ninja.net"}, map[string]string{"neulevel:unspec": "FeeCheck=Y"})
	st.Expect(t, err, nil)
	st.Expect(t, con.buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>hello.com</domain:name><domain:name>foo.domains</domain:name><domain:name>xn--ninja.net</domain:name></domain:check></check><extension><neulevel:extension xmlns:neulevel="urn:ietf:params:xml:ns:neulevel-1.0"><neulevel:unspec>FeeCheck=Y</neulevel:unspec></neulevel:extension></extension></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(con.buf.Bytes(), &v)
	st.Expect(t, err, nil)
}

func TestScanCheckDomainResponseWithCharge(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
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

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "good.memorial")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "premium name")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "good.memorial")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "BBB+")
}
func TestScanCheckDomainResponseWithMultipleChargeSets(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
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
					<charge:set>
						<charge:category>earlyAccess</charge:category>
						<charge:type>fee</charge:type>
						<charge:amount command="create">2500.00</charge:amount>
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

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "good.memorial")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "premium name")
	st.Expect(t, len(dcr.Charges), 2)
	st.Expect(t, dcr.Charges[0].Domain, "good.memorial")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "BBB+")
	st.Expect(t, dcr.Charges[1].Domain, "good.memorial")
	st.Expect(t, dcr.Charges[1].Category, "earlyAccess")
	st.Expect(t, dcr.Charges[1].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee05(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">good.space</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.5">
				<fee:cd>
					<fee:name premium="true">good.space</fee:name>
					<fee:currency>USD</fee:currency>
					<fee:command>create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee description="Premium Registration Fee" grace-period="P5D" refundable="1">100.00</fee:fee>
					<fee:class>premium</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "good.space")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "good.space")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "Premium Registration Fee")
}

func TestScanCheckDomainResponseWithFee06(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">good.space</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.6">
				<fee:cd>
					<fee:name>good.space</fee:name>
					<fee:currency>USD</fee:currency>
					<fee:command>create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee>100.00</fee:fee>
					<fee:class>SPACE Tier 1</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "good.space")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "good.space")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee07(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">austin.green</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee='urn:ietf:params:xml:ns:fee-0.7' xsi:schemaLocation='urn:ietf:params:xml:ns:fee-0.7 fee-0.7.xsd'>
				<fee:cd>
					<fee:name>austin.green</fee:name>
					<fee:currency>USD</fee:currency>
					<fee:command>create</fee:command>
					<fee:period unit='y'>1</fee:period>
					<fee:fee description='Registration fee' refundable='1' grace-period='P5D'>3500.00</fee:fee>
					<fee:class>premium</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "austin.green")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "austin.green")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee08(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">crrc.yln</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.8">
				<fee:cd>
					<fee:name>crrc.yln</fee:name>
					<fee:currency>CNY</fee:currency>
					<fee:command>create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee applied="delayed" description="Registration Fee" grace-period="P5D" refundable="1">100.00</fee:fee>
					<fee:class>premium</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>testnn-domain-check-f193d63b-1ab7-43bc-bc9d-4e835fb0fece</clTRID>
			<svTRID>SERVER-4aafbfa9-cd31-4e89-b585-25a753d3c69a</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "crrc.yln")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "crrc.yln")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee09(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">example.com</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.9">
				<fee:cd>
					<fee:objID>example.com</fee:objID>
					<fee:currency>USD</fee:currency>
					<fee:command phase="sunrise">create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee description="Application Fee" refundable="0">5.00</fee:fee>
					<fee:fee description="Registration Fee" refundable="1" grace-period="P5D">5.00</fee:fee>
					<fee:class>premium-tier1</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>ABC-12345</clTRID>
			<svTRID>54322-XYZ</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "example.com")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "example.com")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee11(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">example.com</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.11">
				<fee:cd>
					<fee:objID>example.com</fee:objID>
					<fee:currency>USD</fee:currency>
					<fee:command phase="sunrise">create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee refundable="true" grace-period="P5D">50.00</fee:fee>
					<fee:class>premium</fee:class>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<clTRID>ABC-12345</clTRID>
			<svTRID>54322-XYZ</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "example.com")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "example.com")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee21(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<chkData xmlns="urn:ietf:params:xml:ns:domain-1.0">
				<cd>
					<name avail="true">example.com</name>
				</cd>
			</chkData>
		</resData>
		<extension>
			<chkData xmlns="urn:ietf:params:xml:ns:fee-0.21">
				<currency>EUR</currency>
				<cd>
					<objID>example.com</objID>
					<command name="create" phase="open">
						<period unit="y">1</period>
						<fee applied="immediate" description="domain creation in phase &#39;open&#39;" grace-period="P5D" refundable="true">25.00</fee>
					</command>
					<command name="create" phase="custom" subphase="open-50">
						<reason>the requested launch phase is not suitable for the domain</reason>
					</command>
				</cd>
			</chkData>
		</extension>
		<trID>
			<svTRID>1501792511080-81912</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "example.com")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "example.com")
	st.Expect(t, dcr.Charges[0].Category, "open")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponseWithFee21Premium(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<chkData xmlns="urn:ietf:params:xml:ns:domain-1.0">
				<cd>
					<name avail="false">example.com</name>
					<reason>not registrable in this phase</reason>
				</cd>
			</chkData>
		</resData>
		<extension>
			<chkData xmlns="urn:ietf:params:xml:ns:fee-0.21">
				<currency>EUR</currency>
				<cd>
					<objID>example.com</objID>
					<command name="create" phase="open">
						<reason>the requested launch phase is not suitable for the domain</reason>
					</command>
					<command name="create" phase="custom" subphase="open-1000">
						<period unit="y">1</period>
						<fee applied="immediate" description="domain creation in phase &#39;open-1000&#39;" grace-period="P5D" refundable="true">800.00</fee>
					</command>
				</cd>
			</chkData>
		</extension>
		<trID>
			<svTRID>1501792511080-81912</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "example.com")
	st.Expect(t, dcr.Checks[0].Available, false)
	st.Expect(t, dcr.Checks[0].Reason, "not registrable in this phase")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "example.com")
	st.Expect(t, dcr.Charges[0].Category, "custom")
	st.Expect(t, dcr.Charges[0].CategoryName, "open-1000")
}

func TestScanCheckDomainResponseWithPremiumAttribute(t *testing.T) {
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
			<fee:chkData xmlns:fee="urn:ietf:params:xml:ns:fee-0.5">
				<fee:cd>
					<fee:name premium="true">zero.work</fee:name>
					<fee:currency>USD</fee:currency>
					<fee:command>create</fee:command>
					<fee:period unit="y">1</fee:period>
					<fee:fee description="Registration Fee">500.000</fee:fee>
				</fee:cd>
			</fee:chkData>
		</extension>
		<trID>
			<svTRID>14470834306141</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
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

func TestScanCheckDomainResponseNeulevelExtension(t *testing.T) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns="urn:ietf:params:xml:ns:domain-1.0" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd">
				<domain:cd>
					<domain:name avail="1">420.earth</domain:name>
				</domain:cd>
			</domain:chkData>
		</resData>
		<extension>
			<neulevel:extension xmlns="urn:ietf:params:xml:ns:neulevel-1.0" xmlns:neulevel="urn:ietf:params:xml:ns:neulevel-1.0" xsi:schemaLocation="urn:ietf:params:xml:ns:neulevel-1.0 neulevel-1.0.xsd">
				<neulevel:unspec>TierName=EARTH_Tier3 AnnualTierPrice=120.00</neulevel:unspec>
			</neulevel:extension>
		</extension>
		<trID>
			<clTRID>0000000000000002</clTRID>
			<svTRID>83fa5767-5624-4be5-9e54-0b3a52f9de5b:1</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 1)
	st.Expect(t, dcr.Checks[0].Domain, "420.earth")
	st.Expect(t, dcr.Checks[0].Available, true)
	st.Expect(t, dcr.Checks[0].Reason, "")
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "420.earth")
	st.Expect(t, dcr.Charges[0].Category, "EARTH_Tier3")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func TestScanCheckDomainResponsePriceExtension(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg lang="en">Command completed successfully</msg>
		</result>
		<extension>
			<chkData xmlns="urn:ar:params:xml:ns:price-1.1">
				<cd>
					<name premium="1">foundations.build</name>
					<period unit="y">1</period>
					<createPrice>1500</createPrice>
					<renewPrice>1500</renewPrice>
					<restorePrice>40</restorePrice>
					<transferPrice>1500</transferPrice>
				</cd>
			</chkData>
		</extension>
		<trID>
			<svTRID>aaa39bf9-12dd-4810-bdb2-98f629cfbbbb</svTRID>
		</trID>
	</response>
</epp>`

	var res response_
	dcr := &res.DomainCheckResponse

	d := decoder(x)
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, len(dcr.Checks), 0)
	st.Expect(t, len(dcr.Charges), 1)
	st.Expect(t, dcr.Charges[0].Domain, "foundations.build")
	st.Expect(t, dcr.Charges[0].Category, "premium")
	st.Expect(t, dcr.Charges[0].CategoryName, "")
}

func BenchmarkEncodeDomainCheck(b *testing.B) {
	con := &Conn{}
	domains := []string{"hello.com"}
	for i := 0; i < b.N; i++ {
		con.encodeDomainCheck(domains, nil)
	}
}

func BenchmarkScanDomainCheckResponse(b *testing.B) {
	x := `<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
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

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		d := decoder(x)
		b.StartTimer()
		var res response_
		scanResponse.Scan(d, &res)
	}
}
