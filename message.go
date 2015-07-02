package epp

import "errors"

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
	ResponseData struct {
		DomainCheckData *domainCheckData `xml:"urn:ietf:params:xml:ns:domain-1.0 chkData,omitempty"`
	} `xml:"resData"`
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

var (
	// ErrResponseNotFound is returned when the EPP XML does not contain a <response> tag.
	ErrResponseNotFound = errors.New("<response> element not found")

	// ErrResponseMalformed is returned when the EPP XML contains a malformed <response> tag.
	ErrResponseMalformed = errors.New("malformed <response> element")

	// ErrResultNotFound is returned when the EPP XML does not contain any <result> tags.
	ErrResultNotFound = errors.New("<result> element not found")
)
