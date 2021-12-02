package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/test"
)

func TestEPPRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`nil`,
			nil,
			``,
			false,
		},
		{
			`empty <epp> element`,
			&epp.EPP{},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"></epp>`,
			false,
		},
		{
			`empty <hello> message`,
			&epp.EPP{Hello: &epp.Hello{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello/></epp>`,
			false,
		},
		{
			`empty <greeting>`,
			&epp.EPP{Greeting: &epp.Greeting{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><greeting></greeting></epp>`,
			false,
		},
		{
			`empty <command>`,
			&epp.EPP{Command: &epp.Command{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
