package workflow

import (
	"fmt"
	"strings"
	"time"

	"example.com/gridpulse/internal/model"
	"example.com/gridpulse/internal/textfmt"
)

type Step struct {
	Name      string
	OK        bool
	Detail    string
	Started   time.Time
	Finished  time.Time
	DependsOn []string
}

type Pipeline struct {
	Name  string
	Steps []Step
}

func New(name string) *Pipeline {
	return &Pipeline{Name: name}
}

func (p *Pipeline) Add(name string, deps ...string) {
	p.Steps = append(p.Steps, Step{Name: name, DependsOn: append([]string(nil), deps...)})
}

func (p *Pipeline) Run(fn func(name string) error) error {
	done := map[string]bool{}
	for i := range p.Steps {
		st := &p.Steps[i]
		for _, d := range st.DependsOn {
			if !done[d] {
				return fmt.Errorf("missing dependency %s for %s", d, st.Name)
			}
		}
		st.Started = time.Now().UTC()
		if err := fn(st.Name); err != nil {
			st.OK = false
			st.Detail = err.Error()
			st.Finished = time.Now().UTC()
			return err
		}
		st.OK = true
		st.Detail = "ok"
		st.Finished = time.Now().UTC()
		done[st.Name] = true
	}
	return nil
}

func (p *Pipeline) Report() string {
	var b strings.Builder
	b.WriteString(textfmt.Heading(1, p.Name))
	b.WriteByte('\n')
	var items []string
	for _, st := range p.Steps {
		mark := "FAIL"
		if st.OK {
			mark = "OK"
		}
		items = append(items, fmt.Sprintf("%s %s (%s)", mark, st.Name, st.Detail))
	}
	b.WriteString(textfmt.Bullet(items))
	return b.String()
}

func IndexRecords(recs []model.Record) map[string]model.Record {
	out := map[string]model.Record{}
	for _, rec := range recs {
		out[rec.Slug] = rec
	}
	return out
}

func DiffTitles(old, next []model.Record) (added, removed []string) {
	a := map[string]bool{}
	b := map[string]bool{}
	for _, rec := range old {
		a[rec.Title] = true
	}
	for _, rec := range next {
		b[rec.Title] = true
	}
	for t := range b {
		if !a[t] {
			added = append(added, t)
		}
	}
	for t := range a {
		if !b[t] {
			removed = append(removed, t)
		}
	}
	return added, removed
}

func HealthySeed() []model.Record {
	return []model.Record{
		{Slug: "welcome", Title: "Welcome", Body: "getting started", Tags: []string{"docs"}},
		{Slug: "ops", Title: "Ops", Body: "runbook", Tags: []string{"ops"}},
		{Slug: "faq", Title: "FAQ", Body: "questions", Tags: []string{"docs"}},
	}
}
