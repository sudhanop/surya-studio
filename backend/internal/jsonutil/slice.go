package jsonutil

// Slice ensures JSON encodes an empty array instead of null when no rows exist.
func Slice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}
