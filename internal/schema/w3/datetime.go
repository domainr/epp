package w3

import (
	"time"
)

// DateTime represents W3C XML dateTime values.
// See https://www.w3.org/TR/xmlschema-2/#dateTime and https://www.rfc-editor.org/rfc/rfc3339.html.
type DateTime struct {
	time.Time
}

// NewDateTime returns a pointer to a DateTime struct.
func NewDateTime(t time.Time) *DateTime {
	return &DateTime{t}
}

// UnmarshalText implements an encoding.TextUnmarshaler that ignores parsing errors.
func (t *DateTime) UnmarshalText(text []byte) error {
	_ = t.Time.UnmarshalText(text)
	return nil
}
