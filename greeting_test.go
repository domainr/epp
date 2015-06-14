package epp

import (
	"testing"

	"github.com/nbio/st"
)

func TestHello(t *testing.T) {
	c, err := NewConn(testDial(t))
	st.Assert(t, err, nil)
	g, err := c.Hello()
	st.Expect(t, err, nil)
	st.Reject(t, g, nil)
}
