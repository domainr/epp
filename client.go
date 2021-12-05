package epp

import (
	"context"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/nbio/xml"
)

// Client represents an EPP Client.
type Client struct {
	t Transport

	cfg *Config

	// hasGreeting is closed when the initial <greeting> is received from the server and stored.
	hasGreeting chan struct{}
	greeting    atomic.Value

	responses chan (response)

	// done is closed when the client receives a fatal error or the connection is closed.
	done chan struct{}
}

// NewClient returns an EPP client from t and cfg.
// A Transport can be created from an io.Reader/Writer pair or a net.Conn,
// typically a tls.Conn. Once used, cfg should not be modified.
func NewClient(t Transport, cfg *Config) (*Client, error) {
	return newClient(t, cfg)
}

func newClient(t Transport, cfg *Config) (*Client, error) {
	c := &Client{
		t:           t,
		cfg:         cfg,
		hasGreeting: make(chan struct{}),
		done:        make(chan struct{}),
	}
	return c, nil
}

func (c *Client) setGreeting(g *epp.Greeting) {
	c.greeting.Store(g)
	select {
	case <-c.hasGreeting:
	default:
		close(c.hasGreeting)
	}
}

func (c *Client) waitForGreeting(ctx context.Context) (*epp.Greeting, error) {
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

// ServerConfig returns the server configuration described in a <greeting> message.
// Will block until the an initial <greeting> is received, or ctx is canceled.
func (c *Client) ServerConfig(ctx context.Context) (*Config, error) {
	g, err := c.waitForGreeting(ctx)
	if err != nil {
		return nil, err
	}
	return configFromGreeting(g), nil
}

// ServerName returns the most recently received server name.
// Will block until an initial <greeting> is received, or ctx is canceled.
func (c *Client) ServerName(ctx context.Context) (string, error) {
	g, err := c.waitForGreeting(ctx)
	if err != nil {
		return "", err
	}
	return g.ServerName, nil
}

// ServerTime returns the most recently received timestamp from the server.
// Will block until an initial <greeting> is received, or ctx is canceled.
// TODO: what is used for?
func (c *Client) ServerTime(ctx context.Context) (time.Time, error) {
	g, err := c.waitForGreeting(ctx)
	if err != nil {
		return time.Time{}, err
	}
	return g.ServerDate.Time, nil
}

// readLoop reads EPP messages from c.t and sends them to c.responses.
// It closes c.responses before returning.
// I/O errors are considered fatal and are returned.
func (c *Client) readLoop() {
	defer close(c.responses)
	for {
		select {
		case <-c.done:
			return
		default:
		}

		data, err := c.t.ReadDataUnit()
		if err != nil {
			c.responses <- response{err: err}
			return
		}

		var e epp.EPP
		err = xml.Unmarshal(data, &e)
		if err != nil {
			c.responses <- response{err: err}
			continue // TODO: should XML parsing errors be considered fatal?
		}

		// TODO: this is not exactly conforming, as a valid <epp> message
		// should not contain both a <greeting> and a <response> element.
		if e.Greeting != nil {
			c.setGreeting(e.Greeting)
		}
		if e.Response != nil {
			c.responses <- response{res: e.Response}
		}

		// TODO: log if server receives a <hello> or <command>.
	}
}

type response struct {
	res *epp.Response
	err error
}

// Close closes a client connection. If the Transport used by this client
// implements io.Closer, Close will be called and any error returned.
func (c *Client) Close() error {
	select {
	case <-c.done:
		return net.ErrClosed
	default:
		close(c.done)
	}
	if closer, ok := c.t.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
