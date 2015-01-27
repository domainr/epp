package epp

// Msg represents EPP request and response messages as defined by RFC 5730.
// https://tools.ietf.org/html/rfc5730#section-2.2
type Msg struct {
	MessageNamespace

	// Response elements. Will be nil if not present in response message.
	Response *Response `xml:"response"`
}

type MessageNamespace struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
}

func (m *MessageNamespace) IsMessage() {}

type Message interface {
	IsMessage()
}
