package textfmt

import (
	"strings"
	"unicode"
)

func Heading(level int, title string) string {
	if level < 1 {
		level = 1
	}
	if level > 6 {
		level = 6
	}
	return strings.Repeat("#", level) + " " + strings.TrimSpace(title)
}

func Bullet(items []string) string {
	var b strings.Builder
	for _, it := range items {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		b.WriteString("- ")
		b.WriteString(it)
		b.WriteByte('\n')
	}
	return b.String()
}

func StripTags(s string) string {
	var b strings.Builder
	in := false
	for _, r := range s {
		if r == '<' {
			in = true
			continue
		}
		if r == '>' {
			in = false
			continue
		}
		if !in {
			b.WriteRune(r)
		}
	}
	return CollapseSpace(b.String())
}

func IsIdent(s string) bool {
	if s == "" {
		return false
	}
	for i, r := range s {
		if i == 0 && unicode.IsDigit(r) {
			return false
		}
		ok := unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_'
		if !ok {
			return false
		}
	}
	return true
}

func Lines(s string) []string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")
	return strings.Split(s, "\n")
}

func NonEmptyLines(s string) []string {
	var out []string
	for _, ln := range Lines(s) {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			out = append(out, ln)
		}
	}
	return out
}
