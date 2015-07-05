package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestConnLogin(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
}

func TestEncodeLogin(t *testing.T) {
	var buf bytes.Buffer
	err := encodeLogin(&buf, "jane", "battery", "", "1.0", "en", nil, nil)
	st.Expect(t, err, nil)
	st.Expect(t, buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(buf.Bytes(), &v)
	st.Expect(t, err, nil)
}

func TestEncodeLoginChangePassword(t *testing.T) {
	var buf bytes.Buffer
	err := encodeLogin(&buf, "jane", "battery", "horse", "1.0", "en", nil, nil)
	st.Expect(t, err, nil)
	st.Expect(t, buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><newPW>horse</newPW><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(buf.Bytes(), &v)
	st.Expect(t, err, nil)
}

var (
	testObjects = []string{
		"urn:ietf:params:xml:ns:domain-1.0",
		"urn:ietf:params:xml:ns:host-1.0",
		"urn:ietf:params:xml:ns:contact-1.0",
		"http://www.unitedtld.com/epp/finance-1.0",
	}
	testExtensions = []string{
		"urn:ietf:params:xml:ns:secDNS-1.1",
		"urn:ietf:params:xml:ns:rgp-1.0",
		"urn:ietf:params:xml:ns:launch-1.0",
		"urn:ietf:params:xml:ns:idn-1.0",
		"http://www.unitedtld.com/epp/charge-1.0",
	}
)

func BenchmarkEncodeLogin(b *testing.B) {
	var buf bytes.Buffer
	for i := 0; i < b.N; i++ {
		encodeLogin(&buf, "jane", "battery", "horse", "1.0", "en", testObjects, testExtensions)
	}
}
