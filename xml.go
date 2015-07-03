package epp

import (
	"bytes"
	"encoding/xml"
)

const (
	startEPP = `<epp xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="urn:ietf:params:xml:ns:epp-1.0 epp-1.0.xsd" xmlns="urn:ietf:params:xml:ns:epp-1.0">`
	endEPP   = `</epp>`
)

var (
	xmlHeader        = []byte(xml.Header)
	xmlPrefix        = []byte(xml.Header + startEPP)
	xmlSuffix        = []byte(endEPP)
	xmlCommandPrefix = []byte(xml.Header + startEPP + `<command>`)
	xmlCommandSuffix = []byte(`</command>` + endEPP)
)

func deleteRange(s, pfx, sfx []byte) []byte {
	start := bytes.Index(s, pfx)
	if start < 0 {
		return s
	}
	end := bytes.Index(s[start+len(pfx):], sfx)
	if end < 0 {
		return s
	}
	end += start + len(pfx) + len(sfx)
	size := len(s) - (end - start)
	copy(s[start:size], s[end:])
	return s[:size]
}

func deleteBufferRange(buf *bytes.Buffer, pfx, sfx []byte) {
	v := deleteRange(buf.Bytes(), pfx, sfx)
	buf.Truncate(len(v))
}
