package report

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"example.com/gridpulse/internal/model"
)

func CSV(recs []model.Record) string {
	var b strings.Builder
	b.WriteString("id,slug,title,owner,bytes,tags\n")
	for _, rec := range recs {
		fmt.Fprintf(&b, "%d,%s,%s,%d,%d,%s\n", rec.ID, rec.Slug, quote(rec.Title), rec.OwnerID, rec.Bytes, strings.Join(rec.Tags, "|"))
	}
	return b.String()
}

func quote(s string) string {
	if strings.ContainsAny(s, ",\"") {
		return `"` + strings.ReplaceAll(s, `"`, `""`) + `"`
	}
	return s
}

func JSON(recs []model.Record) ([]byte, error) {
	return json.Marshal(recs)
}

func SortByTitle(recs []model.Record) []model.Record {
	out := append([]model.Record(nil), recs...)
	sort.Slice(out, func(i, j int) bool { return out[i].Title < out[j].Title })
	return out
}

func FilterTag(recs []model.Record, tag string) []model.Record {
	if tag == "" {
		return recs
	}
	var out []model.Record
	for _, rec := range recs {
		for _, t := range rec.Tags {
			if t == tag {
				out = append(out, rec)
				break
			}
		}
	}
	return out
}
