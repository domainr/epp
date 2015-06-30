package epp

import (
	"encoding/xml"
	"testing"

	"github.com/nbio/st"
)

func TestMarshalEmbeddedStructs(t *testing.T) {
	type EPP struct {
		XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	}

	type Hello struct {
		EPP
		Hello struct{} `xml:"hello"`
	}

	x, _ := xml.Marshal(&Hello{})
	st.Expect(t, string(x), `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello></hello></epp>`)

	type Login struct {
		EPP
		User        string   `xml:"command>login>clID"`
		Password    string   `xml:"command>login>pw"`
		NewPassword string   `xml:"command>login>newPW,omitempty"`
		Version     string   `xml:"command>login>options>version"`
		Language    string   `xml:"command>login>options>lang"`
		Objects     []string `xml:"command>login>svcs>objURI"`
		Extensions  []string `xml:"command>login>svcs>svcExtension>extURI,omitempty"`
	}

	x, _ = xml.Marshal(&Login{})
	st.Expect(t, string(x), `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID></clID><pw></pw><options><version></version><lang></lang></options><svcs><svcExtension></svcExtension></svcs></login></command></epp>`)

	type DomainCheckMsg struct {
		EPP
		Command struct {
			Check struct {
				Domains []string `xml:"name"`
			} `xml:"urn:ietf:params:xml:ns:domain-1.0 check,omitempty"`
		} `xml:"command"`
	}

	x, _ = xml.Marshal(&DomainCheckMsg{})
	st.Expect(t, string(x), `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check xmlns="urn:ietf:params:xml:ns:domain-1.0"></check></command></epp>`)
}
