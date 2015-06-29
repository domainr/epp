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
