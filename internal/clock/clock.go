package clock

import "time"

type Clock interface {
	Now() time.Time
}

type Real struct{}

func (Real) Now() time.Time { return time.Now().UTC() }

type Frozen struct{ T time.Time }

func (f Frozen) Now() time.Time { return f.T }

func RFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func ParseRFC3339(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t, nil
	}
	return time.Parse("2006-01-02 15:04:05", s)
}

func SinceMS(start time.Time) int64 {
	return time.Since(start).Milliseconds()
}
