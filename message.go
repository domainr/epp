package epp

import (
	"encoding/xml"
	"errors"
	"fmt"
)

// Marshal encodes an EPP request or message into XML,
// returning any errors that occur.
func Marshal(msg interface{}) ([]byte, error) {
	return xml.Marshal(msg)
}

// Unmarshal decodes an EPP XML response into res,
// returning any errors, including any EPP errors
// received in the response message.
func Unmarshal(data []byte, msg *Message) error {
	err := xml.Unmarshal(data, msg)
	if err != nil {
		return err
	}
	// color.Fprintf(os.Stderr, "@{y}%s\n", spew.Sprintf("%+v", req))
	res := msg.Response
	if res != nil && len(res.Results) > 0 {
		r := res.Results[0]
		if r.IsError() {
			return r
		}
	}
	return nil
}

// Message represents a single EPP message.
type Message struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Request types. Set to nil if not present in message.
	LoginCommand *loginCommand `xml:"command"`

	// Responses types. Set to nil if not present in message.
	Response *Response `xml:"response,omitempty"`
	Greeting *Greeting `xml:"greeting,omitempty"`
}

// Response represents an EPP response message.
type Response struct {
	XMLName struct{} `xml:"response"`
	Results []Result `xml:"result"`
	Queue   *struct {
		ID    int  `xml:"id,attr"`
		Count int  `xml:"count,attr"`
		Time  Time `xml:"qDate"`
	} `xml:"msgQ,omitempty"`
	TxnID       string       `xml:"trID>clTRID"`
	ServerTxnID string       `xml:"trID>svTRID"`
	DomainCheck *DomainCheck `xml:"resData>chkData,omitempty"`
}

var (
	ErrResponseNotFound  = errors.New("<response> element not found")
	ErrResponseMalformed = errors.New("malformed <response> element")
	ErrResultNotFound    = errors.New("<result> element not found")
)

// Result represents an EPP server <result> element.
type Result struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"msg"`
}

// IsError() determines whether an EPP status code is an error.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsError() bool {
	return r.Code >= 2000
}

// IsFatal() determines whether an EPP status code is a fatal response, and the connection should be closed.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsFatal() bool {
	return r.Code >= 2500
}

// Error() implements the error interface.
func (r Result) Error() string {
	return fmt.Sprintf("EPP result code %d: %s", r.Code, r.Message)
}
