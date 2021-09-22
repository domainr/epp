package protocol

import (
	"testing"
)

func TestUnsafeBytes(t *testing.T) {
	tests := []string{
		"",
		"foo",
		"bar",
		"baz",
		"héllo",
	}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			got := unsafeBytes(tt)
			want := []byte(tt)
			if !cmp(got, want) {
				t.Errorf("unsafeBytes() = %v, want %v", got, want)
			}
		})
	}
}

func cmp(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
