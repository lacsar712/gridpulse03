package notify

import (
	"log"
	"strings"
	"time"
)

type Event struct {
	At      time.Time
	Level   string
	Message string
}

type Sink struct {
	events []Event
}

func (s *Sink) Info(msg string) {
	s.add("info", msg)
}

func (s *Sink) Error(msg string) {
	s.add("error", msg)
}

func (s *Sink) add(level, msg string) {
	s.events = append(s.events, Event{At: time.Now().UTC(), Level: level, Message: strings.TrimSpace(msg)})
	log.Printf("%s %s", strings.ToUpper(level), msg)
}

func (s *Sink) Last() (Event, bool) {
	if len(s.events) == 0 {
		return Event{}, false
	}
	return s.events[len(s.events)-1], true
}

func (s *Sink) Filter(level string) []Event {
	var out []Event
	for _, e := range s.events {
		if e.Level == level {
			out = append(out, e)
		}
	}
	return out
}
