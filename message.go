package epp

import (
	"errors"
	"fmt"
)

func (msg *message) error() error {
	if msg.Response == nil || len(msg.Response.Results) == 0 {
		return nil
	}
	r := msg.Response.Results[0]
	if r.IsError() {
		return r
	}
	return nil
}

// message represents a single EPP message.
type message struct {
	XMLName struct{} `xml:"urn:ietf:params:xml:ns:epp-1.0 epp"`

	// Message types. Set to nil if not present in message.
	Greeting *Greeting `xml:"greeting,omitempty"`
	Command  *command  `xml:"command,omitempty"`
	Response *response `xml:"response,omitempty"`
}

// EPP requests

// command represents an EPP command wrapper.
type command struct {
	// Command types. Set to nil if not present in message.
	Login *login `xml:"login,omitempty"`
	Check *check `xml:"check"`

	// TxnID represents a unique client ID for this transaction.
	TxnID string `xml:"clTRID"`
}

// loginCommand authenticates and authorizes an EPP session.
// Supply a non-empty value in NewPassword to change the password for subsequent sessions.
type login struct {
	User        string   `xml:"clID"`
	Password    string   `xml:"pw"`
	NewPassword string   `xml:"newPW,omitempty"`
	Version     string   `xml:"options>version"`
	Language    string   `xml:"options>lang"`
	Objects     []string `xml:"svcs>objURI"`
	Extensions  []string `xml:"svcs>svcExtension>extURI,omitempty"`
}

type check struct {
	DomainCheck *domainCheck `xml:"urn:ietf:params:xml:ns:domain-1.0 check,omitempty"`
}

type domainCheck struct {
	// DomainNS domainNS `xml:"xmlns:domain,attr"`
	Domains []string `xml:"name"`
}

// EPP responses

// response represents an EPP response message.
type response struct {
	Results      []Result `xml:"result"`
	Queue        *queue   `xml:"msgQ,omitempty"`
	TxnID        string   `xml:"trID>clTRID"`
	ServerTxnID  string   `xml:"trID>svTRID"`
	ResponseData struct {
		DomainCheckData *domainCheckData `xml:"urn:ietf:params:xml:ns:domain-1.0 chkData,omitempty"`
		ChargeCheckData *chargeCheckData `xml:"http://www.unitedtld.com/epp/charge-1.0 chkData,omitempty"`
	} `xml:"resData"`
}

// Result represents an EPP server <result> element.
type Result struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"msg"`
}

// IsError determines whether an EPP status code is an error.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsError() bool {
	return r.Code >= 2000
}

// IsFatal determines whether an EPP status code is a fatal response, and the connection should be closed.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsFatal() bool {
	return r.Code >= 2500
}

// Error implements the error interface.
func (r Result) Error() string {
	return fmt.Sprintf("EPP result code %d: %s", r.Code, r.Message)
}

// queue represents an EPP command queue.
type queue struct {
	ID    int  `xml:"id,attr"`
	Count int  `xml:"count,attr"`
	Time  Time `xml:"qDate"`
}

// domainCheckData represents an EPP <domain:chkData> element.
type domainCheckData struct {
	Results []struct {
		Domain struct {
			Domain      string `xml:",chardata"`
			IsAvailable bool   `xml:"avail,attr"`
		} `xml:"name"`
		Reason string `xml:"reason"`
	} `xml:"cd"`
}

// chargeCheckData represents an EPP <charge:chkData> element.
// http://www.unitedtld.com/epp/charge-1.0
type chargeCheckData struct {
	Results []struct {
		Domain   string `xml:"name"`
		Category struct {
			Value string `xml:",chardata"`
			Name  bool   `xml:"name,attr"`
		} `xml:"set>category"`
		Type string `xml:"set>type"`
	} `xml:"cd"`
}

var (
	// ErrResponseNotFound is returned when the EPP XML does not contain a <response> tag.
	ErrResponseNotFound = errors.New("<response> element not found")

	// ErrResponseMalformed is returned when the EPP XML contains a malformed <response> tag.
	ErrResponseMalformed = errors.New("malformed <response> element")

	// ErrResultNotFound is returned when the EPP XML does not contain any <result> tags.
	ErrResultNotFound = errors.New("<result> element not found")
)
