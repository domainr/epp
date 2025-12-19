package epp

import (
	"bytes"
	"encoding/xml"
	"testing"
	"time"

	"github.com/nbio/st"
)

func TestEncodeDomainTransferRequest(t *testing.T) {
	x, err := encodeDomainTransfer(nil, TransferRequest, "example.com", "auth123", &Period{Value: 1, Unit: "y"})
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer op="request"><domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name><domain:period unit="y">1</domain:period><domain:authInfo><domain:pw>auth123</domain:pw></domain:authInfo></domain:transfer></transfer></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainTransferQueryMinimal(t *testing.T) {
	x, err := encodeDomainTransfer(nil, TransferQuery, "example.com", "", nil)
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><transfer op="query"><domain:transfer xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:transfer></transfer></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeDomainTransferInvalidOp(t *testing.T) {
	_, err := encodeDomainTransfer(nil, TransferOp("nope"), "example.com", "", nil)
	st.Reject(t, err, nil)
}

func TestScanDomainTransferResponse(t *testing.T) {
	x := `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
  <response>
    <result code="1001"><msg>Command completed successfully; action pending</msg></result>
    <resData>
      <domain:trnData xmlns:domain="urn:ietf:params:xml:ns:domain-1.0">
        <domain:name>example.com</domain:name>
        <domain:trStatus>pending</domain:trStatus>
        <domain:reID>ClientX</domain:reID>
        <domain:reDate>2025-01-01T00:00:00Z</domain:reDate>
        <domain:acID>ClientY</domain:acID>
        <domain:acDate>2025-01-02T00:00:00Z</domain:acDate>
        <domain:exDate>2026-01-01T00:00:00Z</domain:exDate>
      </domain:trnData>
    </resData>
  </response>
</epp>`

	var res Response
	d := xml.NewDecoder(bytes.NewBufferString(x))
	err := IgnoreEOF(scanResponse.Scan(d, &res))
	st.Expect(t, err, nil)

	dtr := res.DomainTransferResponse
	st.Expect(t, dtr.Name, "example.com")
	st.Expect(t, dtr.TrStatus, "pending")
	st.Expect(t, dtr.ReID, "ClientX")
	st.Expect(t, dtr.AcID, "ClientY")
	st.Expect(t, dtr.ReDate, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC))
	st.Expect(t, dtr.AcDate, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC))
	st.Reject(t, dtr.ExDate, (*time.Time)(nil))
	st.Expect(t, *dtr.ExDate, time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
}



