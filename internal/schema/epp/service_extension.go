package epp

// ServiceExtension represents an EPP <svcExtension> element as defined in RFC 5730.
// Used in EPP <greeting> and <login> messages.
type ServiceExtension struct {
	Extensions []string `xml:"extURI"`
}
