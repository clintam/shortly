package main

import (
	"crypto/md5"
	"encoding/binary"
	"strconv"
	"log"
)

const SLUG_SIZE = 3 // TODO: This will need to increase to scale up (lets let tests discover issue)

type LinkShortenService struct {
	storage LinkStorage
}

func (s *LinkShortenService) CreateSlug(url string) string {
	generateHash := func(seed int) string {
		hasher := md5.New()
		hasher.Write([]byte(url))
		hasher.Write([]byte(strconv.Itoa(seed)))
		intVal := binary.BigEndian.Uint64(hasher.Sum(nil))
		return strconv.FormatUint(intVal, 36)[:SLUG_SIZE]
	}

	seed := 0
	slug := generateHash(seed)
	for ; !s.storage.Store(slug, url); slug = generateHash(seed) {
		seed++
		log.Printf("Conflicting slug [%s], trying again with seed [%d]", slug, seed)
	}
	log.Printf("Mapped slug [%s] to url [%s]", slug, url)
	return slug
}

func (s *LinkShortenService) ExpandSlug(url string) string {
	return s.storage.GetUrl(url)
}
