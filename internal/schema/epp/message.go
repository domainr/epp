package epp

// Message represents an human-readable message + optional language identifier.
// Used in epp>response>result>msg and epp>response>result>extValue>reason.
type Message struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}
