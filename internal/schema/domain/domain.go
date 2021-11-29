package domain

type Check struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:check"`
	Names   []string `xml:"urn:ietf:params:xml:ns:domain-1.0 domain:name,omitempty"`
}
