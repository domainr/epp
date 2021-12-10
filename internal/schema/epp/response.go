package epp

// Response represents an EPP server <response> as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.6.
type Response struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 response"`

	// Results contain one or more results (success or failure) of an EPP command.
	Results []Result `xml:"result,omitempty"`

	// The OPTIONAL <msgQ> element describes messages queued for client
	// retrieval.
	MessageQueue *MessageQueue `xml:"msgQ"`

	// The <trID> (transaction identifier) element contains a client
	// transaction ID of the command that elicited this response and a
	// server transaction ID that uniquely identifies this response.
	TransactionID TransactionID `xml:"trID"`
}

func (Response) eppBody() {}

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
	Client string `xml:"clTRID"`
	Server string `xml:"svTRID"`
}
