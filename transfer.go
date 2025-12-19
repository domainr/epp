package epp

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"time"

	"github.com/nbio/xx"
)

// TransferOp specifies an EPP transfer operation.
// https://www.rfc-editor.org/rfc/rfc5730#section-2.9.2.1
type TransferOp string

const (
	TransferRequest TransferOp = "request"
	TransferQuery   TransferOp = "query"
	TransferApprove TransferOp = "approve"
	TransferReject  TransferOp = "reject"
	TransferCancel  TransferOp = "cancel"
)

// Period represents an EPP <period> value with a unit (typically "y" for years).
type Period struct {
	Value int
	Unit  string
}

// TransferDomain sends a domain transfer command and returns the resulting transfer data, if present.
// https://www.rfc-editor.org/rfc/rfc5731#section-3.2.3
//
// extData is currently unused (reserved for future extensions like fee/price).
func (c *Conn) TransferDomain(op TransferOp, domain string, authInfo string, period *Period, extData map[string]string) (*DomainTransferResponse, error) {
	x, err := encodeDomainTransfer(&c.Greeting, op, domain, authInfo, period)
	if err != nil {
		return nil, err
	}
	if err := c.writeRequest(x); err != nil {
		return nil, err
	}
	res, err := c.readResponse()
	if err != nil {
		return nil, err
	}
	return &res.DomainTransferResponse, nil
}

func encodeDomainTransfer(greeting *Greeting, op TransferOp, domain string, authInfo string, period *Period) ([]byte, error) {
	switch op {
	case TransferRequest, TransferQuery, TransferApprove, TransferReject, TransferCancel:
	default:
		return nil, fmt.Errorf("invalid transfer op %q", string(op))
	}

	buf := bytes.NewBufferString(xmlCommandPrefix)
	buf.WriteString(`<transfer op="`)
	xml.EscapeText(buf, []byte(op))
	buf.WriteString(`">`)
	buf.WriteString(`<domain:transfer xmlns:domain="`)
	buf.WriteString(ObjDomain)
	buf.WriteString(`">`)
	buf.WriteString(`<domain:name>`)
	xml.EscapeText(buf, []byte(domain))
	buf.WriteString(`</domain:name>`)

	// Period is optional; typically used with request/approve.
	if period != nil && period.Value > 0 {
		unit := period.Unit
		if unit == "" {
			unit = "y"
		}
		buf.WriteString(`<domain:period unit="`)
		xml.EscapeText(buf, []byte(unit))
		buf.WriteString(`">`)
		buf.WriteString(xmlInt(period.Value))
		buf.WriteString(`</domain:period>`)
	}

	// authInfo is often required for request and may be required for other ops.
	if authInfo != "" {
		buf.WriteString(`<domain:authInfo><domain:pw>`)
		xml.EscapeText(buf, []byte(authInfo))
		buf.WriteString(`</domain:pw></domain:authInfo>`)
	}

	buf.WriteString(`</domain:transfer>`)
	buf.WriteString(`</transfer>`)
	buf.WriteString(xmlCommandSuffix)
	return buf.Bytes(), nil
}

// DomainTransferResponse represents an EPP response transfer data block for a domain transfer command.
// https://www.rfc-editor.org/rfc/rfc5731#section-3.2.3
type DomainTransferResponse struct {
	Name     string
	TrStatus string
	ReID     string
	ReDate   time.Time
	AcID     string
	AcDate   time.Time
	ExDate   *time.Time
}

func init() {
	path := "epp > response > resData > " + ObjDomain + " trnData"
	scanResponse.MustHandleCharData(path+">name", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.Name = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">trStatus", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.TrStatus = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">reID", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.ReID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">reDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		var err error
		dtr.ReDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">acID", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		dtr.AcID = string(c.CharData)
		return nil
	})
	scanResponse.MustHandleCharData(path+">acDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		var err error
		dtr.AcDate, err = time.Parse(time.RFC3339, string(c.CharData))
		return err
	})
	scanResponse.MustHandleCharData(path+">exDate", func(c *xx.Context) error {
		dtr := &c.Value.(*Response).DomainTransferResponse
		t, err := time.Parse(time.RFC3339, string(c.CharData))
		if err != nil {
			return err
		}
		dtr.ExDate = &t
		return nil
	})
}



