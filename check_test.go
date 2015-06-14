package epp

import (
	"os"
	"testing"

	"github.com/nbio/st"
)

func init() {
	DebugLogger = os.Stdout
}

func TestConnCheck(t *testing.T) {
	c := testLogin(t)

	dc, err := c.CheckDomain("google.com")
	st.Expect(t, err, nil)
	st.Reject(t, dc, nil)
	st.Expect(t, len(dc.Results), 1)
	st.Expect(t, dc.Results[0].Domain.Domain, "google.com")
	st.Expect(t, dc.Results[0].Domain.IsAvailable, false)

	dc, err = c.CheckDomain("dmnr-test-x759824vim-i2.com")
	st.Expect(t, err, nil)
	st.Reject(t, dc, nil)
	st.Expect(t, len(dc.Results), 1)
	st.Expect(t, dc.Results[0].Domain.Domain, "dmnr-test-x759824vim-i2.com")
	st.Expect(t, dc.Results[0].Domain.IsAvailable, true)

	dc, err = c.CheckDomain("--dmnr-test--.com")
	st.Reject(t, err, nil)
	st.Expect(t, dc, (*DomainCheckData)(nil))
}
