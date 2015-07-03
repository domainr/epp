package epp

import (
	"encoding/xml"
	"fmt"
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

func init() {
	scanResponse.MustHandleStartElement("epp > response > result", func(c *Context) error {
		res := c.Value.(*response_)
		res.Result = Result{}
		res.Result.Code, _ = strconv.Atoi(c.Attr("", "code"))
		return nil
	})
	scanResponse.MustHandleCharData("epp > response > result > msg", func(c *Context) error {
		c.Value.(*response_).Result.Message = string(c.CharData)
		return nil
	})
}

// decodeResult decodes Result r from Decoder d.
// It does not reset the Decoder.
func (d *Decoder) decodeResult(r *Result) error {
	*r = Result{}
	err := d.DecodeElementWith(d.Element(-1), func(t xml.Token) error {
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

		case xml.CharData:
			if d.AtPath("result", "msg") {
				r.Message = string(node)
			}
		}
		return nil
	})
	if r.IsError() {
		return r
	}
	return err
}
