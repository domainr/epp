package epp

import (
	"bytes"
	"encoding/xml"
	"io"
	"strconv"
)

// Login initializes an authenticated EPP session.
func (c *Conn) Login(user, password, newPassword string) error {
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
	err = c.flushDataUnit()
	if err != nil {
		return err
	}
	msg := message{}
	return c.readMessage(&msg)
}

func encodeLogin(buf *bytes.Buffer, user, password, newPassword, version, language string, objects, extensions []string) error {
	buf.Reset()
	buf.Write(xmlCommandPrefix)
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
	buf.Write(xmlCommandSuffix)
	return nil
}

func decodeLoginResponse(d *Decoder) (Result, error) {
	d.Reset()
	return decodeResult(d)
}

func decodeResult(d *Decoder) (Result, error) {
	var r Result
outer:
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return r, err
		}
		switch node := t.(type) {
		case xml.StartElement:
			if node.Name.Local == "result" {
				for _, a := range node.Attr {
					if a.Name.Local == "code" {
						r.Code, _ = strconv.Atoi(a.Value)
						break
					}
				}
			}

		case xml.EndElement:
			// Escape early (skip remaining XML)
			if node.Name.Local == "result" {
				break outer
			}

		case xml.CharData:
			e := d.Element()
			if e != nil && e.Name.Local == "msg" {
				r.Message = string(node)
			}
		}

		// Escape early (skip remaining XML)
		if r.Code > 0 && r.Message != "" {
			break
		}
	}
	return r, nil
}
