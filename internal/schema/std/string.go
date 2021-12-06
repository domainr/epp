package std

// StringPointer returns a pointer to s.
// Used for declaring a pointer to a string literal.
func StringPointer(s string) *string {
	return &s
}
