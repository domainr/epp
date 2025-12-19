package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeDomainRestore(t *testing.T) {
	// Test case 1: Domain restore
	x, err := encodeDomainRestore(nil, "example.com", nil)
	st.Expect(t, err, nil)

	expected := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><update><domain:update xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:update></update><extension><rgp:update xmlns:rgp="urn:ietf:params:xml:ns:rgp-1.0"><rgp:restore op="request"/></rgp:update></extension></command></epp>`

	st.Expect(t, string(x), expected)

	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}
