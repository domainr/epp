package date

import (
	"time"

	"github.com/rickb777/date/period"
)

// Duration represents RFC 3339 duration values.
// See https://www.rfc-editor.org/rfc/rfc3339.html.
type Duration time.Duration

// ParseDuration parses an RFC 3339 duration string.
// It fails silently if an error occurs.
func ParseDuration(s string) Duration {
	var d Duration
	if p, err := period.Parse(s, false); err == nil {
		td, _ := p.Duration()
		d = Duration(td)
	}
	return d
}

// MarshalText implements encoding.TextMarshaler.
func (v Duration) MarshalText() ([]byte, error) {
	p, _ := period.NewOf(time.Duration(v))
	return p.MarshalText()
}

// UnmarshalText implements a custom TextUnmarshaler that ignores parsing errors.
func (v *Duration) UnmarshalText(text []byte) error {
	*v = ParseDuration(string(text))
	return nil
}
