package std

// XMLValuer is the interface implemented by types that need to modify their
// representation before being marshaled into or from XML.
//
// XMLValue returns a value that can be marshaled to or unmarshaled from XML. It
// will be passed a single argument, which is the value being marshaled or
// unmarshaled. If an XMLValuer is embedded in another struct, XMLValue will be
// called with a pointer to the outer struct.
type XMLValuer interface {
	XMLValue(interface{}) interface{}
}

// XMLValue will return v.XMLValue(v) if v implements XMLValuer.
// If v does not implement XMLValuer, it will return v.
func XMLValue(v interface{}) interface{} {
	if v, ok := v.(XMLValuer); ok {
		return v.XMLValue(v)
	}
	return v
}
