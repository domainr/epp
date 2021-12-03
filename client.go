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

	greetingStored chan struct{}
	greetingFunc   atomic.Value
}

// NewClient returns an EPP client using t as transport.
// A Transport can be created from an io.Reader/Writer pair or a net.Conn,
// typically a tls.Conn.
func NewClient(t Transport) (Client, error) {
	return newClient(t)
}

func newClient(t Transport) (Client, error) {
	c := &client{
		t: t,
	}
	c.greetingFunc.Store(func(ctx context.Context) *epp.Greeting {
		select {
		case <-ctx.Done():
			return nil
		case <-c.greetingStored:
			return c.greetingFunc.Load().(greetingFunc)(ctx)
		}
	})
	return c, nil
}

func (c *client) storeGreeting(g *epp.Greeting) {
	c.greetingFunc.Store(func(context.Context) *epp.Greeting {
		return g
	})
	select {
	case <-c.greetingStored:
	default:
		close(c.greetingStored)
	}
}

func (c *client) greeting(ctx context.Context) *epp.Greeting {
	return c.greetingFunc.Load().(greetingFunc)(ctx)
}

type greetingFunc func(context.Context) *epp.Greeting
