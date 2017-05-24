package epp

import (
	"fmt"

	"github.com/nbio/xx"
)

// Result represents an EPP <result> element.
type Result struct {
	Code    int    `xml:"code,attr"`
	Message string `xml:"msg"`
}

// IsError determines whether an EPP status code is an error.
// https://tools.ietf.org/html/rfc5730#section-3
func (r *Result) IsError() bool {
	return r.Code >= 2000
}

// IsFatal determines whether an EPP status code is a fatal response,
// and the connection should be closed.
// https://tools.ietf.org/html/rfc5730#section-3
func (r *Result) IsFatal() bool {
	return r.Code >= 2500
}

// Error implements the error interface.
func (r *Result) Error() string {
	return fmt.Sprintf("EPP result code %d: %s", r.Code, r.Message)
}

func init() {
	path := "epp>response>result"
	scanResponse.MustHandleStartElement(path, func(c *xx.Context) error {
		res := c.Value.(*response_)
		res.Result.Code = c.AttrInt("", "code")
		return nil
	})
	scanResponse.MustHandleCharData(path+">msg", func(c *xx.Context) error {
		c.Value.(*response_).Result.Message = string(c.CharData)
		return nil
	})
}
