package date

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
// It fails silently if an error occurs.
func ParseDuration(s string) Duration {
	if p, err := period.Parse(s, false); err == nil {
		td, _ := p.Duration()
		return Duration{td}
	}
	return Duration{}
}

// MarshalText implements encoding.TextMarshaler.
func (d Duration) MarshalText() ([]byte, error) {
	p, _ := period.NewOf(d.Duration)
	return p.MarshalText()
}

// UnmarshalText implements an encoding.TextUnmarshaler that ignores parsing errors.
func (d *Duration) UnmarshalText(text []byte) error {
	*d = ParseDuration(string(text))
	return nil
}
