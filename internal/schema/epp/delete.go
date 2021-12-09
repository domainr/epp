package epp

// Delete represents an EPP <delete> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.3.1.
type Delete struct {
	// TODO: DomainDelete *domain.Delete
	// TODO: HostDelete *host.Delete
}

func (Delete) eppCommand() {}
