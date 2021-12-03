package epp

import "fmt"

// ResultCode represents a 4-digit EPP result code.
type ResultCode uint16

// MarshalText implements encoding.TextMarshaler to print c as a 4-digit number.
func (c ResultCode) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%04d", c)), nil
}
