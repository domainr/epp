package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeDomainInfo(t *testing.T) {
	// Test case 1: Basic domain info
	x, err := encodeDomainInfo(nil, "example.com", nil)
	st.Expect(t, err, nil)

	expectedDomain := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><info><domain:info xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name hosts="none">example.com</domain:name></domain:info></info></command></epp>`

	st.Expect(t, string(x), expectedDomain)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeContactInfo(t *testing.T) {
	// Test case 1: Contact info without auth
	x, err := encodeContactInfo(nil, "contactID", "", nil)
	st.Expect(t, err, nil)

	expectedContactNoAuth := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><info><contact:info xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>contactID</contact:id></contact:info></info></command></epp>`

	st.Expect(t, string(x), expectedContactNoAuth)

	// Test case 2: Contact info with auth
	x, err = encodeContactInfo(nil, "contactID", "auth123", nil)
	st.Expect(t, err, nil)

	expectedContactWithAuth := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><info><contact:info xmlns:contact="urn:ietf:params:xml:ns:contact-1.0"><contact:id>contactID</contact:id><contact:authInfo><contact:pw>auth123</contact:pw></contact:authInfo></contact:info></info></command></epp>`

	st.Expect(t, string(x), expectedContactWithAuth)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
