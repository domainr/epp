package epp

import "encoding/xml"

// Login initializes an authenticated EPP session.
func (c *Conn) Login(user, password, newPassword string) error {
	err := c.encodeLogin(user, password, newPassword)
	if err != nil {
		return err
	}
	err = c.flushDataUnit()
	if err != nil {
		return err
	}
	msg := message{}
	return c.readMessage(&msg)
}

func (c *Conn) encodeLogin(user, password, newPassword string) error {
	c.buf.Reset()
	c.buf.Write(xmlCommandPrefix)
	c.buf.WriteString(`<login><clID>`)
	xml.EscapeText(&c.buf, []byte(user))
	c.buf.WriteString(`</clID><pw>`)
	xml.EscapeText(&c.buf, []byte(password))
	if len(newPassword) > 0 {
		c.buf.WriteString(`</pw><newPW>`)
		xml.EscapeText(&c.buf, []byte(newPassword))
		c.buf.WriteString(`</newPW><options>`)
	} else {
		c.buf.WriteString(`</pw><options>`)
	}
	if len(c.Greeting.Versions) > 0 {
		c.buf.WriteString(`<version>`)
		xml.EscapeText(&c.buf, []byte(c.Greeting.Versions[0]))
		c.buf.WriteString(`</version>`)
	} else {
		c.buf.WriteString(`<version>1.0</version>`)
	}
	if len(c.Greeting.Languages) > 0 {
		c.buf.WriteString(`<lang>`)
		xml.EscapeText(&c.buf, []byte(c.Greeting.Languages[0]))
		c.buf.WriteString(`</lang>`)
	} else {
		c.buf.WriteString(`<lang>en</lang>`)
	}
	c.buf.WriteString(`</options><svcs>`)
	for _, o := range c.Greeting.Objects {
		c.buf.WriteString(`<objURI>`)
		xml.EscapeText(&c.buf, []byte(o))
		c.buf.WriteString(`</objURI>`)
	}
	if len(c.Greeting.Extensions) > 0 {
		c.buf.WriteString(`<svcExtension>`)
		for _, o := range c.Greeting.Extensions {
			c.buf.WriteString(`<extURI>`)
			xml.EscapeText(&c.buf, []byte(o))
			c.buf.WriteString(`</extURI>`)
		}
		c.buf.WriteString(`</svcExtension>`)
	}
	c.buf.WriteString(`</svcs></login>`)
	c.encodeID()
	c.buf.Write(xmlCommandSuffix)
	return nil
}

var xmlLogin = `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0"><command><login><clID>domainr</clID><pw>lgOEWu5rJA</pw><options><version>1.0</version><lang>en</lang></options><svcs><objURI>urn:ietf:params:xml:ns:domain-1.0</objURI><objURI>urn:ietf:params:xml:ns:host-1.0</objURI><objURI>urn:ietf:params:xml:ns:contact-1.0</objURI><objURI>http://www.unitedtld.com/epp/finance-1.0</objURI><svcExtension><extURI>urn:ietf:params:xml:ns:secDNS-1.1</extURI><extURI>urn:ietf:params:xml:ns:rgp-1.0</extURI><extURI>urn:ietf:params:xml:ns:launch-1.0</extURI><extURI>urn:ietf:params:xml:ns:idn-1.0</extURI><extURI>http://www.unitedtld.com/epp/charge-1.0</extURI></svcExtension></svcs></login><clTRID>0000000000000001</clTRID></command></epp>`
