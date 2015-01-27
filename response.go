package epp

import "errors"

// Response represents an EPP <response> element.
type ResponseMessage struct {
	MessageNamespace
	Results []Result `xml:"response>result"`
	Queue   *struct {
		ID    int  `xml:"id,attr"`
		Count int  `xml:"count,attr"`
		Time  Time `xml:"qDate"`
	} `xml:"response>msgQ"`
	TxnID       string `xml:"response>trID>clTRID"`
	ServerTxnID string `xml:"response>trID>svTRID"`

	// Individual response types. Set to nil if not present in response message.
	DomainCheck *DomainCheckResponse `xml:"response>resData>chkData"`
}

var (
	ErrMissingResponse   = errors.New("EPP message did not contain a valid <response> element")
	ErrMalformedResponse = errors.New("EPP message contained a malformed <response> element")
	ErrMissingResult     = errors.New("EPP response did not contain any valid <result> elements")
)
