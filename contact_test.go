package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeContactCreate(t *testing.T) {
	pi := PostalInfo{
		Name:   "John Doe",
		Org:    "Acme Corp",
		Street: "123 Main St",
		City:   "Metropolis",
		SP:     "NY",
		PC:     "10001",
		CC:     "US",
	}

	// Test case 1: With voice
	x, err := encodeContactCreate(nil, "contactID", "user@example.com", pi, "+1.5555555555", "auth123", nil)
	st.Expect(t, err, nil)

	expectedWithVoice := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create><contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>contactID</contact:id><contact:postalInfo type="int"><contact:name>John Doe</contact:name><contact:org>Acme Corp</contact:org><contact:addr><contact:street>123 Main St</contact:street><contact:city>Metropolis</contact:city><contact:sp>NY</contact:sp><contact:pc>10001</contact:pc><contact:cc>US</contact:cc></contact:addr></contact:postalInfo><contact:voice>+1.5555555555</contact:voice><contact:email>user@example.com</contact:email><contact:authInfo><contact:pw>auth123</contact:pw></contact:authInfo></contact:create></create></command></epp>`

	st.Expect(t, string(x), expectedWithVoice)

	// Test case 2: Without voice
	x, err = encodeContactCreate(nil, "contactID", "user@example.com", pi, "", "auth123", nil)
	st.Expect(t, err, nil)

	expectedWithoutVoice := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create><contact:create xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>contactID</contact:id><contact:postalInfo type="int"><contact:name>John Doe</contact:name><contact:org>Acme Corp</contact:org><contact:addr><contact:street>123 Main St</contact:street><contact:city>Metropolis</contact:city><contact:sp>NY</contact:sp><contact:pc>10001</contact:pc><contact:cc>US</contact:cc></contact:addr></contact:postalInfo><contact:email>user@example.com</contact:email><contact:authInfo><contact:pw>auth123</contact:pw></contact:authInfo></contact:create></create></command></epp>`

	st.Expect(t, string(x), expectedWithoutVoice)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
