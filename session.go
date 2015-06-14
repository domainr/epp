package epp

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
func (c *Conn) Login(user, password, newPassword string) (err error) {
	req := Message{
		Command: &command{
			Login: &login{
				User:        user,
				Password:    password,
				NewPassword: newPassword,
				Version:     "1.0",
				Language:    "en",
			},
			TxnID: c.id(),
		},
	}
	if c.Greeting != nil {
		// FIXME: find the highest protocol version?
		// Do any EPP servers send anything other than 1.0?
		if len(c.Greeting.ServiceVersions) > 0 {
			req.Command.Login.Version = c.Greeting.ServiceVersions[0]
		}
		// FIXME: look for a particular language?
		// Do any EPP servers send anything other than “en”?
		if len(c.Greeting.ServiceLanguages) > 0 {
			req.Command.Login.Language = c.Greeting.ServiceLanguages[0]
		}
		// FIXME: we currently just echo back what’s reported by the server.
		// We may or may not use any of these in a given session. Optimization opportunity?
		req.Command.Login.Objects = c.Greeting.ServiceObjects
		req.Command.Login.Extensions = c.Greeting.ServiceExtensions
	}
	err = c.WriteMessage(&req)
	if err != nil {
		return
	}
	msg := Message{}
	return c.ReadMessage(&msg)
}
