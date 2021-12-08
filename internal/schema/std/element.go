package std

import "github.com/nbio/xml"

// Element is a generic XML element, used for marshaling
// other types into an XML wrapper.
type Element struct {
	XMLName  xml.Name
	Contents []interface{}
}
