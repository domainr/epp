package epp

import (
	"bytes"
	"encoding/xml"
	"net"
	"testing"

	"github.com/domainr/epp/ns"
	"github.com/nbio/st"
)

func TestHello(t *testing.T) {
	ls, err := newLocalServer()
	st.Assert(t, err, nil)
	defer ls.teardown()
	ls.buildup(func(ls *localServer, ln net.Listener) {
		conn, err := ls.Accept()
		st.Assert(t, err, nil)
		// Respond with greeting
		err = writeDataUnit(conn, []byte(testXMLGreeting))
		st.Assert(t, err, nil)
		// Respond with greeting for <hello>
		err = writeDataUnit(conn, []byte(testXMLGreeting))
		st.Assert(t, err, nil)
	})
	nc, err := net.Dial(ls.Listener.Addr().Network(), ls.Listener.Addr().String())
	st.Assert(t, err, nil)

	c, err := NewConn(nc)
	st.Assert(t, err, nil)
	err = c.Hello()
	st.Expect(t, err, nil)
	st.Expect(t, c.Greeting.ServerName, "Example EPP server epp.example.com")
}

func TestGreetingSupportsObject(t *testing.T) {
	g := Greeting{}
	st.Expect(t, g.SupportsObject(ns.Domain), false)
	st.Expect(t, g.SupportsObject(ns.Host), false)
	g.Objects = testObjects
	st.Expect(t, g.SupportsObject(ns.Domain), true)
	st.Expect(t, g.SupportsObject(ns.Host), true)
}

func TestGreetingSupportsExtension(t *testing.T) {
	g := Greeting{}
	st.Expect(t, g.SupportsExtension(ns.Charge), false)
	st.Expect(t, g.SupportsExtension(ns.IDN), false)
	g.Extensions = testExtensions
	st.Expect(t, g.SupportsExtension(ns.Charge), true)
	st.Expect(t, g.SupportsExtension(ns.IDN), true)
}

func TestScanGreeting(t *testing.T) {
	d := decoder(testXMLGreeting)
	var res Response
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, res.Greeting.ServerName, "Example EPP server epp.example.com")
	st.Expect(t, res.Greeting.Objects[0], "urn:ietf:params:xml:ns:obj1")
	st.Expect(t, res.Greeting.Objects[1], "urn:ietf:params:xml:ns:obj2")
	st.Expect(t, res.Greeting.Objects[2], "urn:ietf:params:xml:ns:obj3")
	st.Expect(t, res.Greeting.Extensions[0], "http://custom/obj1ext-1.0")
}

func BenchmarkScanGreeting(b *testing.B) {
	b.StopTimer()
	var buf bytes.Buffer
	d := xml.NewDecoder(&buf)
	saved := *d
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf.Reset()
		buf.WriteString(testXMLGreeting)
		deleteBufferRange(&buf, []byte(`<dcp>`), []byte(`</dcp>`))
		*d = saved
		b.StartTimer()
		var res Response
		scanResponse.Scan(d, &res)
	}
}

var testXMLGreeting = `<?xml version="1.0" encoding="utf-8"?>
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
</epp>`
