package epp

import (
	"encoding/xml"
	"net"
	"testing"

	"github.com/nbio/st"
)

func TestConnLogin(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
}

func TestWriteLogin(t *testing.T) {
	c := newConn(&net.IPConn{})
	err := c.writeLogin("jane", "battery", "")
	st.Expect(t, err, nil)
	st.Expect(t, c.buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login><clTRID>0000000000000001</clTRID></command></epp>`)
	var msg message
	err = xml.Unmarshal(c.buf.Bytes(), &msg)
	st.Expect(t, err, nil)
}

func TestWriteLoginChangePassword(t *testing.T) {
	c := newConn(&net.IPConn{})
	err := c.writeLogin("jane", "battery", "horse")
	st.Expect(t, err, nil)
	st.Expect(t, c.buf.String(), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><newPW>horse</newPW><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login><clTRID>0000000000000001</clTRID></command></epp>`)
	var msg message
	err = xml.Unmarshal(c.buf.Bytes(), &msg)
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

func BenchmarkMarshalLogin(b *testing.B) {
	b.StopTimer()
	msg := message{
		Command: &command{
			Login: &login{
				User:        "jane",
				Password:    "battery",
				NewPassword: "horse",
				Version:     "1.0",
				Language:    "en",
				Objects:     testObjects,
				Extensions:  testExtensions,
			},
			TxnID: "0001",
		},
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		xml.Marshal(&msg)
	}
}

func BenchmarkWriteLogin(b *testing.B) {
	b.StopTimer()
	c := newConn(&net.IPConn{})
	c.Greeting.Objects = testObjects
	c.Greeting.Extensions = testExtensions
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		c.writeLogin("jane", "battery", "horse")
	}
}
