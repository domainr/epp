package epp

// Login represents an EPP <login> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.1.1.
type Login struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 login"`
	ClientID    string   `xml:"clID"`
	Password    string   `xml:"pw"`
	NewPassword *string  `xml:"newPW"`
	Options     Options  `xml:"options"`
	Services    Services `xml:"svcs"`
	command
}

func (Login) eppCommand() {}

// Options represent EPP login options as defined in RFC 5730.
type Options struct {
	Version string `xml:"version"`
	Lang    string `xml:"lang,omitempty"`
}

// Services represent EPP login services as defined in RFC 5730.
type Services struct {
	Objects          []string          `xml:"objURI,omitempty"`
	ServiceExtension *ServiceExtension `xml:"svcExtension"`
}

/*
<?xml version="1.0" encoding="UTF-8" standalone="no"?>
<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
	<command>
		<login>
			<clID>ClientX</clID>
			<pw>foo-BAR2</pw>
			<newPW>bar-FOO2</newPW>
			<options>
				<version>1.0</version>
				<lang>en</lang>
			</options>
			<svcs>
				<objURI>urn:ietf:params:xml:ns:obj1</objURI>
				<objURI>urn:ietf:params:xml:ns:obj2</objURI>
				<objURI>urn:ietf:params:xml:ns:obj3</objURI>
				<svcExtension>
					<extURI>http://custom/obj1ext-1.0</extURI>
				</svcExtension>
			</svcs>
		</login>
		<clTRID>ABC-12345</clTRID>
	</command>
</epp>
*/
