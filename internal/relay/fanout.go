package relay

import "fmt"

func FanoutMark(targets []string, fn func(string) error) (done bool, err error) {
	if len(targets) == 0 {
		return true, nil
	}
	for _, t := range targets {
		if e := fn(t); e != nil {
			return false, fmt.Errorf("%s fanout %s: %w", "gridpulse", t, e)
		}
	}
	return true, nil
}
