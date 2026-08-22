package textfmt

import (
	"strconv"
	"strings"
)

func SplitCSVLine(line string) []string {
	var out []string
	var b strings.Builder
	inQ := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if inQ {
			if c == '"' {
				if i+1 < len(line) && line[i+1] == '"' {
					b.WriteByte('"')
					i++
					continue
				}
				inQ = false
				continue
			}
			b.WriteByte(c)
			continue
		}
		switch c {
		case '"':
			inQ = true
		case ',':
			out = append(out, b.String())
			b.Reset()
		default:
			b.WriteByte(c)
		}
	}
	out = append(out, b.String())
	return out
}

func JoinCSV(fields []string) string {
	var b strings.Builder
	for i, f := range fields {
		if i > 0 {
			b.WriteByte(',')
		}
		need := strings.ContainsAny(f, ",\"\n")
		if need {
			b.WriteByte('"')
			b.WriteString(strings.ReplaceAll(f, `"`, `""`))
			b.WriteByte('"')
		} else {
			b.WriteString(f)
		}
	}
	return b.String()
}

func IntsJoin(ns []int, sep string) string {
	parts := make([]string, len(ns))
	for i, n := range ns {
		parts[i] = strconv.Itoa(n)
	}
	return strings.Join(parts, sep)
}
