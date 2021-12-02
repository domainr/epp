package std

import (
	"time"
)

// Time represents W3C XML dateTime values.
// See https://www.w3.org/TR/xmlschema-2/#dateTime and https://www.rfc-editor.org/rfc/rfc3339.html.
type Time struct {
	time.Time
}

// NewTime returns a pointer to a DateTime struct.
func NewTime(t time.Time) *Time {
	return &Time{t}
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