package option

import (
	"io"
	"strings"

	"github.com/domainr/epp/internal/schema/raw"
	"github.com/nbio/xml"
)

func Values(names map[uint64]string) map[string]uint64 {
	values := make(map[string]uint64, len(names))
	for value, name := range names {
		values[name] = value
	}
	return values
}

func Encode(e *xml.Encoder, start xml.StartElement, v uint64, names map[uint64]string) error {
	var b strings.Builder
	for i := uint64(0); i < 64; i++ {
		j := uint64(1) << i
		if v&j != 0 {
			b.WriteByte('<')
			b.WriteString(names[j])
			b.WriteString("/>")
		}
	}
	return e.EncodeElement(&raw.XML{Value: b.String()}, start)
}

func Decode(v *uint64, d *xml.Decoder, values map[string]uint64) error {
	for {
		var e struct {
			XMLName xml.Name
		}
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		*v |= values[e.XMLName.Local]
	}
	return nil
}
