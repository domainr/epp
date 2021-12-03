package epp

// Response represents an EPP server <response> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.6.
type Response struct {
	Results      []Result      `xml:"result,omitempty"`
	MessageQueue *MessageQueue `xml:"msgQ"`
}

// Result represents an EPP server <result> as defined in RFC 5730.
type Result struct {
	Code    ResultCode `xml:"code"`
	Message Message    `xml:"message"`
	// TODO: Values
	ExtensionValues []ExtensionValue `xml:"extValue,omitempty"`

	// The OPTIONAL <msgQ> element describes messages queued for client
	// retrieval.
	MessageQueue *MessageQueue `xml:"msgQ"`
}

// ExtensionValue represents an extension to an EPP command result.
type ExtensionValue struct {
	// TODO: value
	Reason Message `xml:"reason"`
}
