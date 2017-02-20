package main

import (
"testing"
)

func TestLRedisLinkStorage_CreateAndExpandSlug(t *testing.T) {
	s := NewRedisLinkStorage("redis:6379")
	s.ClearAll()

	url := "http://example.com"
	slug := s.GetUrl(url)
	if slug != "" {
		t.Errorf("Expected no slug")
	}
	slug = "abc"
	stored := s.Store(slug, url)
	if !stored {
		t.Errorf("Exptected to store")
	}

	urlR := s.GetUrl(slug)
	if urlR != url {
		t.Errorf("Expected to get %s, but was %s", url, urlR)
	}

	stored = s.Store(slug, url)
	if stored {
		t.Errorf("Exptected not to store")
	}

}
