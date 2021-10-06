package protocol

type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Command *Command `xml:"command,omitempty"`
}

type Command struct {
	Check *Check `xml:"urn:ietf:params:xml:ns:epp-1.0 check,omitempty"`
}

type Check struct {
	DomainCheck *DomainCheck `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:check,omitempty"`
}

type DomainCheck struct {
	DomainNames []string `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:name,omitempty"`
}
