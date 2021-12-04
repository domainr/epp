package epp

// Config describes an EPP client or server configuration, including
// EPP objects and extensions used for a connection.
type Config struct {
	// The name of the EPP server, sent in <greeting> messages.
	// Only used for by the server.
	ServerName string

	// BCP 47 language code(s) for human-readable messages supported by an EPP client or server.
	// For clients, this describes the desired language(s) in preferred order.
	// If nil, the client will attempt to select "en". If the server
	// does not support any of the client’s preferred languages, the first
	// language advertised by the server will be selected.
	// For servers, this describes its supported language(s).
	// If nil, []string{"en"} will be used.
	Languages []string

	// Namespace URIs of EPP objects supported by a client or server.
	// For clients, this describes the object type(s) the client wants to access.
	// For servers, this describes the object type(s) the server allows clients to access.
	// If nil, a reasonable set of defaults will be used.
	Objects []string

	// EPP extension URIs supported by a client or server.
	// For clients, this is a list of extensions(s) the client wants to use in preferred order.
	// If nil, a client will use the highest version of each supported extension advertised by the server.
	// For servers, this is an advertised list of supported extension(s).
	// If nil, a server will use a reasonable set of defaults.
	Extensions []string

	// EPP extension URIs that will be used by a client or server,
	// regardless of whether the peer advertises it.
	ForcedExtensions []string
}

// Copy returns a deep copy of c.
func (c Config) Copy() Config {
	c.Languages = copySlice(c.Languages)
	c.Objects = copySlice(c.Objects)
	c.Extensions = copySlice(c.Extensions)
	c.ForcedExtensions = copySlice(c.ForcedExtensions)
	return c
}

func copySlice(s []string) []string {
	if s == nil {
		return nil
	}
	dst := make([]string, len(s))
	copy(dst, s)
	return dst
}
