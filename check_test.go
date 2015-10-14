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
	var buf bytes.Buffer
	err := encodeDomainCheck(&buf, []string{"hello.com", "foo.domains", "xn--ninja.net"}, false)
	st.Expect(t, err, nil)
	st.Expect(t, buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>hello.com</domain:name><domain:name>foo.domains</domain:name><domain:name>xn--ninja.net</domain:name></domain:check></check></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(buf.Bytes(), &v)
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

func TestScanCheckDomainResponseWithFee(t *testing.T) {
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
					<fee:name>good.space</fee:name>
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

func BenchmarkEncodeDomainCheck(b *testing.B) {
	var buf bytes.Buffer
	domains := []string{"hello.com"}
	for i := 0; i < b.N; i++ {
		encodeDomainCheck(&buf, domains, false)
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
