package raw

// XML is a container for raw XML. Useful for single, self-closing tags, e.g.:
// <container><value/><container>.
type XML struct {
	Value string `xml:",innerxml"`
}
