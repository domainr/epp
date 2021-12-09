package domain

// Check represents an EPP <domain:check> command.
// See https://www.rfc-editor.org/rfc/rfc5730.html.
type Check struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:check"`
	Names   []string `xml:"domain:name,omitempty"`
}

func (Check) EPPCheck() {}
