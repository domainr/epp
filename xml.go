package epp

import (
	"github.com/nbio/xml"

	"github.com/domainr/epp/ns"
)

const (
	startEPP         = `<epp xmlns="` + ns.EPP + `">`
	endEPP           = `</epp>`
	xmlCommandPrefix = xml.Header + startEPP + `<command>`
	xmlCommandSuffix = `</command>` + endEPP
)
