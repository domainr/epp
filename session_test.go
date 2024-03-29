package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestEncodeLogin(t *testing.T) {
	x, err := encodeLogin("jane", "battery", "", "1.0", "en", nil, nil)
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

func TestEncodeLoginChangePassword(t *testing.T) {
	x, err := encodeLogin("jane", "battery", "horse", "1.0", "en", nil, nil)
	st.Expect(t, err, nil)
	st.Expect(t, string(x), `<?xml version="1.0" encoding="UTF-8"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>jane</clID><pw>battery</pw><newPW>horse</newPW><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`)
	var v struct{}
	err = xml.Unmarshal(x, &v)
	st.Expect(t, err, nil)
}

var (
	testObjects = []string{
		ObjContact,
		ObjDomain,
		ObjFinance,
		ObjHost,
	}
	testExtensions = []string{
		ExtCharge,
		ExtFee05,
		ExtFee06,
		ExtIDN,
		ExtLaunch,
		ExtRGP,
		ExtSecDNS,
	}
)

func BenchmarkEncodeLogin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		encodeLogin("jane", "battery", "horse", "1.0", "en", testObjects, testExtensions)
	}
}
