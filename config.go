package epp

import "github.com/domainr/epp/internal/schema/epp"

// Config describes an EPP client or server configuration, including
// EPP objects and extensions used for a connection.
type Config struct {
	// Supported EPP version(s). Typically this should not be set
	// by either a client or server. If nil, this will default to
	// []string{"1.0"} (currently the only supported version).
	Versions []string

	// BCP 47 language code(s) for human-readable messages.
	// For clients, this describes the desired language(s) in preferred order.
	// If the server does not support any of the client’s preferred languages,
	// the first language advertised by the server will be selected.
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
	// whether or not the peer indicates support for it.
	// This will always be nil when a configuration is read from a peer.
	ForcedExtensions []string
}

func configFromGreeting(g *epp.Greeting) *Config {
	c := &Config{}
	// TODO: should epp.Greeting have getter and setter methods to access deeply-nested data?
	if g.ServiceMenu != nil {
		c.Versions = copySlice(g.ServiceMenu.Versions)
		c.Objects = copySlice(g.ServiceMenu.Objects)
		if g.ServiceMenu.ServiceExtension != nil {
			c.Extensions = copySlice(g.ServiceMenu.ServiceExtension.Extensions)
		}
	}
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
