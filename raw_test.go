package epp

import (
	"net"
	"testing"

	"github.com/nbio/st"
)

func TestRaw(t *testing.T) {
	// Sample greeting response
	greeting := `<?xml version="1.0" encoding="UTF-8"?><epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>TestServer</svID><svcMenu><version>1.0</version></svcMenu></greeting></epp>`

	// Sample info request
	infoReq := []byte(`<?xml version="1.0" encoding="UTF-8"?><epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><info><domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:info></info></command></epp>`)

	// Sample info response
	infoRes := `<?xml version="1.0" encoding="UTF-8"?><epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result code="1000"><msg>Command completed successfully</msg></result><resData><domain:infData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:roid>EXAMPLE1-REP</domain:roid></domain:infData></resData></response></epp>`

	// Mock server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	st.Assert(t, err, nil)
	defer ln.Close()

	go func() {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		defer conn.Close()

		// Send greeting
		writeDataUnit(conn, []byte(greeting))

		// Read info request
		_, err = readDataUnitHeader(conn)
		if err != nil {
			return
		}
		body := make([]byte, 1024) // sufficiently large
		n, _ := conn.Read(body)
		_ = body[:n]

		// Send info response
		writeDataUnit(conn, []byte(infoRes))
	}()

	conn, err := net.Dial("tcp", ln.Addr().String())
	st.Assert(t, err, nil)
	defer conn.Close()

	eppConn, err := NewConn(conn)
	st.Assert(t, err, nil)

	res, err := eppConn.Raw(infoReq)
	st.Expect(t, err, nil)
	st.Expect(t, string(res), infoRes)
}
