package epp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"time"
)

// Response represents an EPP response message.
type Response struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`
	Results []Result `xml:"response>result"`
	Queue   *struct {
		ID    int  `xml:"id,attr"`
		Count int  `xml:"count,attr"`
		Time  Time `xml:"qDate"`
	} `xml:"response>msgQ"`
	TxnID       string `xml:"response>trID>clTRID"`
	ServerTxnID string `xml:"response>trID>svTRID"`

	// Individual response types. Set to nil if not present in response message.
	Greeting    *Greeting    `xml:"greeting"`
	DomainCheck *DomainCheck `xml:"response>resData>chkData"`
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
