package epp

import (
	"context"
	"sync/atomic"

	"github.com/domainr/epp/internal/schema/epp"
)

// Client represents an EPP client.
type Client interface {
}

type client struct {
	t Transport

	cfg *Config

	// hasGreeting is closed when the initial <greeting> is received from the server and stored.
	hasGreeting chan struct{}
	greeting    atomic.Value

	// done is closed when the client receives a fatal error or the connection is closed.
	done chan struct{}
}

// NewClient returns an EPP client from t and cfg.
// A Transport can be created from an io.Reader/Writer pair or a net.Conn,
// typically a tls.Conn. Once used, cfg should not be modified.
func NewClient(t Transport, cfg *Config) (Client, error) {
	return newClient(t, cfg)
}

func newClient(t Transport, cfg *Config) (Client, error) {
	c := &client{
		t:           t,
		cfg:         cfg,
		hasGreeting: make(chan struct{}),
		done:        make(chan struct{}),
	}
	return c, nil
}

func (c *client) setGreeting(g *epp.Greeting) {
	c.greeting.Store(g)
	select {
	case <-c.hasGreeting:
	default:
		close(c.hasGreeting)
	}
}

func (c *client) waitForGreeting(ctx context.Context) (*epp.Greeting, error) {
	g := c.greeting.Load()
	if g != nil {
		return g.(*epp.Greeting), nil
	}
	select {
	case <-c.done:
		return nil, ErrClosedConnection
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.hasGreeting:
		return c.greeting.Load().(*epp.Greeting), nil
	}
}
