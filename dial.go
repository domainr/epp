package epp

import (
	"context"
	"crypto/tls"
)

// Dial opens a TLS connection to the EPP server at addr.
func Dial(ctx context.Context, addr string, cfg *tls.Config) (Client, error) {
	return nil, nil
}
