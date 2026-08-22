package relay

import "fmt"

type Item struct {
	ID, Body string
}

type Validator interface {
	Validate(Item) error
}

type Store interface {
	Commit(Item) error
	Last() string
}

type MemStore struct {
	last string
	log  []Item
}

func (s *MemStore) Commit(it Item) error {
	s.last = it.ID
	s.log = append(s.log, it)
	return nil
}

func (s *MemStore) Last() string { return s.last }

type BodyValidator struct{}

func (BodyValidator) Validate(it Item) error {
	if it.Body == "" {
		return fmt.Errorf("empty Trip body")
	}
	if it.ID == "" {
		return fmt.Errorf("empty Trip id")
	}
	return nil
}

func ApplyValidated(v Validator, s Store, it Item) error {
	if err := v.Validate(it); err != nil {
		return err
	}
	return s.Commit(it)
}
