package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeDomainTransfer(t *testing.T) {
	// Test case 1: Basic domain transfer query
	x, err := encodeDomainTransfer(nil, "query", "example.com", 0, "y", "", nil)
	st.Expect(t, err, nil)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer op="query"><domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:transfer></transfer></command></epp>`

	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainTransferRequest(t *testing.T) {
	// Test case 2: Domain transfer request with auth and period
	x, err := encodeDomainTransfer(nil, "request", "example.com", 1, "y", "secret", nil)
	st.Expect(t, err, nil)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer op="request"><domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">1</domain:period><domain:authInfo><domain:pw>secret</domain:pw></domain:authInfo></domain:transfer></transfer></command></epp>`

	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainTransferWithFee(t *testing.T) {
	// Test case 3: Domain transfer with fee extension
	extData := map[string]string{
		"fee:fee":      "10.00",
		"fee:currency": "USD",
	}
	x, err := encodeDomainTransfer(nil, "request", "example.com", 1, "y", "secret", extData)
	st.Expect(t, err, nil)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer op="request"><domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">1</domain:period><domain:authInfo><domain:pw>secret</domain:pw></domain:authInfo></domain:transfer></transfer><extension><fee:transfer xmlns:fee="urn:ietf:params:xml:ns:epp:fee-1.0"><fee:currency>USD</fee:currency><fee:fee>10.00</fee:fee></fee:transfer></extension></command></epp>`

	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
