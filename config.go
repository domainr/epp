package epp

// Config describes an EPP client or server configuration, including
// EPP objects and extensions used for a connection.
type Config struct {
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
	return Config{
		Objects:          copySlice(c.Objects),
		Extensions:       copySlice(c.Extensions),
		ForcedExtensions: copySlice(c.ForcedExtensions),
	}
}

func copySlice(s []string) []string {
	if s == nil {
		return nil
	}
	dst := make([]string, len(s))
	copy(dst, s)
	return dst
}
