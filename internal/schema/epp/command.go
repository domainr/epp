package epp

// Command represents an EPP client <command> as defined in RFC 5730.
type Command struct {
	Login               *Login `xml:"login"`
	Check               *Check `xml:"check"`
	ClientTransactionID string `xml:"clTRID,omitempty"`
}
