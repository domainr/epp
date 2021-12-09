package epp

// Create represents an EPP <create> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.3.1.
type Create struct {
	// TODO: DomainCreate *domain.Create
	// TODO: HostCreate *host.Create
}

func (Create) eppCommand() {}
