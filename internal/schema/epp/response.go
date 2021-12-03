package epp

// Response represents an EPP server <response> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.6.
type Response struct {
	Results []Result `xml:"result,omitempty"`

	// The OPTIONAL <msgQ> element describes messages queued for client
	// retrieval.
	MessageQueue *MessageQueue `xml:"msgQ"`

	// The <trID> (transaction identifier) element contains the
	// transaction identifier assigned by the server to the command for
	// which the response is being returned.
	TransactionID TransactionID `xml:"trID"`
}

// Result represents an EPP server <result> as defined in RFC 5730.
type Result struct {
	Code    ResultCode `xml:"code"`
	Message Message    `xml:"message"`
	// TODO: Values
	ExtensionValues []ExtensionValue `xml:"extValue,omitempty"`
}

// ExtensionValue represents an extension to an EPP command result.
type ExtensionValue struct {
	// TODO: value
	Reason Message `xml:"reason"`
}

// TransactionID represents an EPP server <trID> as defined in RFC 5730.
type TransactionID struct {
	ClientTransactionID string `xml:"clTRID"`
	ServerTransactionID string `xml:"svTRID"`
}
