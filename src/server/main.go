package main

import (
	"flag"
	"log"
	"os"
)

func main() {
	port := flag.String("port", "8080", "Port to listen on")
	externalBaseUrl := flag.String("external-url", "http://localhost:8080", "Externally exposed url")
	storage := flag.String("storage", "memory", "Persistent storage. Valid values are redis,mongo,memory")
	redisAddr := flag.String("redis", "redis:6379", "Redis address <host>:<port>")
	mongoAddr := flag.String("mongo", "mongo", "Mongo comma seperated hosts addresses")
	log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lmicroseconds)
	flag.Parse()

	s := createServer(getStorage(*storage, *redisAddr, *mongoAddr), *port, *externalBaseUrl)
	s.Start()
}

func getStorage(name string, redisAddr string, mongoAddr string) LinkStorage {
	switch(name) {
	case "memory":
		return NewMemoryLinkStorage()
	case "redis":
		return NewRedisLinkStorage(redisAddr)
	case "mongo":
		return NewMongoLinkStorage(mongoAddr)
	default:
		panic("No such storage: " + name)
	}
}
