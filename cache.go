package libsysinfo

import (
	"sync"
)

type CachedValues struct {
	set map[string]string
	mu  sync.RWMutex
}

func NewCachedValues(l int) *CachedValues {
	return &CachedValues{
		set: make(map[string]string, l),
	}
}

func (s *CachedValues) Set(k string, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[k] = v
}

func (s *CachedValues) Get(k string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess := s.set[k]
	return sess
}

func (s *CachedValues) Delete(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, k)
}

func (s *CachedValues) Exists(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, present := s.set[k]
	return v, present
}

func (s *CachedValues) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set
}

func (s *CachedValues) Empty() {
	for k, _ := range s.All() {
		s.Delete(k)
	}
}
