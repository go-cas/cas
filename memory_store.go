package cas

import (
	"sync"
)

type MemoryStore struct {
	mu    sync.RWMutex
	store map[string]*AuthenticationResponse
}

func (s *MemoryStore) Read(id string) (*AuthenticationResponse, error) {
	s.mu.RLock()

	if s.store == nil {
		s.mu.RUnlock()
		return nil, ErrInvalidTicket
	}

	t, ok := s.store[id]
	s.mu.RUnlock()

	if !ok {
		return nil, ErrInvalidTicket
	}

	return t, nil
}

func (s *MemoryStore) Write(id string, ticket *AuthenticationResponse) error {
	s.mu.Lock()

	if s.store == nil {
		s.store = make(map[string]*AuthenticationResponse)
	}

	s.store[id] = ticket

	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) Delete(id string) error {
	s.mu.Lock()
	delete(s.store, id)
	s.mu.Unlock()
	return nil
}

func (s *MemoryStore) Clear() error {
	s.mu.Lock()
	s.store = nil
	s.mu.Unlock()
	return nil
}
