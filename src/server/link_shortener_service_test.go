package main

import (
	"testing"
	"math/rand"
)


func TestLinkShortenService_CreateAndExpandSlug(t *testing.T) {
	s := LinkShortenService{NewMemoryLinkStorage()}

	url := "http://example.com/foobar"
	slug, _ := s.CreateSlug(url, "")
	// TODO assert on slug
	expaneded := s.ExpandSlug(slug)
	if expaneded != url {
		t.Errorf("expected %s to be %s", expaneded, url)
	}
}

func TestLinkShortenService_CreateTwoFromSameUrl(t *testing.T) {
	s := LinkShortenService{NewMemoryLinkStorage()}
	randSeed := int64(0)
	url := "http://example.com/foobar"
	rand.Seed(randSeed)
	slug1, _ := s.CreateSlug(url, "")
	rand.Seed(randSeed)
	slug2, _ := s.CreateSlug(url, "")
	if slug1 == slug2 {
		t.Errorf("expected different slugs %s and %s", slug1, slug2)
	}
}

func TestLinkShortenService_CreateAndExpandCustomSlug(t *testing.T) {
	s := LinkShortenService{NewMemoryLinkStorage()}

	url := "http://example.com/foobar"
	customSlug := "myCustomName"
	slug, _ := s.CreateSlug(url, customSlug)
	if slug != customSlug {
		t.Errorf("expected %s to be %s", slug, customSlug)
	}
	expaneded := s.ExpandSlug(slug)
	if expaneded != url {
		t.Errorf("expected %s to be %s", expaneded, url)
	}
	_, err := s.CreateSlug(url, customSlug)
	if err == nil {
		t.Errorf("Expected error")
	}
}