package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/test"
)

func TestLoginRoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		v       interface{}
		want    string
		wantErr bool
	}{
		{
			`empty <login>`,
			&epp.EPP{Command: &epp.Command{Login: &epp.Login{}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID></clID><pw></pw><options><version></version></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`simple <login>`,
			&epp.EPP{
				Command: &epp.Command{
					Login: &epp.Login{
						ClientID: "user",
						Password: "password",
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>user</clID><pw>password</pw><options><version></version></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`specify version 1.0`,
			&epp.EPP{
				Command: &epp.Command{
					Login: &epp.Login{
						ClientID: "user",
						Password: "password",
						Options: epp.Options{
							Version: epp.Version,
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>user</clID><pw>password</pw><options><version>1.0</version></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`specify lang=en`,
			&epp.EPP{
				Command: &epp.Command{
					Login: &epp.Login{
						ClientID: "user",
						Password: "password",
						Options: epp.Options{
							Version: epp.Version,
							Lang:    "en",
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>user</clID><pw>password</pw><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`change password`,
			&epp.EPP{
				Command: &epp.Command{
					Login: &epp.Login{
						ClientID:    "user",
						Password:    "password",
						NewPassword: "newpassword",
						Options: epp.Options{
							Version: epp.Version,
							Lang:    "en",
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>user</clID><pw>password</pw><newPW>newpassword</newPW><options><version>1.0</version><lang>en</lang></options><svcs></svcs></login></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
