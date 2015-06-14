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

	// Messages types. Set to nil if not present in message.
	Command  *command
	Response *Response
	Greeting *Greeting `xml:"greeting,omitempty"`
}

// command represents an EPP command wrapper.
type command struct {
	XMLName struct{} `xml:"command,omitempty"`

	// Command types. Set to nil if not present in message.
	Login *login

	// TxnID represents a unique client ID for this transaction.
	TxnID string `xml:"clTRID"`
}

// loginCommand authenticates and authorizes an EPP session.
// Supply a non-empty value in NewPassword to change the password for subsequent sessions.
type login struct {
	XMLName     struct{} `xml:"login,omitempty"`
	User        string   `xml:"clID"`
	Password    string   `xml:"pw"`
	NewPassword string   `xml:"newPW,omitempty"`
	Version     string   `xml:"options>version"`
	Language    string   `xml:"options>lang"`
	Objects     []string `xml:"svcs>objURI"`
	Extensions  []string `xml:"svcs>svcExtension>extURI,omitempty"`
}

// Response represents an EPP response message.
type Response struct {
	XMLName      struct{}      `xml:"response,omitempty"`
	Results      []Result      `xml:"result"`
	Queue        *Queue        `xml:"msgQ,omitempty"`
	TxnID        string        `xml:"trID>clTRID"`
	ServerTxnID  string        `xml:"trID>svTRID"`
	ResponseData *responseData `xml:"resData,omitempty"`
	// DomainCheck  *DomainCheck  `xml:"resData>chkData,omitempty"`
}

// Result represents an EPP server <result> element.
type Result struct {
	XMLName struct{} `xml:"result"`
	Code    int      `xml:"code,attr"`
	Message string   `xml:"msg"`
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

// Queue represents an EPP command queue.
type Queue struct {
	XMLName struct{} `xml:"msgQ,omitempty"`
	ID      int      `xml:"id,attr"`
	Count   int      `xml:"count,attr"`
	Time    Time     `xml:"qDate"`
}

// responseData represents an EPP <resData> element.
type responseData struct {
	XMLName         struct{}         `xml:"resData"`
	DomainCheckData *DomainCheckData `xml:"urn:ietf:params:xml:ns:domain-1.0 chkData,omitempty"`
}

type DomainCheckData struct {
	Results []struct {
		Domain struct {
			Domain      string `xml:",chardata"`
			IsAvailable bool   `xml:"avail,attr"`
		} `xml:"name"`
		Reason string `xml:"reason"`
	} `xml:"cd"`
}

// domainNS type exists solely to emit an xmlns:domain attribute.
type domainNS struct{}

// MarshalText returns a byte slice for the xmlns:domain attribute.
func (n domainNS) MarshalText() (text []byte, err error) {
	return nsDomain, nil
}

var nsDomain = []byte("urn:ietf:params:xml:ns:domain-1.0")

var (
	ErrResponseNotFound  = errors.New("<response> element not found")
	ErrResponseMalformed = errors.New("malformed <response> element")
	ErrResultNotFound    = errors.New("<result> element not found")
)
