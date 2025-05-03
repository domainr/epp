package epp

import (
	"encoding/xml"
	"testing"
)

var contact1 string = `
<?xml version="1.0" encoding="utf-8"?><epp xmlns:_xmlns="xmlns" _xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsi="xsi" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">
        <response>
                <result code="1000">
                        <msg>Command completed successfully</msg>
                </result>
                <resData>
                        <infData xmlns="contact" _xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:contact-1.0 contact-1.0.xsd" _xmlns:contact="urn:ietf:params:xml:ns:contact-1.0">
                                <id xmlns="contact">32-b</id>
                                <roid xmlns="contact">fefefefefefefe-ABCDE</roid>
                                <status xmlns="contact" s="linked"></status>
                                <postalInfo xmlns="contact" type="int">
                                        <name xmlns="contact">John Doe</name>
                                        <addr xmlns="contact">
                                                <street xmlns="contact">1234 Unknown Street</street>
                                                <street xmlns="contact">Apt. 34</street>
                                                <city xmlns="contact">New York</city>
                                                <pc xmlns="contact">99887</pc>
                                                <cc xmlns="contact">US</cc>
                                        </addr>
                                </postalInfo>
                                <email xmlns="contact">testing@example.local</email>
                                <clID xmlns="contact">CLID-1234</clID>
                                <crID xmlns="contact">CLID-2345</crID>
                                <crDate xmlns="contact">2013-01-22T21:23:10.0Z</crDate>
                                <upID xmlns="contact">Registrar119</upID>
                                <upDate xmlns="contact">2020-01-23T23:30:31.0Z</upDate>
                        </infData>
                </resData>
                <trID>
                        <svTRID>fefefefe45964fefefefea1debadf00d</svTRID>
                </trID>
        </response>
</epp>`

func TestContactInfo(t *testing.T) {
	info := ContactInfoResponse{}
	if err := xml.Unmarshal([]byte(contact1), &info); err != nil {
		t.Fatal("Failed to unmarshal info response:", err)
	}

	if info.Result.Code != 1000 {
		t.Error("Expected result code 1000, got:", info.Result.Code)
	}

	if info.PostalInfo.Name != "John Doe" {
		t.Error("Expected John Doe, got:", info.PostalInfo.Name)
	}

	if info.Email != "testing@example.local" {
		t.Error("Expected email 'testing@example.local', got:", info.Email)
	}

	if len(info.PostalInfo.Addr.Street) != 2 {
		t.Error("Expected 2 street elements, got:", len(info.PostalInfo.Addr.Street))
	}

	if info.PostalInfo.Addr.CC != "US" {
		t.Error("Expected country code 'US', got:", info.PostalInfo.Addr.CC)
	}
}
