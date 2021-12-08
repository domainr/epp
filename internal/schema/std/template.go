package std

import (
	"io"
	"reflect"
	"sync"

	"github.com/nbio/xml"
)

// Template maps xml.Name values to Go types. This allows decoding XML into a
// Go struct with one or more interface{} fields, which would otherwise be skipped.
type Template struct {
	types sync.Map
}

// Add maps name to template value v, which must be a pointer to a concrete type.
// If v is nil or points to a nil type, Add will silently fail.
func (d *Template) Add(name xml.Name, v interface{}) {
	t := reflect.TypeOf(v)
	if t == nil {
		return
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t == nil {
			return
		}
	}
	d.types.Store(name, t)
}

// Type returns a reflect.Type for name.
// Returns nil if name is not mapped.
func (t *Template) Type(name xml.Name) reflect.Type {
	v, ok := t.types.Load(name)
	if !ok {
		return nil
	}
	return v.(reflect.Type)
}

// New returns a new instance of the type that matches xml.Name.
// Returns nil if the name does not have a type associated with it.
func (t *Template) New(name xml.Name) interface{} {
	typ := t.Type(name)
	if typ == nil {
		return nil
	}
	return reflect.New(typ).Interface()
}

// DecodeElement attempts to decode the start element using its internal map of xml.Name to reflect.Type.
// It will silently skip unknown tags and return any XML parsing errors encountered.
func (t *Template) DecodeElement(xd *xml.Decoder, start *xml.StartElement) (interface{}, error) {
	v := t.New(start.Name)
	if v == nil {
		// Silently skip unknown tags.
		return nil, nil
	}
	err := xd.DecodeElement(v, start)
	return v, err
}

// DecodeChildren attempts to decode the immediate child elements of start.
// It only evaluates start elements and ignores unknown tags.
func (t *Template) DecodeChildren(d *xml.Decoder, start *xml.StartElement) ([]interface{}, error) {
	var values []interface{}
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return values, err
		}
		if start, ok := tok.(xml.StartElement); ok {
			v, err := t.DecodeElement(d, &start)
			if err != nil {
				return values, err
			}
			values = append(values, v)
		}
	}
	return values, nil
}
