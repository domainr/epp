package epp

import (
	"os"
	"testing"

	"github.com/nbio/st"
)

func init() {
	DebugLogger = os.Stdout
}

func TestConnLogin(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	err = c.Login(user, password, "")
	st.Assert(t, err, nil)
}
