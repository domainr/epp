package epp

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/nbio/st"
)

func logMarshal(t *testing.T, msg *Message) {
	x, err := Marshal(&msg)
	st.Expect(t, err, nil)
	t.Logf("<!-- MARSHALED -->\n%s\n", string(x))
}

func TestUnmarshalGreeting(t *testing.T) {
	x := []byte(`<?xml version="1.0" encoding="utf-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<greeting>
		<svID>Example EPP server epp.example.com</svID>
		<svDate>2000-06-08T22:00:00.0Z</svDate>
		<svcMenu>
			<version>1.0</version>
			<lang>en</lang>
			<lang>fr</lang>
			<objURI>urn:ietf:params:xml:ns:obj1</objURI>
			<objURI>urn:ietf:params:xml:ns:obj2</objURI>
			<objURI>urn:ietf:params:xml:ns:obj3</objURI>
			<svcExtension>
				<extURI>http://custom/obj1ext-1.0</extURI>
			</svcExtension>
		</svcMenu>
		<dcp>
			<access><all/></access>
			<statement>
				<purpose><admin/><prov/></purpose>
				<recipient><ours/><public/></recipient>
				<retention><stated/></retention>
			</statement>
		</dcp>
	</greeting>
</epp>`)

	var msg Message
	err := xml.Unmarshal(x, &msg)
	st.Expect(t, err, nil)
	st.Reject(t, msg.Greeting, nil)
	st.Expect(t, msg.Greeting.ServerName, "Example EPP server epp.example.com")
	tt, _ := time.Parse(time.RFC3339, "2000-06-08T22:00:00.0Z")
	st.Expect(t, msg.Greeting.ServerTime, Time{tt})
	st.Expect(t, msg.Greeting.ServiceMenu.Objects[0], "urn:ietf:params:xml:ns:obj1")
	st.Expect(t, msg.Greeting.ServiceMenu.Objects[1], "urn:ietf:params:xml:ns:obj2")
	st.Expect(t, msg.Greeting.ServiceMenu.Objects[2], "urn:ietf:params:xml:ns:obj3")
	st.Expect(t, msg.Greeting.ServiceMenu.Extensions[0], "http://custom/obj1ext-1.0")
	st.Expect(t, msg.Greeting.DCP.Access.None, (*struct{})(nil))
	st.Reject(t, msg.Greeting.DCP.Access.All, (*struct{})(nil))
	st.Reject(t, msg.Greeting.DCP.Statement[0].Purpose.Admin, (*struct{})(nil))
	st.Expect(t, msg.Greeting.DCP.Statement[0].Purpose.Other, (*struct{})(nil))
	logMarshal(t, &msg)
}

func TestUnmarshalCheckDomainResponse(t *testing.T) {
	x := []byte(`<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<response>
		<result code="1000">
			<msg>Command completed successfully</msg>
		</result>
		<resData>
			<domain:chkData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
				<domain:cd>
					<domain:name avail="1">good.memorial</domain:name>
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
</epp>`)

	var msg Message
	err := Unmarshal(x, &msg)
	st.Expect(t, err, nil)
	st.Reject(t, msg.Response, nil)
	logMarshal(t, &msg)
}
