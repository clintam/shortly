package main

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"log"
)

const REDIS_HASH_NAME = "shortly.slugToUrl"

type RedisLinkStorage struct {
	pool *redis.Pool
}

func NewRedisLinkStorage(addr string) *RedisLinkStorage {
	s := RedisLinkStorage{}
	s.pool = &redis.Pool{
		MaxIdle: 3,
		IdleTimeout: 240 * time.Second,
		Dial: func () (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
	return &s
}

func (s *RedisLinkStorage) Store(slug string, url string) bool {
	c := s.pool.Get()
	if c.Err() != nil {
		log.Printf("ERROR %s", c.Err())
		return false
	}
	defer c.Close()
	n, err := c.Do("HSETNX", REDIS_HASH_NAME, slug, url)
	if err != nil {
		log.Printf("ERROR %s", err)
		return false
	}
	r := n.(int64)
	return r == 1
}

func (s *RedisLinkStorage) GetUrl(slug string) string {
	c := s.pool.Get()
	if c.Err() != nil {
		log.Printf("ERROR %s", c.Err())
		return ""
	}
	defer c.Close()

	n, err := c.Do("HGET", REDIS_HASH_NAME, slug)
	if err != nil {
		log.Printf("ERROR %s", err)
		return ""
	}

	if n == nil {
		return ""
	}

	bytes := n.([]byte)
	return string(bytes)
}

func (s *RedisLinkStorage) ClearAll() {
	c, err := s.pool.Dial()
	if err != nil {
		log.Printf("ERROR %s", err)
		return
	}
	_ , err = c.Do("DEL", REDIS_HASH_NAME)
	if err != nil {
		log.Printf("ERROR %s", err)
	}
}

