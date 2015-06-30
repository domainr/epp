package epp

import "encoding/xml"

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
