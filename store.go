package main

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrNotFound       = errors.New("applicant not found")
	ErrDuplicateEmail = errors.New("email already exists")
)

type Store struct {
	mu      sync.RWMutex
	nextID  int
	byID    map[int]Applicant
	byEmail map[string]int
}

func NewStore() *Store {
	return &Store{nextID: 1, byID: map[int]Applicant{}, byEmail: map[string]int{}}
}

func (s *Store) Create(a Applicant) (Applicant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := strings.ToLower(strings.TrimSpace(a.Email))
	if _, exists := s.byEmail[key]; exists {
		return Applicant{}, ErrDuplicateEmail
	}

	a.ID = s.nextID
	s.nextID++

	now := time.Now().UTC()
	a.CreatedAt = now
	a.UpdatedAt = now

	s.byID[a.ID] = a
	s.byEmail[key] = a.ID

	return a, nil
}

func (s *Store) Get(id int) (Applicant, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	a, ok := s.byID[id]
	if !ok {
		return Applicant{}, ErrNotFound
	}

	return a, nil
}

func (s *Store) List() []Applicant {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Applicant, 0, len(s.byID))
	for id := 1; id < s.nextID; id++ {
		if a, ok := s.byID[id]; ok {
			out = append(out, a)
		}
	}

	return out
}

func (s *Store) Transition(id int, fn func(Applicant) (Applicant, error)) (Applicant, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	a, ok := s.byID[id]
	if !ok {
		return Applicant{}, ErrNotFound
	}

	updated, err := fn(a)
	if err != nil {
		return Applicant{}, err
	}

	updated.UpdatedAt = time.Now().UTC()
	s.byID[id] = updated

	return updated, nil
}
