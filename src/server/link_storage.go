package main

import "sync"

type LinkStorage interface {
	// Store attempts to store a new slug mapped to the url.
	// If the slug already exists, it should return false and not store anything.
	Store(slug string, url string) bool
	GetUrl(slug string) string
}

type MemoryLinkStorage struct {
	slugToUrl map[string]string
	mutex sync.RWMutex
}

func NewMemoryLinkStorage() *MemoryLinkStorage {
	s := MemoryLinkStorage{}
	s.slugToUrl = make(map[string]string)
	return &s
}

func (s *MemoryLinkStorage) Store(slug string, url string) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, ok := s.slugToUrl[slug]; ok {
		return false
	}
	s.slugToUrl[slug] = url
	return true
}

func (s *MemoryLinkStorage) GetUrl(slug string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.slugToUrl[slug]
}
