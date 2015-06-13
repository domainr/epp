package epp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

// Marshal encodes an EPP request or message into XML,
// returning any errors that occur.
func Marshal(msg interface{}) ([]byte, error) {
	return xml.Marshal(msg)
}

// Unmarshal decodes an EPP XML response into res,
// returning any errors.
func Unmarshal(data []byte, res *Response) error {
	err := xml.Unmarshal(data, res)
	if err != nil {
		return err
	}
	// color.Fprintf(os.Stderr, "@{y}%s\n", spew.Sprintf("%+v", req))
	if len(res.Results) != 0 {
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

	// Subordinate message types. Set to nil if not present in message.
	Response *Response2 `xml:"response,omitempty"`
	Greeting *Greeting  `xml:"greeting,omitempty"`
}

// Response represents an EPP response message.
type Response2 struct {
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

// Response represents an EPP response message.
type Response struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Results []Result `xml:"response>result,omitempty"`
	Queue   *struct {
		ID    int  `xml:"id,attr"`
		Count int  `xml:"count,attr"`
		Time  Time `xml:"qDate"`
	} `xml:"response>msgQ,omitempty"`
	TxnID       string `xml:"response>trID>clTRID,omitempty"`
	ServerTxnID string `xml:"response>trID>svTRID,omitempty"`

	// Individual response types. Set to nil if not present in response message.
	Greeting    *Greeting    `xml:"greeting,omitempty"`
	DomainCheck *DomainCheck `xml:"response>resData>chkData,omitempty"`
}

var (
	ErrMalformedResponse = errors.New("EPP message contained a malformed <response> element")
	ErrMissingResult     = errors.New("EPP response did not contain any valid <result> elements")
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

// Time represents EPP date-time values.
type Time struct {
	time.Time
}

// UnmarshalXML implements a custom XML unmarshaler that ignores time parsing errors.
// http://stackoverflow.com/a/25015260
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	d.DecodeElement(&v, &start)
	if parse, err := time.Parse(time.RFC3339, v); err == nil {
		*t = Time{parse}
	}
	return nil
}
