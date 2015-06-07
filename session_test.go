package epp

import (
	"testing"

	"github.com/nbio/st"
)

func TestConnLogin(t *testing.T) {
	c, err := NewConn(dial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
}
