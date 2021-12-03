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

	// done is closed when the client receives a fatal error or the connection is closed.
	done chan struct{}

	// hasGreeting is closed when the initial <greeting> is received from the server and stored.
	hasGreeting chan struct{}
	greeting    atomic.Value
}

// NewClient returns an EPP client using t as transport.
// A Transport can be created from an io.Reader/Writer pair or a net.Conn,
// typically a tls.Conn.
func NewClient(t Transport) (Client, error) {
	return newClient(t)
}

func newClient(t Transport) (Client, error) {
	c := &client{
		t:           t,
		done:        make(chan struct{}),
		hasGreeting: make(chan struct{}),
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

func (c *client) waitForGreeting(ctx context.Context) *epp.Greeting {
	g := c.greeting.Load()
	if g != nil {
		return g.(*epp.Greeting)
	}
	select {
	case <-ctx.Done():
		return nil
	case <-c.hasGreeting:
		return c.greeting.Load().(*epp.Greeting)
	}
}
