package relay

import (
	"fmt"
	"strconv"
	"strings"
)

type Rec struct {
	Title, Body string
	Tags        []string
}

func Sample() Rec {
	return Rec{Title: "gridpulse-sample", Body: "Trip=TG01-12 level=0.5 epoch=4", Tags: []string{"tide"}}
}

func Seed() []Rec {
	return []Rec{
		Sample(),
		{Title: "gridpulse-seed-b", Body: "Trip=TG01-12 level=0.5 epoch=4 b", Tags: []string{"tide"}},
	}
}

func Steps() []string { return []string{"validate", "dispatch", "commit", "export"} }

func Enforce(title, body string, tags []string) error {
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("Trip title required")
	}
	if strings.TrimSpace(body) == "" {
		return fmt.Errorf("Trip body required")
	}
	if len(tags) == 0 {
		return fmt.Errorf("Trip tag required")
	}
	return nil
}

func AfterWrite(getMin func() (string, error), setMin func(string) error, body string) error {
	n := epochOf(body)
	if n < 0 {
		return fmt.Errorf("epoch missing in %q", body)
	}
	cur, err := getMin()
	if err == nil && strings.TrimSpace(cur) != "" {
		old, conv := strconv.Atoi(strings.TrimSpace(cur))
		if conv == nil && n < old {
			return fmt.Errorf("epoch %d < committed %d", n, old)
		}
	}
	return setMin(strconv.Itoa(n))
}

func epochOf(body string) int {
	for _, part := range strings.Fields(body) {
		k, v, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		if k == "epoch" || k == "locked" || k == "bay" {
			n, err := strconv.Atoi(v)
			if err == nil {
				return n
			}
		}
	}
	return -1
}
