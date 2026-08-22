package report

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/textfmt"
)

type Summary struct {
	Generated time.Time      `json:"generated"`
	Total     int            `json:"total"`
	Bytes     int            `json:"bytes"`
	Tags      map[string]int `json:"tags"`
	Owners    map[int64]int  `json:"owners"`
	TopTitles []string       `json:"top_titles"`
	Excerpt   string         `json:"excerpt"`
}

func Build(recs []model.Record) Summary {
	s := Summary{
		Generated: time.Now().UTC(),
		Total:     len(recs),
		Tags:      map[string]int{},
		Owners:    map[int64]int{},
	}
	var bodies []string
	for _, rec := range recs {
		s.Bytes += rec.Bytes
		if rec.Bytes == 0 {
			s.Bytes += len(rec.Body)
		}
		s.Owners[rec.OwnerID]++
		for _, tag := range rec.Tags {
			s.Tags[tag]++
		}
		bodies = append(bodies, rec.Title+": "+textfmt.Excerpt(rec.Body, 40))
		if len(s.TopTitles) < 8 {
			s.TopTitles = append(s.TopTitles, rec.Title)
		}
	}
	s.Excerpt = strings.Join(bodies, " | ")
	return s
}

func (s Summary) Markdown() string {
	var b strings.Builder
	b.WriteString(textfmt.Heading(1, "GridPulse Relay report"))
	b.WriteByte('\n')
	fmt.Fprintf(&b, "Generated: %s\n\n", s.Generated.Format(time.RFC3339))
	fmt.Fprintf(&b, "Total records: %d\nBytes: %d\n\n", s.Total, s.Bytes)
	b.WriteString(textfmt.Heading(2, "Tags"))
	b.WriteByte('\n')
	type kv struct {
		k string
		v int
	}
	var tags []kv
	for k, v := range s.Tags {
		tags = append(tags, kv{k, v})
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].v > tags[j].v })
	var bullets []string
	for _, t := range tags {
		bullets = append(bullets, fmt.Sprintf("%s (%d)", t.k, t.v))
	}
	b.WriteString(textfmt.Bullet(bullets))
	b.WriteString("\n")
	b.WriteString(textfmt.Heading(2, "Titles"))
	b.WriteByte('\n')
	b.WriteString(textfmt.Bullet(s.TopTitles))
	return b.String()
}

func Merge(a, b Summary) Summary {
	out := a
	if out.Tags == nil {
		out.Tags = map[string]int{}
	}
	if out.Owners == nil {
		out.Owners = map[int64]int{}
	}
	out.Total += b.Total
	out.Bytes += b.Bytes
	for k, v := range b.Tags {
		out.Tags[k] += v
	}
	for k, v := range b.Owners {
		out.Owners[k] += v
	}
	out.TopTitles = append(out.TopTitles, b.TopTitles...)
	if len(out.TopTitles) > 16 {
		out.TopTitles = out.TopTitles[:16]
	}
	return out
}
