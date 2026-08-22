package validate

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Title(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return fmt.Errorf("title required")
	}
	if utf8.RuneCountInString(s) > 160 {
		return fmt.Errorf("title too long")
	}
	return nil
}

func Body(s string, max int) error {
	if max <= 0 {
		max = 1 << 20
	}
	if len(s) > max {
		return fmt.Errorf("body too large")
	}
	return nil
}

func Token(s string) error {
	s = strings.TrimSpace(s)
	if len(s) < 4 {
		return fmt.Errorf("token too short")
	}
	return nil
}

func Limit(n, max int) int {
	if n <= 0 {
		return max
	}
	if n > max {
		return max
	}
	return n
}
