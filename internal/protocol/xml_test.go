package protocol

import (
	"encoding/xml"
	"testing"
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
			`empty <epp> tag`,
			&EPP{},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"></epp>`,
			false,
		},
		{
			`empty <command> tag`,
			&EPP{Command: &Command{}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command></command></epp>`,
			false,
		},
		{
			`empty <domain:check> command`,
			&EPP{Command: &Command{Check: &Check{DomainCheck: &DomainCheck{}}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><check><domain:check xmlns:domain="urn:ietf:params:xml:ns:domain-1.0"></domain:check></check></command></epp>`,
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
			got := string(x)
			if got != tt.want {
				t.Errorf("xml.Marshal() = %v, want %v", got, tt.want)
			}
		})
	}
}
