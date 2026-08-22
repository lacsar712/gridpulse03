package model

import "strings"

type ListFilter struct {
	Query    string
	Tag      string
	OwnerID  int64
	Limit    int
	Offset   int
	OrderAsc bool
}

func (f ListFilter) Normalized(pageSize int) ListFilter {
	if f.Limit <= 0 || f.Limit > pageSize {
		f.Limit = pageSize
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	f.Query = strings.TrimSpace(f.Query)
	f.Tag = strings.TrimSpace(f.Tag)
	return f
}

func SplitTags(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	seen := map[string]bool{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" || seen[p] {
			continue
		}
		seen[p] = true
		out = append(out, p)
	}
	return out
}

func JoinTags(tags []string) string {
	return strings.Join(SplitTags(strings.Join(tags, ",")), ",")
}

func Slugify(s string) string {
	s = strings.TrimSpace(strings.ToLower(s))
	var b strings.Builder
	prevDash := false
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			prevDash = false
			continue
		}
		if prevDash {
			continue
		}
		b.WriteByte('-')
		prevDash = true
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "item"
	}
	return out
}
