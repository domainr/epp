package epp

// Login authenticates and authorizes an EPP session.
// Supply a non-empty value in NewPassword to change the password for subsequent sessions.
type Login struct {
	XMLName     struct{} `xml:"login"`
	ClientID    string   `xml:"clID"`
	Password    string   `xml:"pw"`
	NewPassword string   `xml:"newPW,omitempty"`
	Version     string   `xml:"options>version"`
	Language    string   `xml:"options>lang"`
	Objects     []string `xml:"svcs>objURI"`
	Extensions  []string `xml:"svcs>svcExtension>extURI,omitempty"`
}

// <epp xmlns="urn:ietf:params:xml:ns:epp-1.0">
//   <command>
//     <login>
//       <clID>ClientX</clID>
//       <pw>foo-BAR2</pw>
//       <newPW>bar-FOO2</newPW>
//       <options>
//         <version>1.0</version>
//         <lang>en</lang>
//       </options>
//       <svcs>
//         <objURI>urn:ietf:params:xml:ns:obj1</objURI>
//         <objURI>urn:ietf:params:xml:ns:obj2</objURI>
//         <objURI>urn:ietf:params:xml:ns:obj3</objURI>
//         <svcExtension>
//           <extURI>http://custom/obj1ext-1.0</extURI>
//         </svcExtension>
//       </svcs>
//     </login>
//     <clTRID>ABC-12345</clTRID>
//   </command>
// </epp>

// Login initializes an authenticated EPP session.
func (c *Conn) Login(clientID, password, newPassword string) (err error) {
	msg := Msg{Command: NewCommand(c.login(clientID, password, newPassword))}
	err = c.WriteMsg(&msg)
	if err != nil {
		return
	}
	_, err = c.ReadResponse()
	return err
}

// login initializes a <login> command.
func (c *Conn) login(clientID, password, newPassword string) *Login {
	cmd := &Login{
		ClientID:    clientID,
		Password:    password,
		NewPassword: newPassword,
		Version:     "1.0",
		Language:    "en",
	}
	if c.Greeting != nil {
		// FIXME: find the highest protocol version?
		// Do any EPP servers send anything other than 1.0?
		if len(c.Greeting.ServiceMenu.Versions) > 0 {
			cmd.Version = c.Greeting.ServiceMenu.Versions[0]
		}
		// FIXME: look for a particular language?
		// Do any EPP servers send anything other than “en”?
		if len(c.Greeting.ServiceMenu.Languages) > 0 {
			cmd.Language = c.Greeting.ServiceMenu.Languages[0]
		}
		// FIXME: we currently just echo back what’s reported by the server.
		// We may or may not use any of these in a given session. Optimization opportunity?
		cmd.Objects = c.Greeting.ServiceMenu.Objects
		cmd.Extensions = c.Greeting.ServiceMenu.Extensions
	}
	return cmd
}
