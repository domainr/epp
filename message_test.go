package epp

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestUnmarshal(t *testing.T) {
	x := []byte(`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
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
}
