package magic

import (
	"errors"
	"io"
	"reflect"
	"sync"

	"github.com/nbio/xml"
)

// A Decoder maps xml.Name values to Go types. This allows decoding XML into a
// Go struct with one or more interface{} fields, which would otherwise be skipped.
type Decoder struct {
	Types sync.Map
}

func (d *Decoder) RegisterType(name xml.Name, v interface{}) error {
	t := reflect.TypeOf(v)
	if t == nil {
		return ErrNil
	}
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t == nil {
			return ErrNil
		}
	}
	d.Types.Store(name, t)
	return nil
}

var ErrNil = errors.New("magic: nil type")

// DecodeElement attempts to decode the start element using its internal map of xml.Name to reflect.Type.
// It will silently skip unknown tags and return any XML parsing errors encountered.
func (d *Decoder) DecodeElement(xd *xml.Decoder, start *xml.StartElement) (interface{}, error) {
	t, ok := d.Types.Load(start.Name)
	if !ok {
		// Silently skip unknown tags.
		return nil, nil
	}
	v := reflect.New(t.(reflect.Type)).Interface()
	err := xd.DecodeElement(v, start)
	return v, err
}

// DecodeChildren attempts to decode the immediate child elements of start.
// It only evaluates start elements and ignores unknown tags.
func (d *Decoder) DecodeChildren(xd *xml.Decoder, start *xml.StartElement) ([]interface{}, error) {
	var values []interface{}
	for {
		t, err := xd.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return values, err
		}
		if start, ok := t.(xml.StartElement); ok {
			v, err := d.DecodeElement(xd, &start)
			if err != nil {
				return values, err
			}
			values = append(values, v)
		}
	}
	return values, nil
}
