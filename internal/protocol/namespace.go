package protocol

import (
	"encoding"
	"encoding/xml"

	"github.com/domainr/epp/ns"
)

type baseNamespace struct{}

var (
	_ encoding.TextMarshaler = baseNamespace{}
	_ xml.UnmarshalerAttr    = baseNamespace{}
)

func (baseNamespace) MarshalText() ([]byte, error) {
	return nil, nil
}

func (baseNamespace) UnmarshalXMLAttr(attr xml.Attr) error {
	return nil
}

type eppNamespace struct{ baseNamespace }

func (eppNamespace) MarshalText() ([]byte, error) {
	return unsafeBytes(ns.EPP), nil
}

type domainNamespace struct{ baseNamespace }

func (domainNamespace) MarshalText() ([]byte, error) {
	return unsafeBytes(ns.Domain), nil
}
