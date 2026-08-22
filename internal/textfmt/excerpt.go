package textfmt

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func Excerpt(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 {
		return ""
	}
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	runes := []rune(s)
	out := strings.TrimSpace(string(runes[:n]))
	return out + "..."
}

func WordCount(s string) int {
	n := 0
	in := false
	for _, r := range s {
		if unicode.IsSpace(r) {
			in = false
			continue
		}
		if !in {
			n++
			in = true
		}
	}
	return n
}

func Wrap(s string, width int) string {
	if width < 8 {
		width = 8
	}
	words := strings.Fields(s)
	if len(words) == 0 {
		return ""
	}
	var b strings.Builder
	line := 0
	for i, w := range words {
		need := len(w)
		if line > 0 {
			need++
		}
		if line > 0 && line+need > width {
			b.WriteByte('\n')
			b.WriteString(w)
			line = len(w)
			continue
		}
		if i > 0 && line > 0 {
			b.WriteByte(' ')
		}
		b.WriteString(w)
		line += need
	}
	return b.String()
}

func CollapseSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
