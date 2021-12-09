package epp

// Transfer represents an EPP <transfer> command as defined in RFC 5730.
// See https://www.rfc-editor.org/rfc/rfc5730.html#section-2.9.2.4.
type Transfer struct {
	// TODO: DomainTransfer *domain.Transfer
}

func (Transfer) eppCommand() {}
