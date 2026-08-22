package search

import (
	"sort"
	"strings"

	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/textfmt"
)

type Index struct {
	recs []model.Record
	inv  map[string][]int
}

func tokenize(s string) []string {
	s = strings.ToLower(s)
	var out []string
	var b strings.Builder
	flush := func() {
		if b.Len() == 0 {
			return
		}
		out = append(out, b.String())
		b.Reset()
	}
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			continue
		}
		flush()
	}
	flush()
	return out
}

func Build(recs []model.Record) *Index {
	idx := &Index{recs: recs, inv: map[string][]int{}}
	for i, rec := range recs {
		seen := map[string]bool{}
		blob := rec.Title + " " + rec.Body + " " + strings.Join(rec.Tags, " ")
		for _, tok := range tokenize(blob) {
			if seen[tok] {
				continue
			}
			seen[tok] = true
			idx.inv[tok] = append(idx.inv[tok], i)
		}
	}
	return idx
}

func (idx *Index) Find(q string) []model.Record {
	q = strings.TrimSpace(q)
	if q == "" {
		return append([]model.Record(nil), idx.recs...)
	}
	toks := tokenize(q)
	if len(toks) == 0 {
		return nil
	}
	counts := map[int]int{}
	for _, tok := range toks {
		for _, i := range idx.inv[tok] {
			counts[i]++
		}
	}
	var hits []model.Record
	for i, rec := range idx.recs {
		if counts[i] == len(toks) {
			hits = append(hits, rec)
		}
	}
	if len(hits) == 0 {
		ql := strings.ToLower(q)
		for _, rec := range idx.recs {
			blob := strings.ToLower(rec.Title + " " + rec.Body)
			if strings.Contains(blob, ql) {
				hits = append(hits, rec)
			}
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		return textfmt.WordCount(hits[i].Body) > textfmt.WordCount(hits[j].Body)
	})
	return hits
}
