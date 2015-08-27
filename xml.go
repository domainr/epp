package epp

import "encoding/xml"

const (
	startEPP         = `<epp xmlns="urn:ietf:params:xml:ns:epp-1.0">`
	endEPP           = `</epp>`
	xmlCommandPrefix = xml.Header + startEPP + `<command>`
	xmlCommandSuffix = `</command>` + endEPP
)
