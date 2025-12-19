package epp

import (
	"encoding/xml"
	"strconv"
)

func xmlInt(i int) string {
	return strconv.Itoa(i)
}

const (
	// EPP defines the IETF URN for the EPP namespace.
	// https://www.iana.org/assignments/xml-registry/ns/epp-1.0.txt
	EPP = `urn:ietf:params:xml:ns:epp-1.0`

	// EPPCommon defines the IETF URN for the EPP Common namespace.
	// https://www.iana.org/assignments/xml-registry/ns/eppcom-1.0.txt
	EPPCommon = `urn:ietf:params:xml:ns:eppcom-1.0`

	startEPP         = `<epp xmlns="` + EPP + `">`
	endEPP           = `</epp>`
	xmlCommandPrefix = xml.Header + startEPP + `<command>`
	xmlCommandSuffix = `</command>` + endEPP
)
