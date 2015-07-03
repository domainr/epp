package epp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
)

// ScanFunc is a callback that accepts an xml.StartElement, an
// xml.CharData, and an optional interface{} value for private use.
//
// The function is called for two XML tokens: xml.StartElement
// and xml.CharData. The xml.StartElement will always be the last
// StartElement parsed. CharData may be nil.
type ScanFunc func(xml.StartElement, xml.CharData, interface{}) error

// Scanner scans XML from an xml.Decoder, looking for specific paths.
type Scanner struct {
	tree map[xml.Name]*Scanner
	f    ScanFunc
}

// NewScanner returns an initialized Scanner, ready to use.
func NewScanner() *Scanner {
	return &Scanner{tree: make(map[xml.Name]*Scanner)}
}

// Handle sets up an XML path handler for a Scanner.
//
// Paths must be in the form of "foo>bar" or "uri foo>uri bar".
// The path is delimited by > characters, and individual path
// elements are split on whitespace into a namespace and local
// tag name.
//
// Handle returns ErrInvalidPath if the specified path
// is malformed.
func (s *Scanner) Handle(path string, f ScanFunc) error {
	names := strings.SplitN(path, ">", 2)
	fields := strings.Fields(names[0])
	var name xml.Name
	switch len(fields) {
	case 0:
		if len(names) > 1 {
			return ErrInvalidPath
		}
		s.f = f
		return nil
	case 1:
		name.Local = fields[0]
	case 2:
		name.Space = fields[0]
		name.Local = fields[1]
	default:
		return ErrInvalidPath
	}
	s2, ok := s.tree[name]
	if !ok {
		s2 = NewScanner()
		s.tree[name] = s2
	}
	if len(names) == 1 {
		return s2.Handle("", f)
	}
	return s2.Handle(names[1], f)
}

// ErrInvalidPath describes an invalid Scanner path
// returned by Scanner.Handle.
var ErrInvalidPath = errors.New("invalid scan path")

// Scan reads xml.Tokens from Decoder d, passing matching
// xml.StartElement and xml.CharData tokens to the matching
// ScanFuncs in s. It returns any errors encountered.
// It will return if it encounters an xml.EndElement that matches
// the corresponding xml.StartElement it attempted to scan.
func (s *Scanner) Scan(d *xml.Decoder, v interface{}) error {
	var e xml.StartElement
	for {
		t, err := d.Token()
		if err != nil {
			return err
		}
		switch node := t.(type) {
		case xml.StartElement:
			e = node
			s2, ok := s.tree[node.Name]
			if !ok {
				s2, ok = s.tree[xml.Name{"", node.Name.Local}]
				if !ok {
					fmt.Printf("Skipping: %s\n", node.Name.Local)
					err = d.Skip()
					break
				}
			}
			if s2.f != nil {
				err = s2.f(e, nil, v)
				if err != nil {
					return err
				}
			}
			err = s2.Scan(d, v)

		case xml.EndElement:
			return nil

		case xml.CharData:
			if s.f != nil {
				err = s.f(e, node, v)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
