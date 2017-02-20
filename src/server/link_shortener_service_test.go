package main

import (
	"testing"
)

func TestLinkShortenService_CreateAndExpandSlug(t *testing.T) {
	s := LinkShortenService{NewMemoryLinkStorage()}

	url := "http://example.com/foobar"
	slug := s.CreateSlug(url)
	// TODO assert on slug
	expaneded := s.ExpandSlug(slug)
	if expaneded != url {
		t.Errorf("expected %s to be %s", expaneded, url)
	}
}

func TestLinkShortenService_CreateTwoFromSameUrl(t *testing.T) {
	s := LinkShortenService{NewMemoryLinkStorage()}

	url := "http://example.com/foobar"
	slug1 := s.CreateSlug(url)
	slug2 := s.CreateSlug(url)
	if slug1 == slug2 {
		t.Errorf("expected different slugs %s and %s", slug1, slug2)
	}
}
