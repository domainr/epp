package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/std"
	"github.com/domainr/epp/internal/schema/test"
	"github.com/domainr/epp/ns"
)

func TestGreetingRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`empty <greeting>`,
			&epp.EPP{Greeting: &epp.Greeting{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting></greeting></epp>`,
			false,
		},
		{
			`simple <greeting>`,
			&epp.EPP{
				Greeting: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate></greeting></epp>`,
			false,
		},
		{
			`complex <greeting>`,
			&epp.EPP{
				Greeting: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{ns.Contact, ns.Domain, ns.Host},
					},
					DCP: &epp.DCP{},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI></svcMenu><dcp><access></access></dcp></greeting></epp>`,
			false,
		},
		{
			`complex <greeting> with complex <dcp>`,
			&epp.EPP{
				Greeting: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{ns.Contact, ns.Domain, ns.Host},
					},
					DCP: &epp.DCP{
						Access: epp.AccessPersonalAndOther,
						Statements: []epp.Statement{
							{
								Purpose:   epp.PurposeAdmin,
								Recipient: epp.Recipient{Ours: &epp.Ours{Recipient: "Domainr"}, Public: std.True},
							},
							{
								Purpose:   epp.Purpose{Contact: std.True, Other: std.True},
								Recipient: epp.Recipient{Other: std.True, Ours: &epp.Ours{}, Public: std.True},
							},
						},
						Expiry: &epp.Expiry{
							Relative: std.ParseDuration("P1Y").Pointer(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI></svcMenu><dcp><access><personalAndOther/></access><statement><purpose><admin/></purpose><recipient><ours><recDesc>Domainr</recDesc></ours><public/></recipient></statement><statement><purpose><contact/><other/></purpose><recipient><other/><ours/><public/></recipient></statement><expiry><relative>P365DT5H49M12S</relative></expiry></dcp></greeting></epp>`,
			false,
		},
		{
			`<greeting> with <dcp> with absolute expiry`,
			&epp.EPP{
				Greeting: &epp.Greeting{
					DCP: &epp.DCP{
						Expiry: &epp.Expiry{
							Absolute: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><dcp><access></access><expiry><absolute>2000-01-01T00:00:00Z</absolute></expiry></dcp></greeting></epp>`,
			false,
		},
		{
			`complex <greeting> with extensions`,
			&epp.EPP{
				Greeting: &epp.Greeting{
					ServerName: "Test EPP Server",
					ServerDate: std.ParseTime("2000-01-01T00:00:00Z").Pointer(),
					ServiceMenu: &epp.ServiceMenu{
						Versions:  []string{"1.0"},
						Languages: []string{"en", "fr"},
						Objects:   []string{ns.Contact, ns.Domain, ns.Host},
						ServiceExtension: &epp.ServiceExtension{
							Extensions: []string{ns.Fee08, ns.Fee10},
						},
					},
					DCP: &epp.DCP{
						Access: epp.AccessNull,
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting><svID>Test EPP Server</svID><svDate>2000-01-01T00:00:00Z</svDate><svcMenu><version>1.0</version><lang>en</lang><lang>fr</lang><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI><svcExtension><extURI>urn:ietf:params:xml:ns:fee-0.8</extURI><extURI>urn:ietf:params:xml:ns:epp:fee-1.0</extURI></svcExtension></svcMenu><dcp><access><null/></access></dcp></greeting></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
