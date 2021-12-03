package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/domain"
	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/test"
)

func TestCommandRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`empty <command>`,
			&epp.EPP{Command: &epp.Command{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
		{
			`empty <domain:check> command`,
			&epp.EPP{
				Command: &epp.Command{
					Check: &epp.Check{
						DomainCheck: &domain.Check{},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"></domain:check></check></command></epp>`,
			false,
		},
		{
			`single <domain:check> command`,
			&epp.EPP{
				Command: &epp.Command{
					Check: &epp.Check{
						DomainCheck: &domain.Check{
							Names: []string{"example.com"},
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:check></check></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
