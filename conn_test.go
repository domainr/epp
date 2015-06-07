package epp

import (
	"crypto/tls"
	"net"
	"testing"

	"github.com/nbio/st"
)

const (
	// https://wiki.hexonet.net/wiki/Domain_API
	addr     = "api.1api.net:1700"
	user     = "test.user"
	password = "test.passw0rd"
)

func dial(t *testing.T) net.Conn {
	conn, err := tls.Dial("tcp", addr, nil)
	st.Assert(t, err, nil)
	return conn
}

func login(t *testing.T) *Conn {
	c, err := NewConn(dial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
	return c
}

func TestNewConn(t *testing.T) {
	c, err := NewConn(dial(t))
	st.Expect(t, err, nil)
	st.Reject(t, c, nil)
	st.Reject(t, c.Greeting, nil)
}

func TestConnID(t *testing.T) {
	var c Conn
	st.Expect(t, c.id(), "0000000000000001")
	st.Expect(t, c.id(), "0000000000000002")
	st.Expect(t, c.id(), "0000000000000003")
	st.Expect(t, c.id(), "0000000000000004")
	st.Expect(t, c.id(), "0000000000000005")
	st.Expect(t, c.id(), "0000000000000006")
	st.Expect(t, c.id(), "0000000000000007")
	st.Expect(t, c.id(), "0000000000000008")
	st.Expect(t, c.id(), "0000000000000009")
	st.Expect(t, c.id(), "000000000000000a")
}
