package std

import (
	"time"

	"github.com/rickb777/date/period"
)

// Duration represents W3C XML duration values.
// See https://www.w3.org/TR/xmlschema-2/#duration and https://www.rfc-editor.org/rfc/rfc3339.html.
type Duration struct {
	time.Duration
}

// ParseDuration parses an RFC 3339 duration string.
// It returns an empty value if unable to parse s.
func ParseDuration(s string) Duration {
	p, _ := period.Parse(s)
	d, _ := p.Duration()
	return Duration{d}
}

// Pointer returns a pointer to d, useful for declaring composite literals.
func (d Duration) Pointer() *Duration {
	return &d
}

// MarshalText implements encoding.TextMarshaler.
func (d *Duration) MarshalText() ([]byte, error) {
	p, _ := period.NewOf(d.Duration)
	return p.MarshalText()
}

// UnmarshalText implements an encoding.TextUnmarshaler that ignores parsing errors.
func (d *Duration) UnmarshalText(text []byte) error {
	*d = ParseDuration(string(text))
	return nil
}
