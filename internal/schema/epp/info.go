package epp

// Info represents an EPP <info> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.2.
type Info struct {
	// TODO: DomainInfo *domain.Info
}

func (Info) eppCommand() {}
