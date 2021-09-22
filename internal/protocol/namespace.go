package protocol

import (
	"encoding"

	"github.com/domainr/epp/ns"
)

type eppNamespace struct{}

var _ encoding.TextMarshaler = eppNamespace{}

func (eppNamespace) MarshalText() ([]byte, error) {
	return unsafeBytes(ns.EPP), nil
}

type domainNamespace struct{}

var _ encoding.TextMarshaler = domainNamespace{}

func (domainNamespace) MarshalText() ([]byte, error) {
	return unsafeBytes(ns.Domain), nil
}
