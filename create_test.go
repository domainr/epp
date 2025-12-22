package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeDomainCreate(t *testing.T) {
	// Test case 1: Basic domain create
	x, err := encodeDomainCreate(nil, "example.com", 1, "y", "auth123", "", nil, nil, nil)
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create><domain:create xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">1</domain:period><domain:authInfo><domain:pw>auth123</domain:pw></domain:authInfo></domain:create></create></command></epp>`)

	// Test case 2: Full domain create with registrant, contacts, and nameservers
	contacts := map[string]string{
		"admin":   "adminID",
		"tech":    "techID",
		"billing": "billingID",
	}
	ns := []string{"ns1.example.com", "ns2.example.com"}

	x, err = encodeDomainCreate(nil, "example.com", 2, "y", "secret", "regID", contacts, ns, nil)
	st.Expect(t, err, nil)

	// Since map iteration order is random, we can't strict string match easily for contacts.
	// But let's check basic structure or use a fixed order in implementation (we enforced order in implementation: admin, tech, billing).

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create><domain:create xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">2</domain:period><domain:ns><domain:hostObj>ns1.example.com</domain:hostObj><domain:hostObj>ns2.example.com</domain:hostObj></domain:ns><domain:registrant>regID</domain:registrant><domain:contact type="admin">adminID</domain:contact><domain:contact type="tech">techID</domain:contact><domain:contact type="billing">billingID</domain:contact><domain:authInfo><domain:pw>secret</domain:pw></domain:authInfo></domain:create></create></command></epp>`

	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainCreateWithFee(t *testing.T) {
	extData := map[string]string{
		"fee:fee":      "100.00",
		"fee:currency": "USD",
	}
	x, err := encodeDomainCreate(nil, "example.com", 1, "y", "auth123", "", nil, nil, extData)
	st.Expect(t, err, nil)
	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><create><domain:create xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">1</domain:period><domain:authInfo><domain:pw>auth123</domain:pw></domain:authInfo></domain:create></create><extension><fee:create xmlns:fee="urn:ietf:params:xml:ns:epp:fee-1.0"><fee:currency>USD</fee:currency><fee:fee>100.00</fee:fee></fee:create></extension></command></epp>`
	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
