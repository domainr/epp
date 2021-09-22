package protocol

type EPP struct {
	XMLName      struct{}     `xml:"epp"`
	XMLNamespace eppNamespace `xml:"xmlns,attr"`
	Command      *Command     `xml:",omitempty"`
}

type Command struct {
	XMLName struct{} `xml:"command"`
	Check   *Check   `xml:",omitempty"`
}

type Check struct {
	XMLName     struct{}     `xml:"check"`
	DomainCheck *DomainCheck `xml:",omitempty"`
}

type DomainCheck struct {
	XMLName      struct{}        `xml:"domain:check"`
	XMLNamespace domainNamespace `xml:"xmlns:domain,attr"`
}
