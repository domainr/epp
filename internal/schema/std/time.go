package std

import (
	"time"
)

// Time represents an W3C XML date-time value.
// See https://www.w3.org/TR/xmlschema-2/#dateTime and https://www.rfc-editor.org/rfc/rfc3339.html.
type Time struct {
	time.Time
}

// ParseTime parses an RFC 3339 date-time string.
// It returns an empty value if unable to parse s.
func ParseTime(s string) Time {
	tt, _ := time.Parse(time.RFC3339, s)
	return Time{tt}
}

// Pointer returns a pointer to t, useful for declaring composite literals.
func (t Time) Pointer() *Time {
	return &t
}

// MarshalText implements encoding.TextMarshaler.
func (t *Time) MarshalText() ([]byte, error) {
	if t == nil {
		return nil, nil
	}
	return t.Time.MarshalText()
}

// UnmarshalText implements an encoding.TextUnmarshaler that ignores parsing errors.
func (t *Time) UnmarshalText(text []byte) error {
	_ = t.Time.UnmarshalText(text)
	return nil
}
