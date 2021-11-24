package protocol

type EPP struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Command *Command
}

type Command struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 command"`
	Check   *Check
}

type Check struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 check"`
	DomainCheck *DomainCheck
}

type DomainCheck struct {
	XMLName     struct{} `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:check"`
	DomainNames []string `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:name,omitempty"`
}
