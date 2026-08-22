package precopy

// ClonePrefix returns an independent prefix of src.
func ClonePrefix(src []byte, n int) []byte {
	if n < 0 {
		n = 0
	}
	if n > len(src) {
		n = len(src)
	}
	out := make([]byte, n)
	copy(out, src[:n])
	return out
}

// CloneAll returns a copy of src that does not share a backing array.
func CloneAll(src []byte) []byte {
	return ClonePrefix(src, len(src))
}

// HeadStrings copies the first n strings.
func HeadStrings(src []string, n int) []string {
	if n < 0 {
		n = 0
	}
	if n > len(src) {
		n = len(src)
	}
	out := make([]string, n)
	copy(out, src[:n])
	return out
}
