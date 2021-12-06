package epp

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/domainr/epp/internal/schema/epp"
	"github.com/nbio/xml"
)

// Client represents an EPP Client.
type Client struct {
	t Transport

	cfg *Config

	newID func() string

	// hasGreeting is closed when the initial <greeting> is received from the server and stored.
	hasGreeting chan struct{}
	greeting    atomic.Value

	m            sync.Mutex
	transactions map[string]transaction

	// done is closed when the client receives a fatal error or the connection is closed.
	done chan struct{}
}

type transaction struct {
	ctx context.Context
	c   chan resErr
}

type resErr struct {
	res *epp.Response
	err error
}

// NewClient returns an EPP client from t and cfg.
// A Transport can be created from an io.Reader/Writer pair or a net.Conn,
// typically a tls.Conn. Once used, cfg should not be modified.
func NewClient(t Transport, cfg *Config) (*Client, error) {
	return newClient(t, cfg)
}

func newClient(t Transport, cfg *Config) (*Client, error) {
	c := &Client{
		t:            t,
		cfg:          cfg,
		newID:        cfg.TransactionID,
		hasGreeting:  make(chan struct{}),
		transactions: make(map[string]transaction),
		done:         make(chan struct{}),
	}
	if c.newID == nil {
		src, err := newSeqSource("")
		if err != nil {
			return nil, err
		}
		c.newID = src.ID
	}
	go c.readLoop()
	return c, nil
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

// newCommand returns a new epp.Command with a unique transaction ID.
func (c *Client) newCommand() *epp.Command {
	return &epp.Command{
		ClientTransactionID: c.newID(),
	}
}

// readLoop reads EPP messages from c.t and sends them to c.responses.
// It closes c.responses before returning.
// I/O errors are considered fatal and are returned.
func (c *Client) readLoop() {
	var err error
	defer func() {
		c.closeTransactions(err)
	}()
	for {
		select {
		case <-c.done:
			return
		default:
		}

		var data []byte
		data, err = c.t.ReadDataUnit()
		if err != nil {
			// TODO: log I/O errors.
			return
		}

		var e epp.EPP
		err = xml.Unmarshal(data, &e)
		if err != nil {
			// TODO: log XML parsing errors.
			// TODO: should XML parsing errors be considered fatal?
			continue
		}

		// TODO: this is not exactly conforming, as a valid <epp> message
		// should not contain both a <greeting> and a <response> element.
		if e.Greeting != nil {
			c.setGreeting(e.Greeting)
		}
		if e.Response != nil {
			c.processResponse(e.Response)
		}

		// TODO: log if server receives a <hello> or <command>.
	}
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

func (c *Client) processResponse(r *epp.Response) {
	if r.TransactionID.Client == "" {
		// TODO: log when server responds with an empty client transaction ID.
		return
	}
	t, ok := c.popTransaction(r.TransactionID.Client)
	if !ok {
		// TODO: log when server responds with unknown transaction ID.
		// TODO: keep abandoned transactions around for some period of time.
		return
	}
	select {
	case <-t.ctx.Done():
		// TODO: log context cancellation error.
		// TODO: keep abandoned transactions around for some period of time.
	case t.c <- resErr{res: r}:
	}
}

func (c *Client) closeTransactions(err error) {
	c.m.Lock()
	transactions := c.transactions
	c.transactions = nil
	c.m.Unlock()
	for _, t := range transactions {
		select {
		case <-t.ctx.Done():
			// TODO: log context cancellation error.
			// TODO: keep abandoned transactions around for some period of time.
		case t.c <- resErr{err: err}:
		}
	}
}

func (c *Client) pushTransaction(id string, t transaction) error {
	c.m.Lock()
	defer c.m.Unlock()
	_, ok := c.transactions[id]
	if ok {
		return fmt.Errorf("epp: transaction already exists: %s", id)
	}
	c.transactions[id] = t
	return nil
}

func (c *Client) popTransaction(id string) (transaction, bool) {
	c.m.Lock()
	defer c.m.Unlock()
	t, ok := c.transactions[id]
	if ok {
		delete(c.transactions, id)
	}
	return t, ok
}
