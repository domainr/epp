package epp_test

import (
	"testing"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/domainr/epp/internal/schema/std"
	"github.com/domainr/epp/internal/schema/test"
	"github.com/domainr/epp/ns"
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
			&epp.EPP{Body: &epp.Command{Command: &epp.Login{}}},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID></clID><pw></pw><options><version></version></options><svcs></svcs></login></command></epp>`,
			false,
		},
		{
			`simple <login>`,
			&epp.EPP{
				Body: &epp.Command{
					Command: &epp.Login{
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
				Body: &epp.Command{
					Command: &epp.Login{
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
				Body: &epp.Command{
					Command: &epp.Login{
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
				Body: &epp.Command{
					Command: &epp.Login{
						ClientID:    "user",
						Password:    "password",
						NewPassword: std.StringPointer("newpassword"),
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
		{
			`complex <login>`,
			&epp.EPP{
				Body: &epp.Command{
					Command: &epp.Login{
						ClientID:    "user",
						NewPassword: std.StringPointer("newpassword"),
						Options: epp.Options{
							Version: epp.Version,
							Lang:    "en",
						},
						Services: epp.Services{
							Objects: []string{ns.Domain, ns.Contact, ns.Host},
							ServiceExtension: &epp.ServiceExtension{
								Extensions: []string{
									"urn:ietf:params:xml:ns:epp:fee-0.8",
									"urn:ietf:params:xml:ns:epp:fee-1.0",
									"urn:ietf:params:xml:ns:idn-1.0",
								},
							},
						},
					},
				},
			},
			`<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>user</clID><pw></pw><newPW>newpassword</newPW><options><version>1.0</version><lang>en</lang></options><svcs><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI><svcExtension><extURI>urn:ietf:params:xml:ns:epp:fee-0.8</extURI><extURI>urn:ietf:params:xml:ns:epp:fee-1.0</extURI><extURI>urn:ietf:params:xml:ns:idn-1.0</extURI></svcExtension></svcs></login></command></epp>`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.RoundTrip(t, tt.v, tt.want, tt.wantErr)
		})
	}
}
