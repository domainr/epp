package epp

import (
	"bytes"
	"encoding/xml"
)

// Login initializes an authenticated EPP session.
// https://tools.ietf.org/html/rfc5730#section-2.9.1.1
func (c *Conn) Login(user, password, newPassword string) error {
	err := c.writeLogin(user, password, newPassword)
	if err != nil {
		return err
	}
	var res Response
	err = c.readResponse(&res)
	// We always have a .Result in our non-pointer, but it might be meaningless.
	// We might not have read anything.  We think that the worst case is we
	// have the same zero values we'd get without the assignment-even-in-error-case.
	c.LoginResult = res.Result
	return err
}

func (c *Conn) writeLogin(user, password, newPassword string) error {
	ver, lang := "1.0", "en"
	if len(c.Greeting.Versions) > 0 {
		ver = c.Greeting.Versions[0]
	}
	if len(c.Greeting.Languages) > 0 {
		lang = c.Greeting.Languages[0]
	}
	err := encodeLogin(&c.buf, user, password, newPassword, ver, lang, c.Greeting.Objects, c.Greeting.Extensions)
	if err != nil {
		return err
	}
	return c.flushDataUnit()
}

func encodeLogin(buf *bytes.Buffer, user, password, newPassword, version, language string, objects, extensions []string) error {
	buf.Reset()
	buf.WriteString(xmlCommandPrefix)
	buf.WriteString(`<login><clID>`)
	xml.EscapeText(buf, []byte(user))
	buf.WriteString(`</clID><pw>`)
	xml.EscapeText(buf, []byte(password))
	if len(newPassword) > 0 {
		buf.WriteString(`</pw><newPW>`)
		xml.EscapeText(buf, []byte(newPassword))
		buf.WriteString(`</newPW><options><version>`)
	} else {
		buf.WriteString(`</pw><options><version>`)
	}
	xml.EscapeText(buf, []byte(version))
	buf.WriteString(`</version><lang>`)
	xml.EscapeText(buf, []byte(language))
	buf.WriteString(`</lang></options><svcs>`)
	for _, o := range objects {
		buf.WriteString(`<objURI>`)
		xml.EscapeText(buf, []byte(o))
		buf.WriteString(`</objURI>`)
	}
	if len(extensions) > 0 {
		buf.WriteString(`<svcExtension>`)
		for _, o := range extensions {
			buf.WriteString(`<extURI>`)
			xml.EscapeText(buf, []byte(o))
			buf.WriteString(`</extURI>`)
		}
		buf.WriteString(`</svcExtension>`)
	}
	buf.WriteString(`</svcs></login>`)
	buf.WriteString(xmlCommandSuffix)
	return nil
}

// Logout sends a <logout> command to terminate an EPP session.
// https://tools.ietf.org/html/rfc5730#section-2.9.1.2
func (c *Conn) Logout() error {
	err := c.writeDataUnit(xmlLogout)
	if err != nil {
		return err
	}
	var res Response
	return c.readResponse(&res)
}

var xmlLogout = []byte(xmlCommandPrefix + `<logout/>` + xmlCommandSuffix)
