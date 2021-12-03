package epp

import (
	"strconv"

	"github.com/domainr/epp/internal/schema/std"
	"github.com/nbio/xml"
)

// MessageQueue represents an EPP server <msgQ> as defined in RFC 5730.
type MessageQueue struct {
	// The count attribute describes the number of messages that exist in
	// the queue.
	Count uint64 `xml:"count,attr"`

	// The id attribute is used to uniquely identify the message at the head
	// of the queue.
	ID string `xml:"id,attr"`

	// The <msgQ> element contains the following OPTIONAL child elements
	// that MUST be returned in response to a <poll> request command and
	// MUST NOT be returned in response to any other command, including a
	// <poll> acknowledgement.

	// The <qDate> element that contains the date and time that the message
	// was enqueued.
	Date *std.Time `xml:"qDate"`

	// The <msg> element contains a human-readable message.
	// TODO: This element MAY contain XML content for formatting purposes,
	// but the XML content is not specified by the protocol and will thus
	// not be processed for validity.
	Message *Message `xml:"msg"`
}

// MarshalXML impements the xml.Marshaler interface.
// Writes a single self-closing tag if q.Date and q.Message are not set.
func (q *MessageQueue) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if q.Date == nil && q.Message == nil {
		start.Attr = []xml.Attr{
			{Name: xml.Name{Local: "count"}, Value: strconv.FormatUint(uint64(q.Count), 10)},
			{Name: xml.Name{Local: "id"}, Value: q.ID},
		}
		return e.EncodeToken(xml.SelfClosingElement(start))
	}
	type T MessageQueue
	return e.EncodeElement((*T)(q), start)
}
