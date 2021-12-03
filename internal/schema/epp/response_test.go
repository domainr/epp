package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/test"
)

func TestResponseRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`empty <response>`,
			&epp.EPP{Response: &epp.Response{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`simple code 1000`,
			&epp.EPP{
				Response: &epp.Response{
					Results: []epp.Result{
						{
							Code:    1000,
							Message: epp.Message{Lang: "en", Value: "Command completed successfully"},
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result><code>1000</code><message lang="en">Command completed successfully</message></result><trID><clTRID></clTRID><svTRID></svTRID></trID></response></epp>`,
			false,
		},
		{
			`with transaction IDs`,
			&epp.EPP{
				Response: &epp.Response{
					Results: []epp.Result{
						{
							Code:    1000,
							Message: epp.Message{Lang: "en", Value: "Command completed successfully"},
						},
					},
					TransactionID: epp.TransactionID{
						ClientTransactionID: "12345",
						ServerTransactionID: "abcde",
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><response><result><code>1000</code><message lang="en">Command completed successfully</message></result><trID><clTRID>12345</clTRID><svTRID>abcde</svTRID></trID></response></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
