package libsysinfo

import (
	"sync"
)

type cachedValues struct {
	set map[string]string
	mu  sync.RWMutex
}

func newCachedValues(l int) *cachedValues {
	return &cachedValues{
		set: make(map[string]string, l),
	}
}

func (s *cachedValues) Set(k string, v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.set[k] = v
}

func (s *cachedValues) Get(k string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess := s.set[k]
	return sess
}

func (s *cachedValues) Delete(k string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.set, k)
}

func (s *cachedValues) Exists(k string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	v, present := s.set[k]
	return v, present
}

func (s *cachedValues) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.set
}

func (s *cachedValues) Empty() {
	for k, _ := range s.All() {
		s.Delete(k)
	}
}
