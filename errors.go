package epp

import "errors"

// ErrClosedConnection indicates a read or write operation on a closed connection.
var ErrClosedConnection = errors.New("epp: operation on closed connection")

// ErrUnexpectedHello indicates an EPP message contained an unexpected <hello> element.
var ErrUnexpectedHello = errors.New("epp: unexpected <hello>")

// ErrUnexpectedCommand indicates an EPP message contained a <command> element.
var ErrUnexpectedCommand = errors.New("epp: unexpected <command>")

// ErrNoResponse indicates an EPP message did not contain an expected <response> element.
var ErrNoResponse = errors.New("epp: missing <response>")

// ErrUnexpectedResponse indicates an EPP message contained an unexpected <response> element.
var ErrUnexpectedResponse = errors.New("epp: unexpected <response>")
