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
	// mWrite protects writes on t.
	mWrite sync.Mutex
	t      Transport

	cfg *Config

	// newID is a function that generates unique client transaction IDs.
	// It is copied from cfg, and set to a reasonable default if nil.
	newID func() string

	// greeting stores the most recently received <greeting> from the server.
	greeting atomic.Value

	// hasGreeting is closed when the client receives an initial <greeting> from the server.
	hasGreeting chan struct{}

	mHellos sync.Mutex
	hellos  []transaction

	mCommands sync.Mutex
	commands  map[string]transaction

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
		newID:       cfg.TransactionID,
		hasGreeting: make(chan struct{}),
		commands:    make(map[string]transaction),
		done:        make(chan struct{}),
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

// Hello sends an EPP <hello> message to the server.
// It will block until the next <greeting> message is received or ctx is canceled.
func (c *Client) Hello(ctx context.Context) error {
	t, cancel := newTransaction(ctx)
	defer cancel()
	c.pushHello(t)

	err := c.write(&epp.EPP{Body: &epp.Hello{}})
	if err != nil {
		return err
	}

	select {
	case <-c.done:
		return ErrClosedConnection
	case <-ctx.Done():
		return ctx.Err()
	case reply := <-t.reply:
		// TODO: do something with the greeting
		return reply.err
	}
}

// newCommand returns a new epp.Command with a unique transaction ID.
func (c *Client) newCommand() *epp.Command {
	return &epp.Command{
		ClientTransactionID: c.newID(),
	}
}

// write writes e to the underlying Transport.
// Writes are synchronized, so it is safe to call this from multiple goroutines.
func (c *Client) write(e *epp.EPP) error {
	x, err := xml.Marshal(&e)
	if err != nil {
		return err
	}
	return c.writeDataUnit(x)
}

// writeDataUnit writes a single EPP data unit to the underlying Transport.
// Writes are synchronized, so it is safe to call this from multiple goroutines.
func (c *Client) writeDataUnit(p []byte) error {
	c.mWrite.Lock()
	defer c.mWrite.Unlock()
	return c.t.WriteDataUnit(p)
}

// readLoop reads EPP messages from c.t and sends them to c.responses.
// It closes c.responses before returning.
// I/O errors are considered fatal and are returned.
func (c *Client) readLoop() {
	var err error
	defer func() {
		c.cleanup(err)
	}()
	for {
		select {
		case <-c.done:
			return
		default:
		}

		var p []byte
		p, err = c.t.ReadDataUnit()
		if err != nil {
			// TODO: log I/O errors.
			return
		}

		err = c.handleDataUnit(p)
		if err != nil {
			// TODO: log XML and processing errors.
		}
	}
}

func (c *Client) handleDataUnit(p []byte) error {
	e := &epp.EPP{}
	err := xml.Unmarshal(p, e)
	if err != nil {
		// TODO: log XML parsing errors.
		// TODO: should XML parsing errors be considered fatal?
		return err
	}

	// TODO: log processing errors.
	return c.handleReply(e)
}

func (c *Client) handleReply(e *epp.EPP) error {
	// TODO: this is not exactly conforming, as a valid <epp> message
	// should not contain both a <greeting> and a <response> element.

	switch body := e.Body.(type) {
	case *epp.Response:
		id := body.TransactionID.Client
		if id == "" {
			// TODO: log when server responds with an empty client transaction ID.
			return TransactionIDError(id)
		}
		t, ok := c.popCommand(id)
		if !ok {
			// TODO: log when server responds with unknown transaction ID.
			// TODO: keep abandoned transactions around for some period of time.
			return TransactionIDError(id)
		}
		err := c.finalize(t, e, nil)
		if err != nil {
			return err
		}

	case *epp.Greeting:
		// Always store the last <greeting> received from the server.
		c.greeting.Store(body)

		// Close hasGreeting this is the first <greeting> recieved.
		select {
		case <-c.hasGreeting:
		default:
			close(c.hasGreeting)
		}

		// Pass the <greeting> to a caller waiting on it.
		t, ok := c.popHello()
		if ok {
			err := c.finalize(t, e, nil)
			if err != nil {
				return err
			}
		}

	case *epp.Hello:
		// TODO: log if server receives a <hello> or <command>.

	case *epp.Command:
		// TODO: log if server receives a <hello> or <command>.
	}

	return nil
}

func (c *Client) finalize(t transaction, e *epp.EPP, err error) error {
	select {
	case <-c.done:
		return ErrClosedConnection
	case <-t.ctx.Done():
		return t.ctx.Err()
	case t.reply <- reply{e: e, err: err}:
	}
	return nil
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

// pushHello adds a <hello> transaction to the end of the stack.
func (c *Client) pushHello(t transaction) {
	c.mHellos.Lock()
	defer c.mHellos.Unlock()
	c.hellos = append(c.hellos, t)
}

// popHello pops the oldest <hello> transaction off the front of the stack.
func (c *Client) popHello() (transaction, bool) {
	c.mHellos.Lock()
	defer c.mHellos.Unlock()
	if len(c.hellos) == 0 {
		return transaction{}, false
	}
	t := c.hellos[0]
	c.hellos = c.hellos[1:]
	return t, true
}

// pushCommand adds a <command> transaction to the map of in-flight commands.
func (c *Client) pushCommand(id string, t transaction) error {
	c.mCommands.Lock()
	defer c.mCommands.Unlock()
	_, ok := c.commands[id]
	if ok {
		return fmt.Errorf("epp: transaction already exists: %s", id)
	}
	c.commands[id] = t
	return nil
}

// popCommand removes a <command> transaction from the map of in-flight commands.
func (c *Client) popCommand(id string) (transaction, bool) {
	c.mCommands.Lock()
	defer c.mCommands.Unlock()
	t, ok := c.commands[id]
	if ok {
		delete(c.commands, id)
	}
	return t, ok
}

// cleanup cleans up and responds to all in-flight <hello> and <command> transactions.
// Each transaction will be finalized with err, which may be nil.
func (c *Client) cleanup(err error) {
	c.mHellos.Lock()
	hellos := c.hellos
	c.hellos = nil
	c.mHellos.Unlock()
	for _, t := range hellos {
		c.finalize(t, nil, err)
	}

	c.mCommands.Lock()
	commands := c.commands
	c.commands = nil
	c.mCommands.Unlock()
	for _, t := range commands {
		c.finalize(t, nil, err)
	}
}
