package main

import (
	"crypto/md5"
	"encoding/binary"
	"strconv"
	"log"
	"math/rand"
	"errors"
)

const SLUG_SIZE = 7

type LinkShortenService struct {
	storage LinkStorage
}



func (s *LinkShortenService) CreateSlug(url string, cutstomSlug string) (string, error) {
	isCustom := cutstomSlug != ""
	generateSlug := func() string {
		if isCustom {
			return cutstomSlug
		}
		hasher := md5.New()
		hasher.Write([]byte(url))
		hasher.Write([]byte(strconv.Itoa(rand.Int())))
		intVal := binary.BigEndian.Uint64(hasher.Sum(nil))
		return strconv.FormatUint(intVal, 36)[:SLUG_SIZE]
	}

	maxAttempts := 10
	maxAttemptMessage := "Pardon me, but we are experiancing some difficutlites"
	if isCustom {
		maxAttemptMessage = "Custom url is already taken"
		maxAttempts = 1
	}
	attempt := 0
	slug := generateSlug()
	for ; !s.storage.Store(slug, url); slug = generateSlug() {
		attempt++
		if attempt >= maxAttempts {
			return "", errors.New(maxAttemptMessage)
		}
		log.Printf("Conflicting slug [%s] on attempt [%d], trying again with another seed", slug, attempt)
	}
	log.Printf("Mapped slug [%s] to url [%s]", slug, url)
	return slug, nil
}

func (s *LinkShortenService) ExpandSlug(url string) string {
	return s.storage.GetUrl(url)
}
