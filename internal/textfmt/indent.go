package textfmt

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
)

func Fingerprint(s string) string {
	sum := sha1.Sum([]byte(s))
	return hex.EncodeToString(sum[:8])
}

func NormalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}

func Indent(s string, n int) string {
	if n < 0 {
		n = 0
	}
	pad := strings.Repeat(" ", n)
	lines := Lines(NormalizeNewlines(s))
	for i, ln := range lines {
		if ln == "" {
			continue
		}
		lines[i] = pad + ln
	}
	return strings.Join(lines, "\n")
}

func Dedent(s string) string {
	lines := Lines(NormalizeNewlines(s))
	min := -1
	for _, ln := range lines {
		if strings.TrimSpace(ln) == "" {
			continue
		}
		n := 0
		for _, r := range ln {
			if r == ' ' || r == '\t' {
				n++
				continue
			}
			break
		}
		if min < 0 || n < min {
			min = n
		}
	}
	if min <= 0 {
		return strings.Join(lines, "\n")
	}
	for i, ln := range lines {
		if len(ln) >= min {
			lines[i] = ln[min:]
		}
	}
	return strings.Join(lines, "\n")
}

func QuoteLines(s string) string {
	var b strings.Builder
	for _, ln := range Lines(s) {
		b.WriteString("> ")
		b.WriteString(ln)
		b.WriteByte('\n')
	}
	return b.String()
}
