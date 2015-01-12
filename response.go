package epp

import "errors"

// Response represents an EPP <response> element.
type Response struct {
	Results []Result `xml:"result"`
	Queue   *struct {
		ID    int  `xml:"id,attr"`
		Count int  `xml:"count,attr"`
		Time  Time `xml:"qDate"`
	} `xml:"msgQ"`
	TxnID       string `xml:"trID>clTRID"`
	ServerTxnID string `xml:"trID>svTRID"`

	// Individual response types. Set to nil if not present in response message.
	DomainCheck *DomainCheckResponse `xml:"resData>chkData"`
}

// GetResponse returns a response from an EPP msg or an error if
// the message doesn’t contain a valid response.
func (msg *Msg) GetResponse() (*Response, error) {
	switch {
	case msg.Response == nil:
		return nil, ErrMissingResponse
	case len(msg.Response.Results) == 0:
		return nil, ErrMissingResult
	}
	r := msg.Response.Results[0]
	if r.IsError() {
		return nil, r
	}
	return msg.Response, nil
}

var (
	ErrMissingResponse   = errors.New("EPP message did not contain a valid <response> element")
	ErrMalformedResponse = errors.New("EPP message contained a malformed <response> element")
	ErrMissingResult     = errors.New("EPP response did not contain any valid <result> elements")
)
