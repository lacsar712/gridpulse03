package textfmt

import (
	"strings"
	"unicode"
)

func TitleCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	prevSpace := true
	for _, r := range strings.ToLower(s) {
		if unicode.IsSpace(r) {
			if !prevSpace {
				b.WriteByte(' ')
			}
			prevSpace = true
			continue
		}
		if prevSpace {
			b.WriteRune(unicode.ToTitle(r))
		} else {
			b.WriteRune(r)
		}
		prevSpace = false
	}
	return b.String()
}

func PadRight(s string, n int) string {
	if n <= len(s) {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func PadLeft(s string, n int) string {
	if n <= len(s) {
		return s
	}
	return strings.Repeat(" ", n-len(s)) + s
}

func RepeatJoin(parts []string, sep string) string {
	return strings.Join(parts, sep)
}

func TruncateBytes(s string, n int) string {
	if n < 0 {
		n = 0
	}
	if len(s) <= n {
		return s
	}
	return s[:n]
}
