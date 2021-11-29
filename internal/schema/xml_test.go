package epp

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/domainr/epp/internal/schema/domain"
	"github.com/domainr/epp/internal/schema/epp"
	"github.com/nbio/xml"
)

func TestMarshalXML(t *testing.T) {
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
			`empty <epp> message`,
			&epp.EPP{},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"></epp>`,
			false,
		},
		{
			`empty <hello> message`,
			&epp.EPP{Hello: &epp.Hello{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><hello></hello></epp>`,
			false,
		},
		{
			`empty <command> message`,
			&epp.EPP{Command: &epp.Command{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
		{
			`empty <domain:check> command`,
			&epp.EPP{Command: &epp.Command{Check: &epp.Check{DomainCheck: &domain.Check{}}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"></domain:check></check></command></epp>`,
			false,
		},
		{
			`single <domain:check> command`,
			&epp.EPP{Command: &epp.Command{Check: &epp.Check{DomainCheck: &domain.Check{
				Names: []string{"example.com"},
			}}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"><domain:name>example.com</domain:name></domain:check></check></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, err := xml.Marshal(tt.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("xml.Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if string(x) != tt.want {
				t.Errorf("xml.Marshal()\nGot:  %v\nWant: %v", string(x), tt.want)
			}

			if tt.v == nil {
				return
			}

			v := &epp.EPP{}
			err = xml.Unmarshal(x, v)
			if err != nil {
				t.Errorf("xml.Unmarshal() error = %v", err)
				return
			}
			if !reflect.DeepEqual(v, tt.v) {
				// y, _ := xml.Marshal(v)
				t.Errorf("xml.Unmarshal()\nGot:  %v\nWant: %v", asJSON(v), asJSON(tt.v))
			}
		})
	}
}

func asJSON(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
