package workflow

import (
	"fmt"
	"sort"
	"strings"

	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/textfmt"
)

func Plan(names []string) *Pipeline {
	p := New("gridpulse-plan")
	prev := ""
	for _, n := range names {
		if prev == "" {
			p.Add(n)
		} else {
			p.Add(n, prev)
		}
		prev = n
	}
	return p
}

func GroupByTag(recs []model.Record) map[string][]model.Record {
	out := map[string][]model.Record{}
	for _, rec := range recs {
		if len(rec.Tags) == 0 {
			out["untagged"] = append(out["untagged"], rec)
			continue
		}
		for _, tag := range rec.Tags {
			out[tag] = append(out[tag], rec)
		}
	}
	return out
}

func TagList(recs []model.Record) []string {
	seen := map[string]bool{}
	var out []string
	for _, rec := range recs {
		for _, tag := range rec.Tags {
			if seen[tag] {
				continue
			}
			seen[tag] = true
			out = append(out, tag)
		}
	}
	sort.Strings(out)
	return out
}

func FormatGroup(recs []model.Record) string {
	g := GroupByTag(recs)
	keys := make([]string, 0, len(g))
	for k := range g {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(textfmt.Heading(2, k))
		b.WriteByte('\n')
		var titles []string
		for _, rec := range g[k] {
			titles = append(titles, fmt.Sprintf("%s (%s)", rec.Title, rec.Slug))
		}
		b.WriteString(textfmt.Bullet(titles))
		b.WriteByte('\n')
	}
	return b.String()
}

func Titles(recs []model.Record) []string {
	out := make([]string, 0, len(recs))
	for _, rec := range recs {
		out = append(out, rec.Title)
	}
	return out
}
