package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	externalBaseUrl := flag.String("external-url", "http://localhost:8080", "Externally exposed url")
	storage := flag.String("storage", "redis", "Persistent storage. Valid values are redis,memory")
	redisAddr := flag.String("redis", "redis:6379", "Reids address <host>:<port>")
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	flag.Parse()
	s := createServer(getStorage(*storage, *redisAddr), *port, *externalBaseUrl)
	s.Start()
}

func getStorage(name string, redisAddr string) LinkStorage {
	switch(name) {
	case "memory":
		return NewMemoryLinkStorage()
	case "redis":
		return NewRedisLinkStorage(redisAddr)
	default:
		panic("No such storage: " + name)
	}
}
