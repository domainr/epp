package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeDomainDelete(t *testing.T) {
	// Test case 1: Basic domain delete
	x, err := encodeDomainDelete(nil, "example.com", nil)
	st.Expect(t, err, nil)

	expectedDomain := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><delete><domain:delete xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:delete></delete></command></epp>`

	st.Expect(t, string(x), expectedDomain)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeContactDelete(t *testing.T) {
	// Test case 1: Contact delete
	x, err := encodeContactDelete(nil, "contactID", nil)
	st.Expect(t, err, nil)

	expectedContact := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><delete><contact:delete xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>contactID</contact:id></contact:delete></delete></command></epp>`

	st.Expect(t, string(x), expectedContact)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
