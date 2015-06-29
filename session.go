package epp

// Login initializes an authenticated EPP session.
func (c *Conn) Login(user, password, newPassword string) (err error) {
	req := message{
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
	// FIXME: find the highest protocol version?
	// Do any EPP servers send anything other than 1.0?
	if len(c.Greeting.Versions) > 0 {
		req.Command.Login.Version = c.Greeting.Versions[0]
	}
	// FIXME: look for a particular language?
	// Do any EPP servers send anything other than “en”?
	if len(c.Greeting.Languages) > 0 {
		req.Command.Login.Language = c.Greeting.Languages[0]
	}
	// FIXME: we currently just echo back what’s reported by the server.
	// We may or may not use any of these in a given session. Optimization opportunity?
	req.Command.Login.Objects = c.Greeting.Objects
	req.Command.Login.Extensions = c.Greeting.Extensions
	err = c.writeMessage(&req)
	if err != nil {
		return
	}
	msg := message{}
	return c.readMessage(&msg)
}
