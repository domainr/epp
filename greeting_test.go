package epp

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestHello(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	err = c.Hello()
	st.Expect(t, err, nil)
	st.Expect(t, c.Greeting.ServerName, "ISPAPI EPP Server") // FIXME: brittle external dependency
}

func TestDecodeGreeting(t *testing.T) {
	d := NewDecoder(bytes.NewBuffer(testXMLGreeting))
	var msg message
	err := d.DecodeMessage(&msg)
	st.Expect(t, err, nil)
	st.Reject(t, msg.Greeting, nil)
	st.Expect(t, msg.Greeting.ServerName, "Example EPP server epp.example.com")
	st.Expect(t, msg.Greeting.Objects[0], "urn:ietf:params:xml:ns:obj1")
	st.Expect(t, msg.Greeting.Objects[1], "urn:ietf:params:xml:ns:obj2")
	st.Expect(t, msg.Greeting.Objects[2], "urn:ietf:params:xml:ns:obj3")
	st.Expect(t, msg.Greeting.Extensions[0], "http://custom/obj1ext-1.0")
	logMarshal(t, &msg)

	d = NewDecoder(bytes.NewBuffer(testXMLGreeting))
	var g Greeting
	err = IgnoreEOF(decodeGreeting(&d, &g))
	st.Expect(t, err, nil)
	st.Expect(t, g.ServerName, "Example EPP server epp.example.com")
	st.Expect(t, g.Objects[0], "urn:ietf:params:xml:ns:obj1")
	st.Expect(t, g.Objects[1], "urn:ietf:params:xml:ns:obj2")
	st.Expect(t, g.Objects[2], "urn:ietf:params:xml:ns:obj3")
	st.Expect(t, g.Extensions[0], "http://custom/obj1ext-1.0")
}

func TestScanGreeting(t *testing.T) {
	d := xml.NewDecoder(bytes.NewBuffer(testXMLGreeting))
	var res response_
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)
	st.Expect(t, res.greeting.ServerName, "Example EPP server epp.example.com")
	st.Expect(t, res.greeting.Objects[0], "urn:ietf:params:xml:ns:obj1")
	st.Expect(t, res.greeting.Objects[1], "urn:ietf:params:xml:ns:obj2")
	st.Expect(t, res.greeting.Objects[2], "urn:ietf:params:xml:ns:obj3")
	st.Expect(t, res.greeting.Extensions[0], "http://custom/obj1ext-1.0")
}

func BenchmarkScanGreeting(b *testing.B) {
	b.StopTimer()
	var buf bytes.Buffer
	d := NewDecoder(&buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf.Reset()
		buf.Write(testXMLGreeting)
		deleteBufferRange(&buf, []byte(`<dcp>`), []byte(`</dcp>`))
		d.Reset()
		b.StartTimer()
		var res response_
		scanResponse.Scan(&d.decoder, &res)
	}
}

func BenchmarkDecodeGreeting(b *testing.B) {
	b.StopTimer()
	var buf bytes.Buffer
	d := NewDecoder(&buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf.Reset()
		buf.Write(testXMLGreeting)
		deleteBufferRange(&buf, []byte(`<dcp>`), []byte(`</dcp>`))
		d.Reset()
		b.StartTimer()
		var g Greeting
		decodeGreeting(&d, &g)
	}
}

func BenchmarkDecoderDecodeGreeting(b *testing.B) {
	b.StopTimer()
	var buf bytes.Buffer
	d := NewDecoder(&buf)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		buf.Reset()
		buf.Write(testXMLGreeting)
		deleteBufferRange(&buf, []byte(`<dcp>`), []byte(`</dcp>`))
		d.Reset()
		b.StartTimer()
		var msg message
		d.decoder.Decode(&msg)
	}
}

var testXMLGreeting = []byte(`<?xml version="1.0" encoding="utf-8"?>
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
