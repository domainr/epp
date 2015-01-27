package epp

import (
	"crypto/tls"
	"encoding/json"
	"sync"
)

// Limit concurrent EPP connections per Client.
const MaxConnections = 20

// Client allows concurrent access to an EPP server up to MaxConnections.
type Client struct {
	addr      string
	clientID  string
	password  string
	tlsConfig *tls.Config
	pool      *sync.Pool
	limiter   Limiter
}

// A Limiter acts as a semaphore to limit concurrency of a Client.
type Limiter chan struct{}

func NewLimiter(n int) Limiter { return make(chan struct{}, n) }
func (l Limiter) Start()       { l <- struct{}{} }
func (l Limiter) Done()        { <-l }

// NewClient returns a new EPP Client for addr, authenticated with clientID and password.
// The returned client is safe for concurrent use.
func NewClient(addr, clientID, password string, cfg *tls.Config) (*Client, error) {
	c := &Client{
		addr:      addr,
		clientID:  clientID,
		password:  password,
		tlsConfig: cfg,
		pool:      &sync.Pool{},
		limiter:   NewLimiter(MaxConnections),
	}
	return c, nil
}

// MarshalJSON is used by expvar for metrics.
func (c Client) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Addr": c.addr,
		"User": c.clientID,
		"Conn": len(c.limiter),
		"Pool": cap(c.limiter),
	})
}

// CheckDomain checks the domains.
func (c *Client) CheckDomain(domains ...string) (*DomainCheckResponse, error) {
	c.limiter.Start()
	defer c.limiter.Done()
	conn, err := c.getFreeConn()
	if err != nil {
		return nil, err
	}
	dcr, err := conn.CheckDomain(domains...)
	c.release(conn, err)
	return dcr, err
}

// Get an EPP connection or create one.
func (c *Client) getFreeConn() (*Conn, error) {
	conn, ok := c.pool.Get().(*Conn)
	if conn == nil || !ok {
		return c.dial()
	}
	return conn, nil
}

// Create and authenticate a new connection to the EPP server.
func (c *Client) dial() (conn *Conn, err error) {
	if c.tlsConfig != nil {
		conn, err = DialTLS(c.addr, c.tlsConfig)
	} else {
		conn, err = Dial(c.addr)
	}
	if err != nil {
		return nil, err
	}
	err = conn.Login(c.clientID, c.password, "")
	if err != nil {
		return nil, err
	}
	return
}

// Release conn to the connection pool.
func (c *Client) release(conn *Conn, err error) {
	if !resumable(err) {
		conn.Close()
		return
	}
	c.pool.Put(conn)
}

func resumable(err error) bool {
	if err == nil {
		return true
	}
	if r, ok := err.(*Result); ok && !r.IsFatal() {
		return true
	}
	return false
}
