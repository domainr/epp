package epp

import (
	"encoding/xml"
	"fmt"
	"io"
	"strconv"
)

// Result represents an EPP <result> element.
type Result struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"msg"`
}

// IsError determines whether an EPP status code is an error.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsError() bool {
	return r.Code >= 2000
}

// IsFatal determines whether an EPP status code is a fatal response,
// and the connection should be closed.
// https://tools.ietf.org/html/rfc5730#section-3
func (r Result) IsFatal() bool {
	return r.Code >= 2500
}

// Error implements the error interface.
func (r Result) Error() string {
	return fmt.Sprintf("EPP result code %d: %s", r.Code, r.Message)
}

// decodeResult decodes a Result from a Decoder.
// It optimistically stops decoding, and may leave
// the Decoder in a half-finished state.
// It does not reset the Decoder.
func decodeResult(d *Decoder, r *Result) error {
outer:
	for {
		t, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		switch node := t.(type) {
		case xml.StartElement:
			if node.Name.Local == "result" {
				for _, a := range node.Attr {
					if a.Name.Local == "code" {
						r.Code, _ = strconv.Atoi(a.Value)
						break
					}
				}
			}

		case xml.EndElement:
			// Escape early (skip remaining XML)
			if node.Name.Local == "result" {
				break outer
			}

		case xml.CharData:
			if d.AtPath("result", "msg") {
				r.Message = string(node)
			}
		}

		// Escape early (skip remaining XML)
		if r.Code > 0 && r.Message != "" {
			break
		}
	}
	if r.IsError() {
		return r
	}
	return nil
}
