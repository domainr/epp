package epp

import (
	"encoding/xml"
	"testing"
)

var info1 string = `
<?xml version="1.0" encoding="utf-8"?>
<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
<response>
	<result code="1000">
		<msg>Command completed successfully</msg>
	</result>
	<resData>
		<domain:infData xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:domain-1.0 domain-1.0.xsd" xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
			<domain:name>example.local</domain:name>
			<domain:roid>fefefefe456ae747badf00d-ABCDEFG</domain:roid>
			<domain:status s="ok" />
			<domain:clID>REG-1234</domain:clID>
			<domain:crDate>1996-05-19T05:00:00.0Z</domain:crDate>
			<domain:upDate>2020-01-20T01:28:49.406Z</domain:upDate>
			<domain:exDate>2025-05-19T05:00:00.0Z</domain:exDate>
			<domain:trDate>2005-11-01T20:12:08.0Z</domain:trDate>
		</domain:infData>
	</resData>
	<extension>
		<secDNS:infData xmlns:secDNS="urn:ietf:params:xml:ns:secDNS-1.1">
		<secDNS:dsData>
		<secDNS:keyTag>12345</secDNS:keyTag>
		<secDNS:alg>5</secDNS:alg>
		<secDNS:digestType>1</secDNS:digestType>
		<secDNS:digest>BADF00DCB65CFEFEFEFEB3640062E8CBADF00D</secDNS:digest>
		</secDNS:dsData>
		<secDNS:dsData>
		<secDNS:keyTag>98765</secDNS:keyTag>
		<secDNS:alg>5</secDNS:alg>
		<secDNS:digestType>2</secDNS:digestType>
		<secDNS:digest>BADF00D7F1D07B231344FEFEFEFEE7519ADDAE180E20B1B1EC52E7FBADF00D</secDNS:digest>
		</secDNS:dsData>
		</secDNS:infData>
	</extension>
	<trID>
		<svTRID>badf00de4cbadf00d3e32c6fefefefe</svTRID>
	</trID>
</response>
</epp>`

func TestInfoResponse(t *testing.T) {
	info := DomainInfoResponse{}
	if err := xml.Unmarshal([]byte(info1), &info); err != nil {
		t.Fatal("Failed to unmarshal info response:", err)
	}

	if info.Name != "example.local" {
		t.Error("Expected example.local, got:", info.Name)
	}

	if info.ROID != "fefefefe456ae747badf00d-ABCDEFG" {
		t.Error("Unexpected roid value:", info.ROID)
	}

	if info.CLID != "REG-1234" {
		t.Error("Expected CLID:REG-1234, got:", info.CLID)
	}
}
