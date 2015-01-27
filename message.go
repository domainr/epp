package epp

type MessageNamespace struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
}

func (m *MessageNamespace) IsMessage() {}

type Message interface {
	IsMessage()
}
