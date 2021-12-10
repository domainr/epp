package epp

import (
	"context"

	"github.com/domainr/epp/internal/schema/epp"
)

// Server is an EPP version 1.0 server.
type Server struct {
	// Name is the name of this EPP server. It is sent to clients in a EPP
	// <greeting> message. If empty, a reasonable default will be used.
	Name string

	// Config describes the EPP server configuration. Configuration
	// parameters are announced to EPP clients in an EPP <greeting> message.
	Config Config

	// Handler is called in a goroutine for each incoming EPP connection.
	// The connection will be closed when Handler returns.
	Handler func(Session) error
}

type Session interface {
	// Context returns the connection Context for this session. The Context
	// will be canceled if the underlying Transport goes away or is closed.
	Context() context.Context

	// ReadCommand reads the next EPP command from the client. An error will
	// be returned if the underlying connection is closed or an error occurs
	// reading from the connection.
	ReadCommand() (*epp.Command, error)

	// WriteResponse sends an EPP response to the client. An error will
	// be returned if the underlying connection is closed or an error occurs
	// writing to the connection.
	WriteResponse(*epp.Response) error
}
