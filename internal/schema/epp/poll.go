package epp

// Poll represents an EPP <poll> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.3.
type Poll struct {
}

func (Poll) eppCommand() {}
