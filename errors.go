package epp

import "errors"

// ErrClosedConnection is the error used for read or write operations on a closed connection.
var ErrClosedConnection = errors.New("epp: operation on closed connection")
